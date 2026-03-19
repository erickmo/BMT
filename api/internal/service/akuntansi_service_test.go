package service_test

import (
	"context"
	"testing"

	"github.com/bmt-saas/api/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAkuntansiRepo struct {
	mock.Mock
}

func (m *MockAkuntansiRepo) CreateJurnal(ctx context.Context, j *service.Jurnal) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MockAkuntansiRepo) CreateEntries(ctx context.Context, entries []*service.EntryJurnal) error {
	args := m.Called(ctx, entries)
	return args.Error(0)
}

func (m *MockAkuntansiRepo) GenerateNomorJurnal(ctx context.Context, bmtID uuid.UUID) (string, error) {
	args := m.Called(ctx, bmtID)
	return args.String(0), args.Error(1)
}

func TestAkuntansi_PostJurnal_Balance(t *testing.T) {
	repo := new(MockAkuntansiRepo)
	svc := service.NewAkuntansiService(repo)

	bmtID := uuid.New()
	repo.On("GenerateNomorJurnal", mock.Anything, bmtID).Return("JRN-2025-0001", nil)
	repo.On("CreateJurnal", mock.Anything, mock.Anything).Return(nil)
	repo.On("CreateEntries", mock.Anything, mock.Anything).Return(nil)

	err := svc.PostJurnal(context.Background(), service.PostJurnalInput{
		BMTID:      bmtID,
		CabangID:   uuid.New(),
		Keterangan: "Setoran tunai",
		Entries: []service.JurnalEntry{
			{KodeAkun: "101", Posisi: "DEBIT", Nominal: 1000000},
			{KodeAkun: "202", Posisi: "KREDIT", Nominal: 1000000},
		},
	})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAkuntansi_PostJurnal_TidakBalance(t *testing.T) {
	repo := new(MockAkuntansiRepo)
	svc := service.NewAkuntansiService(repo)

	err := svc.PostJurnal(context.Background(), service.PostJurnalInput{
		BMTID:      uuid.New(),
		CabangID:   uuid.New(),
		Keterangan: "Test tidak balance",
		Entries: []service.JurnalEntry{
			{KodeAkun: "101", Posisi: "DEBIT", Nominal: 1000000},
			{KodeAkun: "202", Posisi: "KREDIT", Nominal: 900000}, // tidak balance
		},
	})

	assert.ErrorIs(t, err, service.ErrJurnalTidakBalance)
}

// TestJurnal_SemuaTransaksi_DoubleEntryBalance adalah test wajib per CLAUDE.md.
// Memastikan SETIAP transaksi keuangan mengikuti prinsip double-entry akuntansi
// (total debit HARUS selalu sama dengan total kredit).
func TestJurnal_SemuaTransaksi_DoubleEntryBalance(t *testing.T) {
	skenarioTransaksi := []struct {
		nama    string
		entries []service.JurnalEntry
		balance bool
	}{
		{
			nama: "setoran tunai: kas debit, simpanan sukarela kredit",
			entries: []service.JurnalEntry{
				{KodeAkun: "101", Posisi: "DEBIT", Nominal: 500000},
				{KodeAkun: "202", Posisi: "KREDIT", Nominal: 500000},
			},
			balance: true,
		},
		{
			nama: "penarikan tunai: simpanan debit, kas kredit",
			entries: []service.JurnalEntry{
				{KodeAkun: "202", Posisi: "DEBIT", Nominal: 200000},
				{KodeAkun: "101", Posisi: "KREDIT", Nominal: 200000},
			},
			balance: true,
		},
		{
			nama: "angsuran pembiayaan: kas debit, piutang kredit, margin pendapatan kredit",
			entries: []service.JurnalEntry{
				{KodeAkun: "101", Posisi: "DEBIT", Nominal: 1100000},
				{KodeAkun: "111", Posisi: "KREDIT", Nominal: 1000000},
				{KodeAkun: "401", Posisi: "KREDIT", Nominal: 100000},
			},
			balance: true,
		},
		{
			nama: "transaksi tidak balance: harus ditolak",
			entries: []service.JurnalEntry{
				{KodeAkun: "101", Posisi: "DEBIT", Nominal: 500000},
				{KodeAkun: "202", Posisi: "KREDIT", Nominal: 400000}, // tidak balance!
			},
			balance: false,
		},
		{
			nama: "biaya admin rekening: debit simpanan, kredit pendapatan admin",
			entries: []service.JurnalEntry{
				{KodeAkun: "202", Posisi: "DEBIT", Nominal: 10000},
				{KodeAkun: "405", Posisi: "KREDIT", Nominal: 10000},
			},
			balance: true,
		},
		{
			nama: "transfer antar rekening: debit rek asal, kredit rek tujuan",
			entries: []service.JurnalEntry{
				{KodeAkun: "202", Posisi: "DEBIT", Nominal: 300000},
				{KodeAkun: "202", Posisi: "KREDIT", Nominal: 300000},
			},
			balance: true,
		},
	}

	for _, sk := range skenarioTransaksi {
		t.Run(sk.nama, func(t *testing.T) {
			repo := new(MockAkuntansiRepo)
			svc := service.NewAkuntansiService(repo)

			if sk.balance {
				// Siapkan mock untuk transaksi yang balance
				repo.On("GenerateNomorJurnal", mock.Anything, mock.Anything).Return("JRN-TEST", nil)
				repo.On("CreateJurnal", mock.Anything, mock.Anything).Return(nil)
				repo.On("CreateEntries", mock.Anything, mock.Anything).Return(nil)
			}

			err := svc.PostJurnal(context.Background(), service.PostJurnalInput{
				BMTID:      uuid.New(),
				CabangID:   uuid.New(),
				Keterangan: sk.nama,
				Entries:    sk.entries,
			})

			if sk.balance {
				assert.NoError(t, err, "transaksi balance harus diterima")
				repo.AssertExpectations(t)
			} else {
				assert.ErrorIs(t, err, service.ErrJurnalTidakBalance,
					"transaksi tidak balance harus ditolak dengan ErrJurnalTidakBalance")
			}
		})
	}
}
