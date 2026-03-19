package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrJurnalTidakBalance = errors.New("jurnal tidak balance: debit ≠ kredit")
)

// JurnalEntry adalah satu baris entri jurnal double-entry.
type JurnalEntry struct {
	KodeAkun string
	Posisi   string // "DEBIT" | "KREDIT"
	Nominal  int64
}

// PostJurnalInput adalah input untuk memposting jurnal.
type PostJurnalInput struct {
	BMTID      uuid.UUID
	CabangID   uuid.UUID
	Keterangan string
	Referensi  string // UUID string referensi transaksi (opsional)
	Entries    []JurnalEntry
}

// Jurnal adalah representasi jurnal akuntansi yang tersimpan.
type Jurnal struct {
	ID         uuid.UUID
	BMTID      uuid.UUID
	CabangID   uuid.UUID
	Nomor      string
	Keterangan string
}

// EntryJurnal adalah representasi baris jurnal yang tersimpan.
type EntryJurnal struct {
	JurnalID uuid.UUID
	KodeAkun string
	Posisi   string
	Nominal  int64
}

// AkuntansiRepository adalah kontrak repository untuk akuntansi.
type AkuntansiRepository interface {
	CreateJurnal(ctx context.Context, j *Jurnal) error
	CreateEntries(ctx context.Context, entries []*EntryJurnal) error
	GenerateNomorJurnal(ctx context.Context, bmtID uuid.UUID) (string, error)
}

// AkuntansiService mengelola posting jurnal double-entry.
type AkuntansiService struct {
	repo AkuntansiRepository
}

// NewAkuntansiService membuat instance baru AkuntansiService.
func NewAkuntansiService(repo AkuntansiRepository) *AkuntansiService {
	return &AkuntansiService{repo: repo}
}

// PostJurnal memvalidasi balance debit=kredit lalu menyimpan jurnal.
func (s *AkuntansiService) PostJurnal(ctx context.Context, input PostJurnalInput) error {
	var totalDebit, totalKredit int64
	for _, e := range input.Entries {
		switch e.Posisi {
		case "DEBIT":
			totalDebit += e.Nominal
		case "KREDIT":
			totalKredit += e.Nominal
		}
	}
	if totalDebit != totalKredit {
		return ErrJurnalTidakBalance
	}

	nomor, err := s.repo.GenerateNomorJurnal(ctx, input.BMTID)
	if err != nil {
		return err
	}

	j := &Jurnal{
		ID:         uuid.New(),
		BMTID:      input.BMTID,
		CabangID:   input.CabangID,
		Nomor:      nomor,
		Keterangan: input.Keterangan,
	}
	if err := s.repo.CreateJurnal(ctx, j); err != nil {
		return err
	}

	var entries []*EntryJurnal
	for _, e := range input.Entries {
		entries = append(entries, &EntryJurnal{
			JurnalID: j.ID,
			KodeAkun: e.KodeAkun,
			Posisi:   e.Posisi,
			Nominal:  e.Nominal,
		})
	}
	return s.repo.CreateEntries(ctx, entries)
}
