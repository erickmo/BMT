package accounting_test

import (
	"context"
	"testing"
	"time"

	accounting "github.com/bmt-saas/module-vernon-accounting"
	"github.com/google/uuid"
)

// MockRepo implements accounting.Repository for testing
type MockRepo struct {
	journals []*accounting.JournalRecord
	counter  int
}

func (m *MockRepo) SaveJournal(ctx context.Context, record *accounting.JournalRecord) error {
	m.journals = append(m.journals, record)
	return nil
}

func (m *MockRepo) GenerateNomor(ctx context.Context, bmtID uuid.UUID) (string, error) {
	m.counter++
	return "JRN-TEST-" + string(rune('0'+m.counter)), nil
}

func (m *MockRepo) GetByID(ctx context.Context, id uuid.UUID) (*accounting.JournalRecord, error) {
	for _, j := range m.journals {
		if j.ID == id {
			return j, nil
		}
	}
	return nil, nil
}

func (m *MockRepo) ListByPeriode(ctx context.Context, bmtID, cabangID uuid.UUID, dari, sampai time.Time) ([]*accounting.JournalRecord, error) {
	return m.journals, nil
}

func (m *MockRepo) GetSaldoAkun(ctx context.Context, bmtID uuid.UUID, akunKode string, sampai time.Time) (accounting.Money, error) {
	return 0, nil
}

func TestJurnal_Balance_Valid(t *testing.T) {
	repo := &MockRepo{}
	svc := accounting.NewService(repo)

	record, err := svc.Post(context.Background(), accounting.Journal{
		BMTID:      uuid.New(),
		CabangID:   uuid.New(),
		Tanggal:    time.Now(),
		Keterangan: "Setoran tunai nasabah",
		Entries: []accounting.Entry{
			{AkunKode: "101", Posisi: accounting.DEBIT, Nominal: 500000},
			{AkunKode: "202", Posisi: accounting.KREDIT, Nominal: 500000},
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !record.IsBalanced {
		t.Error("jurnal seharusnya balance")
	}
	if record.TotalDebit != record.TotalKredit {
		t.Errorf("total debit %d != total kredit %d", record.TotalDebit, record.TotalKredit)
	}
}

func TestJurnal_TidakBalance_Error(t *testing.T) {
	repo := &MockRepo{}
	svc := accounting.NewService(repo)

	_, err := svc.Post(context.Background(), accounting.Journal{
		BMTID:    uuid.New(),
		CabangID: uuid.New(),
		Tanggal:  time.Now(),
		Keterangan: "Jurnal tidak balance",
		Entries: []accounting.Entry{
			{AkunKode: "101", Posisi: accounting.DEBIT, Nominal: 500000},
			{AkunKode: "202", Posisi: accounting.KREDIT, Nominal: 400000}, // tidak balance
		},
	})

	if err == nil {
		t.Fatal("seharusnya ada error")
	}
}

func TestJurnal_EntryKosong_Error(t *testing.T) {
	repo := &MockRepo{}
	svc := accounting.NewService(repo)

	_, err := svc.Post(context.Background(), accounting.Journal{
		BMTID:    uuid.New(),
		CabangID: uuid.New(),
		Tanggal:  time.Now(),
		Keterangan: "Jurnal tanpa entry",
		Entries: []accounting.Entry{},
	})

	if err == nil {
		t.Fatal("seharusnya ada error")
	}
}

func TestJurnal_NominalNol_Error(t *testing.T) {
	repo := &MockRepo{}
	svc := accounting.NewService(repo)

	_, err := svc.Post(context.Background(), accounting.Journal{
		BMTID:    uuid.New(),
		CabangID: uuid.New(),
		Tanggal:  time.Now(),
		Keterangan: "Jurnal nominal nol",
		Entries: []accounting.Entry{
			{AkunKode: "101", Posisi: accounting.DEBIT, Nominal: 0},
			{AkunKode: "202", Posisi: accounting.KREDIT, Nominal: 0},
		},
	})

	if err == nil {
		t.Fatal("seharusnya ada error")
	}
}

func TestJurnal_MultiEntry_Balance(t *testing.T) {
	repo := &MockRepo{}
	svc := accounting.NewService(repo)

	// 3 debit, 1 kredit - total harus sama
	record, err := svc.Post(context.Background(), accounting.Journal{
		BMTID:    uuid.New(),
		CabangID: uuid.New(),
		Tanggal:  time.Now(),
		Keterangan: "Angsuran pembiayaan murabahah",
		Entries: []accounting.Entry{
			{AkunKode: "101", Posisi: accounting.DEBIT, Nominal: 1100000},
			{AkunKode: "111", Posisi: accounting.KREDIT, Nominal: 1000000},
			{AkunKode: "401", Posisi: accounting.KREDIT, Nominal: 100000},
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !record.IsBalanced {
		t.Error("jurnal seharusnya balance")
	}
}

func TestValidate_DoubleEntry_SemuaTransaksi(t *testing.T) {
	tests := []struct {
		nama    string
		entries []accounting.Entry
		wantErr bool
	}{
		{
			nama: "setoran_simpanan",
			entries: []accounting.Entry{
				{AkunKode: "101", Posisi: accounting.DEBIT, Nominal: 100000},
				{AkunKode: "203", Posisi: accounting.KREDIT, Nominal: 100000},
			},
			wantErr: false,
		},
		{
			nama: "penarikan_simpanan",
			entries: []accounting.Entry{
				{AkunKode: "203", Posisi: accounting.DEBIT, Nominal: 50000},
				{AkunKode: "101", Posisi: accounting.KREDIT, Nominal: 50000},
			},
			wantErr: false,
		},
		{
			nama: "pencairan_pembiayaan",
			entries: []accounting.Entry{
				{AkunKode: "111", Posisi: accounting.DEBIT, Nominal: 5000000},
				{AkunKode: "101", Posisi: accounting.KREDIT, Nominal: 5000000},
			},
			wantErr: false,
		},
		{
			nama: "tidak_balance",
			entries: []accounting.Entry{
				{AkunKode: "101", Posisi: accounting.DEBIT, Nominal: 100000},
				{AkunKode: "203", Posisi: accounting.KREDIT, Nominal: 99000},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			j := accounting.Journal{
				BMTID:    uuid.New(),
				CabangID: uuid.New(),
				Tanggal:  time.Now(),
				Keterangan: tt.nama,
				Entries: tt.entries,
			}
			err := accounting.Validate(j)
			if tt.wantErr && err == nil {
				t.Error("seharusnya ada error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("tidak seharusnya ada error: %v", err)
			}
		})
	}
}
