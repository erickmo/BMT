package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/finance"
	"github.com/bmt-saas/api/internal/domain/pembiayaan"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/google/uuid"
)

// ─── Laporan types ────────────────────────────────────────────────────────────

// LaporanNeraca adalah ringkasan aset, kewajiban, dan ekuitas.
type LaporanNeraca struct {
	Aset      int64  `json:"aset"`
	Kewajiban int64  `json:"kewajiban"`
	Ekuitas   int64  `json:"ekuitas"`
	Dari      string `json:"dari"`
	Sampai    string `json:"sampai"`
}

// LaporanSHU adalah ringkasan pendapatan dan biaya (Sisa Hasil Usaha).
type LaporanSHU struct {
	Pendapatan int64  `json:"pendapatan"`
	Biaya      int64  `json:"biaya"`
	SHU        int64  `json:"shu"`
	Dari       string `json:"dari"`
	Sampai     string `json:"sampai"`
}

// LaporanArusKas adalah ringkasan arus kas masuk dan keluar.
type LaporanArusKas struct {
	KasMasuk  int64  `json:"kas_masuk"`
	KasKeluar int64  `json:"kas_keluar"`
	NetKas    int64  `json:"net_kas"`
	Dari      string `json:"dari"`
	Sampai    string `json:"sampai"`
}

// KolektibilitasRingkasan merangkum jumlah dan outstanding per level kolektibilitas.
type KolektibilitasRingkasan struct {
	Level       int16  `json:"level"`
	Keterangan  string `json:"keterangan"`
	Jumlah      int    `json:"jumlah"`
	Outstanding int64  `json:"outstanding"`
}

// LaporanKolektibilitas adalah laporan kolektibilitas pembiayaan per BMT.
type LaporanKolektibilitas struct {
	BMTID  uuid.UUID                `json:"bmt_id"`
	Data   []KolektibilitasRingkasan `json:"data"`
	Total  int64                    `json:"total_outstanding"`
}

// ─── FinanceRepoExtended ──────────────────────────────────────────────────────

// FinanceRepoExtended menggabungkan finance.Repository dengan metode laporan tambahan.
type FinanceRepoExtended interface {
	finance.Repository
	SumByAkunPrefix(ctx context.Context, bmtID uuid.UUID, dari, sampai time.Time) (map[string]int64, error)
}

// ─── FinanceService ───────────────────────────────────────────────────────────

// FinanceService mengelola jurnal manual, transaksi operasional, dan laporan keuangan.
type FinanceService struct {
	repo             FinanceRepoExtended
	pembiayaanRepo   pembiayaan.Repository
	settingsResolver *settings.Resolver
}

func NewFinanceService(
	repo FinanceRepoExtended,
	pembiayaanRepo pembiayaan.Repository,
	settingsResolver *settings.Resolver,
) *FinanceService {
	return &FinanceService{
		repo:             repo,
		pembiayaanRepo:   pembiayaanRepo,
		settingsResolver: settingsResolver,
	}
}

// ─── Jurnal Manual ───────────────────────────────────────────────────────────

// CreateJurnal membuat jurnal manual baru (status DRAFT).
func (s *FinanceService) CreateJurnal(ctx context.Context, input finance.CreateJurnalInput) (*finance.JurnalManual, error) {
	j, err := finance.NewJurnalManual(input)
	if err != nil {
		return nil, fmt.Errorf("jurnal tidak valid: %w", err)
	}
	if err := s.repo.CreateJurnal(ctx, j); err != nil {
		return nil, fmt.Errorf("gagal simpan jurnal: %w", err)
	}
	return j, nil
}

// GetJurnal mengambil jurnal berdasarkan ID dengan validasi kepemilikan BMT.
func (s *FinanceService) GetJurnal(ctx context.Context, id, bmtID uuid.UUID) (*finance.JurnalManual, error) {
	j, err := s.repo.GetJurnalByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if j.BMTID != bmtID {
		return nil, finance.ErrJurnalNotFound
	}
	return j, nil
}

// ListJurnal mengembalikan daftar jurnal dengan filter.
func (s *FinanceService) ListJurnal(ctx context.Context, filter finance.ListJurnalFilter) ([]*finance.JurnalManual, int64, error) {
	return s.repo.ListJurnal(ctx, filter)
}

// PostJurnal mengubah status jurnal menjadi POSTED (disetujui oleh userID).
func (s *FinanceService) PostJurnal(ctx context.Context, id, bmtID, userID uuid.UUID) (*finance.JurnalManual, error) {
	j, err := s.GetJurnal(ctx, id, bmtID)
	if err != nil {
		return nil, err
	}
	if j.Status != finance.StatusJurnalDraft {
		return nil, fmt.Errorf("jurnal berstatus %s tidak bisa diposting", j.Status)
	}
	if err := s.repo.PostJurnal(ctx, id, userID); err != nil {
		return nil, fmt.Errorf("gagal post jurnal: %w", err)
	}
	j.Status = finance.StatusJurnalPosted
	j.DisetujuiOleh = &userID
	return j, nil
}

// ─── Transaksi Operasional ───────────────────────────────────────────────────

// CreateTransaksiOperasional membuat transaksi operasional baru.
func (s *FinanceService) CreateTransaksiOperasional(ctx context.Context, t *finance.TransaksiOperasional) (*finance.TransaksiOperasional, error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()
	if err := s.repo.CreateTransaksiOperasional(ctx, t); err != nil {
		return nil, fmt.Errorf("gagal simpan transaksi operasional: %w", err)
	}
	return t, nil
}

// ListTransaksiOperasional mengembalikan daftar transaksi operasional.
func (s *FinanceService) ListTransaksiOperasional(
	ctx context.Context,
	bmtID, cabangID uuid.UUID,
	dari, sampai time.Time,
	page, perPage int,
) ([]*finance.TransaksiOperasional, int64, error) {
	return s.repo.ListTransaksiOperasional(ctx, bmtID, cabangID, dari, sampai, page, perPage)
}

// ─── Laporan ─────────────────────────────────────────────────────────────────

// GetLaporanNeraca menghitung laporan neraca sederhana dari jurnal POSTED.
// Aset (prefix 1), Kewajiban (prefix 2), Ekuitas (prefix 3).
func (s *FinanceService) GetLaporanNeraca(ctx context.Context, bmtID uuid.UUID, dari, sampai time.Time) (*LaporanNeraca, error) {
	sums, err := s.repo.SumByAkunPrefix(ctx, bmtID, dari, sampai)
	if err != nil {
		return nil, fmt.Errorf("gagal hitung neraca: %w", err)
	}
	return &LaporanNeraca{
		Aset:      sums["1_DEBIT"] - sums["1_KREDIT"],
		Kewajiban: sums["2_KREDIT"] - sums["2_DEBIT"],
		Ekuitas:   sums["3_KREDIT"] - sums["3_DEBIT"],
		Dari:      dari.Format("2006-01-02"),
		Sampai:    sampai.Format("2006-01-02"),
	}, nil
}

// GetLaporanSHU menghitung Sisa Hasil Usaha dari jurnal POSTED.
// Pendapatan (prefix 4), Biaya (prefix 5).
func (s *FinanceService) GetLaporanSHU(ctx context.Context, bmtID uuid.UUID, dari, sampai time.Time) (*LaporanSHU, error) {
	sums, err := s.repo.SumByAkunPrefix(ctx, bmtID, dari, sampai)
	if err != nil {
		return nil, fmt.Errorf("gagal hitung SHU: %w", err)
	}
	pendapatan := sums["4_KREDIT"] - sums["4_DEBIT"]
	biaya := sums["5_DEBIT"] - sums["5_KREDIT"]
	return &LaporanSHU{
		Pendapatan: pendapatan,
		Biaya:      biaya,
		SHU:        pendapatan - biaya,
		Dari:       dari.Format("2006-01-02"),
		Sampai:     sampai.Format("2006-01-02"),
	}, nil
}

// GetLaporanArusKas menghitung arus kas dari jurnal POSTED.
// Akun Kas (prefix 1): debit = masuk, kredit = keluar.
func (s *FinanceService) GetLaporanArusKas(ctx context.Context, bmtID uuid.UUID, dari, sampai time.Time) (*LaporanArusKas, error) {
	sums, err := s.repo.SumByAkunPrefix(ctx, bmtID, dari, sampai)
	if err != nil {
		return nil, fmt.Errorf("gagal hitung arus kas: %w", err)
	}
	return &LaporanArusKas{
		KasMasuk:  sums["1_DEBIT"],
		KasKeluar: sums["1_KREDIT"],
		NetKas:    sums["1_DEBIT"] - sums["1_KREDIT"],
		Dari:      dari.Format("2006-01-02"),
		Sampai:    sampai.Format("2006-01-02"),
	}, nil
}

var keteranganKolektibilitas = map[int16]string{
	1: "Lancar",
	2: "Dalam Perhatian Khusus",
	3: "Kurang Lancar",
	4: "Diragukan",
	5: "Macet",
}

// GetLaporanKolektibilitas mengelompokkan pembiayaan aktif per level kolektibilitas OJK.
func (s *FinanceService) GetLaporanKolektibilitas(ctx context.Context, bmtID uuid.UUID) (*LaporanKolektibilitas, error) {
	list, err := s.pembiayaanRepo.ListAktifByBMT(ctx, bmtID)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil pembiayaan aktif: %w", err)
	}

	countByLevel := make(map[int16]int)
	outstandingByLevel := make(map[int16]int64)
	var totalOutstanding int64

	for _, p := range list {
		lvl := p.Kolektibilitas
		if lvl <= 0 {
			lvl = 1 // default lancar
		}
		countByLevel[lvl]++
		outstandingByLevel[lvl] += p.SaldoPokok
		totalOutstanding += p.SaldoPokok
	}

	data := make([]KolektibilitasRingkasan, 0, 5)
	for lvl := int16(1); lvl <= 5; lvl++ {
		data = append(data, KolektibilitasRingkasan{
			Level:       lvl,
			Keterangan:  keteranganKolektibilitas[lvl],
			Jumlah:      countByLevel[lvl],
			Outstanding: outstandingByLevel[lvl],
		})
	}

	return &LaporanKolektibilitas{
		BMTID: bmtID,
		Data:  data,
		Total: totalOutstanding,
	}, nil
}

// GetLaporanBagiHasilDeposito mengembalikan parameter distribusi bagi hasil deposito.
// Untuk laporan historis, jalankan distribusi_service.DistribusiBagiHasil terlebih dahulu.
func (s *FinanceService) GetLaporanBagiHasilDeposito(ctx context.Context, bmtID uuid.UUID, bulan time.Time) (map[string]interface{}, error) {
	rateStr := s.settingsResolver.ResolveWithDefault(ctx, bmtID, uuid.Nil, "DEPOSITO_RATE_BULANAN_PERSEN", "0.5")
	return map[string]interface{}{
		"bulan":       bulan.Format("2006-01"),
		"rate_persen": rateStr,
		"catatan":     "Laporan distribusi detail tersedia setelah eksekusi akhir bulan",
	}, nil
}
