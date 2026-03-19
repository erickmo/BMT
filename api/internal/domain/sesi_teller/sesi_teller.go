package sesi_teller

import (
	"context"
	"errors"
	"time"

	"github.com/bmt-saas/api/pkg/money"
	"github.com/google/uuid"
)

var (
	ErrSesiSudahAktif = errors.New("teller sudah memiliki sesi aktif hari ini")
	ErrSesiTidakAktif = errors.New("tidak ada sesi teller aktif")
	ErrSesiSelisih    = errors.New("saldo fisik tidak sesuai, sesi ditolak")
	ErrSesiSudahTutup = errors.New("sesi teller sudah ditutup")
)

type StatusSesi string

const (
	StatusAktif   StatusSesi = "AKTIF"
	StatusTutup   StatusSesi = "TUTUP"
	StatusSelisih StatusSesi = "SELISIH"
)

type ItemPecahan struct {
	PecahanID uuid.UUID `json:"pecahan_id"`
	Nominal   int64     `json:"nominal"`
	Jenis     string    `json:"jenis"`
	Label     string    `json:"label"`
	Jumlah    int       `json:"jumlah"`
	Subtotal  int64     `json:"subtotal"`
}

type SesiTeller struct {
	ID                uuid.UUID     `json:"id"`
	BMTID             uuid.UUID     `json:"bmt_id"`
	CabangID          uuid.UUID     `json:"cabang_id"`
	TellerID          uuid.UUID     `json:"teller_id"`
	Tanggal           time.Time     `json:"tanggal"`
	SaldoAwal         money.Money   `json:"saldo_awal"`
	Redenominasi      []ItemPecahan `json:"redenominasi"`
	SaldoAkhir        *money.Money  `json:"saldo_akhir,omitempty"`
	RedenominasiAkhir []ItemPecahan `json:"redenominasi_akhir,omitempty"`
	Status            StatusSesi    `json:"status"`
	ToleransiSelisih  money.Money   `json:"toleransi_selisih"`
	Selisih           *money.Money  `json:"selisih,omitempty"`
	DibukaPada        time.Time     `json:"dibuka_pada"`
	DitutupPada       *time.Time    `json:"ditutup_pada,omitempty"`
}

func (s *SesiTeller) HitungSaldoAwal() money.Money {
	total := money.Zero
	for _, item := range s.Redenominasi {
		total = total.Add(money.New(item.Subtotal))
	}
	return total
}

func (s *SesiTeller) TutupSesi(redenominasiAkhir []ItemPecahan, toleransiSelisih money.Money) error {
	if s.Status != StatusAktif {
		return ErrSesiSudahTutup
	}

	saldoFisik := money.Zero
	for _, item := range redenominasiAkhir {
		saldoFisik = saldoFisik.Add(money.New(item.Subtotal))
	}

	selisih := s.SaldoAwal.Sub(saldoFisik)
	if selisih < 0 {
		selisih = -selisih
	}

	if selisih > toleransiSelisih {
		s.Status = StatusSelisih
		s.Selisih = &selisih
		s.RedenominasiAkhir = redenominasiAkhir
		return ErrSesiSelisih
	}

	now := time.Now()
	s.SaldoAkhir = &saldoFisik
	s.RedenominasiAkhir = redenominasiAkhir
	s.Status = StatusTutup
	s.Selisih = &selisih
	s.DitutupPada = &now
	return nil
}

type Repository interface {
	Create(ctx context.Context, s *SesiTeller) error
	GetAktifByTeller(ctx context.Context, tellerID uuid.UUID) (*SesiTeller, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SesiTeller, error)
	Update(ctx context.Context, s *SesiTeller) error
	ListByBMT(ctx context.Context, bmtID, cabangID uuid.UUID, tanggal time.Time) ([]*SesiTeller, error)
}
