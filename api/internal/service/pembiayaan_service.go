package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/pembiayaan"
	"github.com/bmt-saas/api/internal/domain/rekening"
	"github.com/google/uuid"
)

// PembiayaanService mengelola state machine pembiayaan.
//
// State machine:
//
//	PENGAJUAN → ANALISIS → KOMITE → AKAD → PENCAIRAN → AKTIF → LUNAS/MACET
type PembiayaanService struct {
	repo            pembiayaan.Repository
	rekeningService *RekeningService
	akuntansiSvc    *AkuntansiService
	notifikasiSvc   *NotifikasiService
}

// NewPembiayaanService membuat instance baru PembiayaanService.
func NewPembiayaanService(
	repo pembiayaan.Repository,
	rekeningService *RekeningService,
	akuntansiSvc *AkuntansiService,
	notifikasiSvc *NotifikasiService,
) *PembiayaanService {
	return &PembiayaanService{
		repo:            repo,
		rekeningService: rekeningService,
		akuntansiSvc:    akuntansiSvc,
		notifikasiSvc:   notifikasiSvc,
	}
}

// ─── Transisi yang diizinkan ─────────────────────────────────────────────────

// transiziDiizinkan mendefinisikan state machine pembiayaan.
// Setiap key adalah status asal, nilai adalah slice status tujuan yang valid.
var transiziDiizinkan = map[pembiayaan.StatusPembiayaan][]pembiayaan.StatusPembiayaan{
	pembiayaan.StatusPengajuan: {"ANALISIS"},
	"ANALISIS":                 {"KOMITE", pembiayaan.StatusPengajuan},
	"KOMITE":                   {"AKAD", "ANALISIS"},
	"AKAD":                     {"PENCAIRAN"},
	"PENCAIRAN":                {pembiayaan.StatusAktif},
}

func bolehTransisi(dari, ke pembiayaan.StatusPembiayaan) bool {
	allowed, ok := transiziDiizinkan[dari]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == ke {
			return true
		}
	}
	return false
}

// ─── Public Methods ───────────────────────────────────────────────────────────

// AjukanPembiayaan membuat pengajuan pembiayaan baru dengan status PENGAJUAN.
func (s *PembiayaanService) AjukanPembiayaan(ctx context.Context, input pembiayaan.CreatePembiayaanInput) (*pembiayaan.Pembiayaan, error) {
	nomor, err := s.repo.GenerateNomor(ctx, input.BMTID, input.CabangID)
	if err != nil {
		return nil, fmt.Errorf("gagal generate nomor pembiayaan: %w", err)
	}

	p, err := pembiayaan.NewPembiayaan(input, nomor)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat pembiayaan: %w", err)
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("gagal simpan pembiayaan: %w", err)
	}

	return p, nil
}

// GetByID mengambil pembiayaan berdasarkan ID dan memvalidasi kepemilikan BMT.
func (s *PembiayaanService) GetByID(ctx context.Context, id, bmtID uuid.UUID) (*pembiayaan.Pembiayaan, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil pembiayaan: %w", err)
	}
	if p.BMTID != bmtID {
		return nil, errors.New("pembiayaan tidak ditemukan")
	}
	return p, nil
}

// ListPembiayaan mengembalikan daftar pembiayaan sesuai filter.
func (s *PembiayaanService) ListPembiayaan(ctx context.Context, filter pembiayaan.ListPembiayaanFilter) ([]*pembiayaan.Pembiayaan, int64, error) {
	list, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal ambil daftar pembiayaan: %w", err)
	}
	return list, total, nil
}

// MajukanStatus memvalidasi transisi status lalu memprosesnya.
// State machine: PENGAJUAN → ANALISIS → KOMITE → AKAD → PENCAIRAN → AKTIF
func (s *PembiayaanService) MajukanStatus(
	ctx context.Context,
	id uuid.UUID,
	statusBaru pembiayaan.StatusPembiayaan,
	oleh uuid.UUID,
) (*pembiayaan.Pembiayaan, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil pembiayaan: %w", err)
	}

	if !bolehTransisi(p.Status, statusBaru) {
		return nil, fmt.Errorf("transisi status dari %s ke %s tidak diizinkan", p.Status, statusBaru)
	}

	if err := s.repo.UpdateStatus(ctx, id, statusBaru); err != nil {
		return nil, fmt.Errorf("gagal update status pembiayaan: %w", err)
	}

	p.Status = statusBaru
	p.UpdatedBy = oleh

	// Jika status menjadi AKTIF, generate jadwal angsuran
	if statusBaru == pembiayaan.StatusAktif {
		if err := s.GenerateJadwalAngsuran(ctx, p); err != nil {
			return nil, fmt.Errorf("gagal generate jadwal angsuran: %w", err)
		}
	}

	return p, nil
}

// GenerateJadwalAngsuran membuat seluruh baris angsuran berdasarkan jangka bulan pembiayaan.
// Untuk Murabahah: setiap angsuran berisi pokok + margin flat.
// Untuk Mudharabah/Musyarakah/Ijarah/Qardh: hanya pokok (margin dari realisasi bagi hasil).
func (s *PembiayaanService) GenerateJadwalAngsuran(ctx context.Context, p *pembiayaan.Pembiayaan) error {
	now := time.Now()
	// Hitung nominal pokok per bulan (flat)
	nominalPokok := p.Pokok / int64(p.JangkaBulan)
	// Sisanya ke angsuran terakhir agar total tepat
	sisaPokok := p.Pokok - nominalPokok*int64(p.JangkaBulan)

	// Nominal margin hanya untuk Murabahah
	nominalMargin := int64(0)
	sisaMargin := int64(0)
	if p.Akad == pembiayaan.AkadMurabahah && p.JangkaBulan > 0 {
		totalMargin := p.TotalKewajiban - p.Pokok
		nominalMargin = totalMargin / int64(p.JangkaBulan)
		sisaMargin = totalMargin - nominalMargin*int64(p.JangkaBulan)
	}

	for i := int16(1); i <= p.JangkaBulan; i++ {
		pokok := nominalPokok
		margin := nominalMargin

		// Sisa dimasukkan ke angsuran terakhir
		if i == p.JangkaBulan {
			pokok += sisaPokok
			margin += sisaMargin
		}

		tanggalJatuhTempo := now.AddDate(0, int(i), 0)

		a := &pembiayaan.AngsuranPembiayaan{
			ID:                uuid.New(),
			BMTID:             p.BMTID,
			PembiayaanID:      p.ID,
			PeriodeBulan:      i,
			NominalPokok:      pokok,
			NominalMargin:     margin,
			TotalAngsuran:     pokok + margin,
			TanggalJatuhTempo: tanggalJatuhTempo,
			NominalTerbayar:   0,
			Status:            "MENUNGGU",
			CreatedAt:         now,
		}

		if err := s.repo.CreateAngsuran(ctx, a); err != nil {
			return fmt.Errorf("gagal simpan angsuran periode %d: %w", i, err)
		}
	}
	return nil
}

// BayarAngsuran memproses pembayaran angsuran dari rekening nasabah.
// Urutan: lock pembiayaan → validasi aktif → cari angsuran MENUNGGU tertua →
// tarik rekening → update angsuran → update saldo pembiayaan → post jurnal → cek lunas.
func (s *PembiayaanService) BayarAngsuran(
	ctx context.Context,
	pembiayaanID, rekeningID uuid.UUID,
	nominal int64,
	oleh uuid.UUID,
) (*pembiayaan.AngsuranPembiayaan, error) {
	// 1. Lock pembiayaan agar tidak ada race condition
	p, err := s.repo.LockForUpdate(ctx, pembiayaanID)
	if err != nil {
		return nil, fmt.Errorf("gagal lock pembiayaan: %w", err)
	}

	// 2. Validasi pembiayaan harus berstatus AKTIF
	if err := p.ValidasiAktif(); err != nil {
		return nil, err
	}

	// 3. Ambil semua angsuran, cari yang MENUNGGU tertua (periode_bulan terkecil)
	angsuranList, err := s.repo.ListAngsuran(ctx, pembiayaanID)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil jadwal angsuran: %w", err)
	}

	var targetAngsuran *pembiayaan.AngsuranPembiayaan
	for _, a := range angsuranList {
		if a.Status == "MENUNGGU" || a.Status == "SEBAGIAN" || a.Status == "LEWAT" {
			targetAngsuran = a
			break
		}
	}
	if targetAngsuran == nil {
		return nil, errors.New("tidak ada angsuran yang perlu dibayar")
	}

	// 4. Validasi nominal tidak melebihi kewajiban angsuran
	sisaAngsuran := targetAngsuran.TotalAngsuran - targetAngsuran.NominalTerbayar
	if nominal > sisaAngsuran {
		return nil, fmt.Errorf("%w: nominal %d melebihi sisa angsuran %d",
			pembiayaan.ErrAngsuranMelebihiSaldo, nominal, sisaAngsuran)
	}

	// 5. Tarik dana dari rekening nasabah
	trx, err := s.rekeningService.Tarik(ctx, rekening.PenarikanInput{
		RekeningID: rekeningID,
		Nominal:    nominal,
		Keterangan: fmt.Sprintf("Angsuran pembiayaan %s periode %d", p.NomorPembiayaan, targetAngsuran.PeriodeBulan),
		CreatedBy:  oleh,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal tarik rekening untuk angsuran: %w", err)
	}

	// 6. Hitung proporsi pokok dan margin dari pembayaran
	var pokoBayar, marginBayar int64
	if targetAngsuran.TotalAngsuran > 0 {
		rasio := float64(nominal) / float64(targetAngsuran.TotalAngsuran)
		pokoBayar = int64(float64(targetAngsuran.NominalPokok) * rasio)
		marginBayar = nominal - pokoBayar
		// Pastikan margin tidak negatif
		if marginBayar < 0 {
			marginBayar = 0
			pokoBayar = nominal
		}
	} else {
		pokoBayar = nominal
	}

	// 7. Perbarui data angsuran di DB
	tanggalBayar := time.Now()
	if err := s.repo.UpdateAngsuranTerbayar(ctx, targetAngsuran.ID, nominal, tanggalBayar, trx.ID); err != nil {
		return nil, fmt.Errorf("gagal update angsuran terbayar: %w", err)
	}

	// Reload angsuran setelah update
	angsuranUpdated, err := s.repo.GetAngsuranByID(ctx, targetAngsuran.ID)
	if err != nil {
		return nil, fmt.Errorf("gagal reload angsuran: %w", err)
	}

	// 8. Perbarui saldo pembiayaan
	saldoPokokBaru := p.SaldoPokok - pokoBayar
	saldoMarginBaru := p.SaldoMargin - marginBayar
	if saldoPokokBaru < 0 {
		saldoPokokBaru = 0
	}
	if saldoMarginBaru < 0 {
		saldoMarginBaru = 0
	}
	if err := s.repo.UpdateSaldo(ctx, pembiayaanID, saldoPokokBaru, saldoMarginBaru); err != nil {
		return nil, fmt.Errorf("gagal update saldo pembiayaan: %w", err)
	}

	// 9. Post jurnal double-entry
	// Prinsip: Debit Kas(101) / Kredit Piutang Pembiayaan(131) untuk pokok
	//          Debit Kas(101) / Kredit Pendapatan Margin(411) untuk margin (jika ada)
	entries := []JurnalEntry{}
	if pokoBayar > 0 {
		entries = append(entries,
			JurnalEntry{KodeAkun: "101", Posisi: "DEBIT", Nominal: pokoBayar},
			JurnalEntry{KodeAkun: "131", Posisi: "KREDIT", Nominal: pokoBayar},
		)
	}
	if marginBayar > 0 && p.SaldoMargin > 0 {
		entries = append(entries,
			JurnalEntry{KodeAkun: "101", Posisi: "DEBIT", Nominal: marginBayar},
			JurnalEntry{KodeAkun: "411", Posisi: "KREDIT", Nominal: marginBayar},
		)
	}
	if len(entries) > 0 {
		_ = s.akuntansiSvc.PostJurnal(ctx, PostJurnalInput{
			BMTID:      p.BMTID,
			CabangID:   p.CabangID,
			Keterangan: fmt.Sprintf("Bayar angsuran %s periode %d", p.NomorPembiayaan, targetAngsuran.PeriodeBulan),
			Referensi:  trx.ID.String(),
			Entries:    entries,
		})
	}

	// 10. Cek apakah pembiayaan sudah lunas (saldo pokok == 0)
	if saldoPokokBaru == 0 {
		if err := s.repo.UpdateStatus(ctx, pembiayaanID, pembiayaan.StatusLunas); err != nil {
			return nil, fmt.Errorf("gagal update status lunas: %w", err)
		}

		// 11. Kirim notifikasi jika service tersedia
		if s.notifikasiSvc != nil {
			_ = s.notifikasiSvc.Kirim(ctx, KirimInput{
				BMTID:        p.BMTID,
				Channel:      "FCM",
				TemplateKode: "PEMBIAYAAN_LUNAS",
				Variables: map[string]string{
					"nomor_pembiayaan": p.NomorPembiayaan,
					"pesan":            fmt.Sprintf("Selamat! Pembiayaan %s Anda telah lunas.", p.NomorPembiayaan),
				},
			})
		}
	}

	return angsuranUpdated, nil
}

// GetJadwalAngsuran mengembalikan seluruh jadwal angsuran milik pembiayaan tertentu.
func (s *PembiayaanService) GetJadwalAngsuran(ctx context.Context, pembiayaanID, bmtID uuid.UUID) ([]*pembiayaan.AngsuranPembiayaan, error) {
	// Validasi kepemilikan BMT
	p, err := s.repo.GetByID(ctx, pembiayaanID)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil pembiayaan: %w", err)
	}
	if p.BMTID != bmtID {
		return nil, errors.New("pembiayaan tidak ditemukan")
	}

	angsuranList, err := s.repo.ListAngsuran(ctx, pembiayaanID)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil jadwal angsuran: %w", err)
	}
	return angsuranList, nil
}
