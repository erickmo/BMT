package accounting

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrJurnalTidakBalance = errors.New("jurnal tidak balance: debit ≠ kredit")
	ErrEntryKosong        = errors.New("jurnal harus memiliki minimal 2 entry")
	ErrNominalNol         = errors.New("nominal entry tidak boleh nol")
)

type Posisi string

const (
	DEBIT  Posisi = "DEBIT"
	KREDIT Posisi = "KREDIT"
)

// Money adalah tipe untuk uang dalam rupiah (int64, bukan float)
type Money int64

func (m Money) Int64() int64 { return int64(m) }

// Entry adalah satu baris dalam jurnal (debit atau kredit)
type Entry struct {
	AkunKode   string `json:"akun_kode"`
	Posisi     Posisi `json:"posisi"`
	Nominal    Money  `json:"nominal"`
	Keterangan string `json:"keterangan,omitempty"`
}

// Journal adalah satu transaksi akuntansi double-entry
type Journal struct {
	BMTID         uuid.UUID  `json:"bmt_id"`
	CabangID      uuid.UUID  `json:"cabang_id"`
	Tanggal       time.Time  `json:"tanggal"`
	Keterangan    string     `json:"keterangan"`
	Referensi     string     `json:"referensi,omitempty"`
	ReferensiTipe string     `json:"referensi_tipe,omitempty"`
	Entries       []Entry    `json:"entries"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
}

// JournalRecord adalah journal yang sudah disimpan (memiliki ID dan nomor)
type JournalRecord struct {
	ID          uuid.UUID  `json:"id"`
	NomorJurnal string     `json:"nomor_jurnal"`
	Journal                // embed
	TotalDebit  Money      `json:"total_debit"`
	TotalKredit Money      `json:"total_kredit"`
	IsBalanced  bool       `json:"is_balanced"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Repository adalah interface untuk menyimpan jurnal ke DB
type Repository interface {
	// SaveJournal menyimpan jurnal dan semua entries-nya dalam satu transaksi DB
	SaveJournal(ctx context.Context, record *JournalRecord) error
	// GenerateNomor menghasilkan nomor jurnal unik per BMT
	GenerateNomor(ctx context.Context, bmtID uuid.UUID) (string, error)
	// GetByID mengambil jurnal berdasarkan ID
	GetByID(ctx context.Context, id uuid.UUID) (*JournalRecord, error)
	// ListByPeriode mengambil jurnal berdasarkan periode
	ListByPeriode(ctx context.Context, bmtID, cabangID uuid.UUID, dari, sampai time.Time) ([]*JournalRecord, error)
	// GetSaldoAkun mengembalikan saldo akun pada periode tertentu
	GetSaldoAkun(ctx context.Context, bmtID uuid.UUID, akunKode string, sampai time.Time) (Money, error)
}

// Poster adalah interface utama untuk posting jurnal
type Poster interface {
	Post(ctx context.Context, j Journal) (*JournalRecord, error)
}

// Service implementasi Poster
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Post memvalidasi dan menyimpan jurnal double-entry
func (s *Service) Post(ctx context.Context, j Journal) (*JournalRecord, error) {
	if err := Validate(j); err != nil {
		return nil, err
	}

	var totalDebit, totalKredit Money
	for _, e := range j.Entries {
		if e.Posisi == DEBIT {
			totalDebit += e.Nominal
		} else {
			totalKredit += e.Nominal
		}
	}

	nomor, err := s.repo.GenerateNomor(ctx, j.BMTID)
	if err != nil {
		return nil, fmt.Errorf("gagal generate nomor jurnal: %w", err)
	}

	record := &JournalRecord{
		ID:          uuid.New(),
		NomorJurnal: nomor,
		Journal:     j,
		TotalDebit:  totalDebit,
		TotalKredit: totalKredit,
		IsBalanced:  totalDebit == totalKredit,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.SaveJournal(ctx, record); err != nil {
		return nil, fmt.Errorf("gagal simpan jurnal: %w", err)
	}

	return record, nil
}

// Validate memvalidasi jurnal sebelum disimpan
func Validate(j Journal) error {
	if len(j.Entries) < 2 {
		return ErrEntryKosong
	}

	var totalDebit, totalKredit Money
	for _, e := range j.Entries {
		if e.Nominal <= 0 {
			return ErrNominalNol
		}
		if e.Posisi == DEBIT {
			totalDebit += e.Nominal
		} else {
			totalKredit += e.Nominal
		}
	}

	if totalDebit != totalKredit {
		return fmt.Errorf("%w: debit=%d kredit=%d", ErrJurnalTidakBalance, totalDebit, totalKredit)
	}

	return nil
}

// ChartOfAccounts akun standar BMT syariah
var ChartOfAccounts = map[string]string{
	// ASET
	"101": "Kas",
	"102": "Kas Bank",
	"111": "Piutang Murabahah",
	"112": "Piutang Musyarakah",
	"113": "Piutang Mudharabah",
	"114": "Piutang Ijarah",
	// KEWAJIBAN
	"201": "Simpanan Pokok",
	"202": "Simpanan Wajib",
	"203": "Simpanan Sukarela",
	"204": "Deposito Mudharabah",
	"211": "Dana Sosial (Ta'zir/Infaq)",
	// EKUITAS
	"301": "Modal Disetor",
	"302": "Cadangan",
	"303": "Cadangan Umum",
	"304": "SHU Ditahan",
	// PENDAPATAN
	"401": "Pendapatan Margin Murabahah",
	"402": "Pendapatan Bagi Hasil Mudharabah",
	"403": "Pendapatan Bagi Hasil Musyarakah",
	"404": "Pendapatan Ujrah/Ijarah",
	"405": "Biaya Admin Rekening",
	"406": "Komisi NFC",
	"407": "Komisi OPOP",
	// BEBAN
	"501": "Beban Bagi Hasil",
	"502": "Beban Operasional",
	"503": "Penyisihan Penghapusan Aktiva",
	"504": "Beban Gaji",
	"505": "Beban Utilitas",
	// DANA SOSIAL
	"601": "Penerimaan Donasi",
	"602": "Penerimaan Wakaf",
	"603": "Penerimaan Infaq/Shadaqah",
	"604": "Penerimaan Zakat",
	"611": "Penyaluran Donasi",
	"612": "Penyaluran Wakaf",
	"613": "Penyaluran Infaq/Shadaqah",
	"614": "Penyaluran Zakat",
}
