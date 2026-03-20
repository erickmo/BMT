package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/bmt-saas/api/internal/domain/autodebet"
	"github.com/bmt-saas/api/internal/domain/rekening"
	"github.com/bmt-saas/api/internal/service"
	"github.com/bmt-saas/api/pkg/money"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ── Mock: rekening.Repository ────────────────────────────────────────────────

type MockRekeningRepo struct {
	mock.Mock
}

func (m *MockRekeningRepo) CreateJenis(ctx context.Context, jr *rekening.JenisRekening) error {
	args := m.Called(ctx, jr)
	return args.Error(0)
}

func (m *MockRekeningRepo) GetJenisByID(ctx context.Context, id uuid.UUID) (*rekening.JenisRekening, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rekening.JenisRekening), args.Error(1)
}

func (m *MockRekeningRepo) ListJenisByBMT(ctx context.Context, bmtID uuid.UUID) ([]*rekening.JenisRekening, error) {
	args := m.Called(ctx, bmtID)
	return args.Get(0).([]*rekening.JenisRekening), args.Error(1)
}

func (m *MockRekeningRepo) UpdateJenis(ctx context.Context, jr *rekening.JenisRekening) error {
	args := m.Called(ctx, jr)
	return args.Error(0)
}

func (m *MockRekeningRepo) Create(ctx context.Context, r *rekening.Rekening) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockRekeningRepo) GetByID(ctx context.Context, id uuid.UUID) (*rekening.Rekening, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rekening.Rekening), args.Error(1)
}

func (m *MockRekeningRepo) GetByNomor(ctx context.Context, nomor string) (*rekening.Rekening, error) {
	args := m.Called(ctx, nomor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rekening.Rekening), args.Error(1)
}

func (m *MockRekeningRepo) ListByNasabah(ctx context.Context, nasabahID uuid.UUID) ([]*rekening.Rekening, error) {
	args := m.Called(ctx, nasabahID)
	return args.Get(0).([]*rekening.Rekening), args.Error(1)
}

func (m *MockRekeningRepo) ListByBMT(ctx context.Context, bmtID, cabangID uuid.UUID, page, perPage int) ([]*rekening.Rekening, int64, error) {
	args := m.Called(ctx, bmtID, cabangID, page, perPage)
	return args.Get(0).([]*rekening.Rekening), args.Get(1).(int64), args.Error(2)
}

func (m *MockRekeningRepo) UpdateSaldo(ctx context.Context, id uuid.UUID, saldoBaru int64) error {
	args := m.Called(ctx, id, saldoBaru)
	return args.Error(0)
}

func (m *MockRekeningRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status rekening.StatusRekening, alasan string) error {
	args := m.Called(ctx, id, status, alasan)
	return args.Error(0)
}

func (m *MockRekeningRepo) LockForUpdate(ctx context.Context, id uuid.UUID) (*rekening.Rekening, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rekening.Rekening), args.Error(1)
}

func (m *MockRekeningRepo) CreateTransaksi(ctx context.Context, t *rekening.TransaksiRekening) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockRekeningRepo) ListTransaksi(ctx context.Context, rekeningID uuid.UUID, limit, offset int) ([]*rekening.TransaksiRekening, int64, error) {
	args := m.Called(ctx, rekeningID, limit, offset)
	return args.Get(0).([]*rekening.TransaksiRekening), args.Get(1).(int64), args.Error(2)
}

func (m *MockRekeningRepo) GetTransaksiByIdempotency(ctx context.Context, key uuid.UUID) (*rekening.TransaksiRekening, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rekening.TransaksiRekening), args.Error(1)
}

func (m *MockRekeningRepo) GenerateNomorRekening(ctx context.Context, bmtID, cabangID uuid.UUID, kodeJenis string) (string, error) {
	args := m.Called(ctx, bmtID, cabangID, kodeJenis)
	return args.String(0), args.Error(1)
}

func (m *MockRekeningRepo) ListDepositoAktif(ctx context.Context, bmtID uuid.UUID) ([]*rekening.Rekening, error) {
	args := m.Called(ctx, bmtID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*rekening.Rekening), args.Error(1)
}

// ── Mock: autodebet.Repository ────────────────────────────────────────────────

type MockAutodebetRepo struct {
	mock.Mock
}

func (m *MockAutodebetRepo) CreateConfig(ctx context.Context, c *autodebet.Config) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockAutodebetRepo) GetConfig(ctx context.Context, id uuid.UUID) (*autodebet.Config, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*autodebet.Config), args.Error(1)
}

func (m *MockAutodebetRepo) ListConfigByRekening(ctx context.Context, rekeningID uuid.UUID) ([]*autodebet.Config, error) {
	args := m.Called(ctx, rekeningID)
	return args.Get(0).([]*autodebet.Config), args.Error(1)
}

func (m *MockAutodebetRepo) ListConfigAktifByBMT(ctx context.Context, bmtID uuid.UUID) ([]*autodebet.Config, error) {
	args := m.Called(ctx, bmtID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*autodebet.Config), args.Error(1)
}

func (m *MockAutodebetRepo) UpdateConfig(ctx context.Context, c *autodebet.Config) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockAutodebetRepo) CreateJadwal(ctx context.Context, j *autodebet.Jadwal) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *MockAutodebetRepo) ListJadwalByTanggal(ctx context.Context, bmtID uuid.UUID, tanggal time.Time) ([]*autodebet.Jadwal, error) {
	args := m.Called(ctx, bmtID, tanggal)
	return args.Get(0).([]*autodebet.Jadwal), args.Error(1)
}

func (m *MockAutodebetRepo) UpdateJadwalStatus(ctx context.Context, id uuid.UUID, status autodebet.StatusJadwal) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockAutodebetRepo) CreateTunggakan(ctx context.Context, t *autodebet.Tunggakan) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockAutodebetRepo) ListTunggakanByRekening(ctx context.Context, rekeningID uuid.UUID) ([]*autodebet.Tunggakan, error) {
	args := m.Called(ctx, rekeningID)
	return args.Get(0).([]*autodebet.Tunggakan), args.Error(1)
}

func (m *MockAutodebetRepo) UpdateTunggakan(ctx context.Context, t *autodebet.Tunggakan) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

// ── Mock: JurnalService ──────────────────────────────────────────────────────

type MockJurnalService struct {
	mock.Mock
}

func (m *MockJurnalService) PostJurnal(ctx context.Context, input service.PostJurnalInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

// ── Helper ───────────────────────────────────────────────────────────────────

func buatRekeningAktif(saldo int64) *rekening.Rekening {
	return &rekening.Rekening{
		ID:              uuid.New(),
		BMTID:           uuid.New(),
		CabangID:        uuid.New(),
		NasabahID:       uuid.New(),
		JenisRekeningID: uuid.New(),
		NomorRekening:   "001-001-000001",
		Saldo:           money.New(saldo),
		Status:          rekening.StatusAktif,
		TanggalBuka:     time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func buatJenisSimpananSukarela() *rekening.JenisRekening {
	return &rekening.JenisRekening{
		ID:          uuid.New(),
		SetoranMin:  10000,
		BisaDitarik: true,
	}
}

func buatService(repo *MockRekeningRepo, autodebetRepo *MockAutodebetRepo, jurnal *MockJurnalService) *service.RekeningService {
	return service.NewRekeningService(repo, autodebetRepo, nil, jurnal)
}

// ── Tests: Setor ─────────────────────────────────────────────────────────────

func TestRekeningService_Setor_Berhasil(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	rek := buatRekeningAktif(100000)
	jenis := buatJenisSimpananSukarela()

	repo.On("GetTransaksiByIdempotency", mock.Anything, mock.Anything).Return(nil, rekening.ErrRekeningNotFound)
	repo.On("LockForUpdate", mock.Anything, rek.ID).Return(rek, nil)
	repo.On("GetJenisByID", mock.Anything, rek.JenisRekeningID).Return(jenis, nil)
	repo.On("UpdateSaldo", mock.Anything, rek.ID, mock.Anything).Return(nil)
	repo.On("CreateTransaksi", mock.Anything, mock.Anything).Return(nil)
	jurnal.On("PostJurnal", mock.Anything, mock.Anything).Return(nil)

	idemKey := uuid.New()
	result, err := svc.Setor(context.Background(), rekening.SetoranInput{
		RekeningID:     rek.ID,
		Nominal:        50000,
		Keterangan:     "Setoran tunai",
		IdempotencyKey: &idemKey,
		CreatedBy:      uuid.New(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "SETOR", result.Jenis)
	assert.Equal(t, int64(50000), result.Nominal)
	assert.Equal(t, int64(100000), result.SaldoSebelum)
	assert.Equal(t, int64(150000), result.SaldoSesudah)
	repo.AssertExpectations(t)
}

func TestRekeningService_Setor_Idempotency_SudahAda(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	idemKey := uuid.New()
	existingTr := &rekening.TransaksiRekening{
		ID:     uuid.New(),
		Jenis:  "SETOR",
		Nominal: 50000,
	}
	repo.On("GetTransaksiByIdempotency", mock.Anything, idemKey).Return(existingTr, nil)

	result, err := svc.Setor(context.Background(), rekening.SetoranInput{
		RekeningID:     uuid.New(),
		Nominal:        50000,
		IdempotencyKey: &idemKey,
	})

	assert.NoError(t, err)
	assert.Equal(t, existingTr.ID, result.ID)
	// Pastikan LockForUpdate tidak dipanggil karena idempotency return early
	repo.AssertNotCalled(t, "LockForUpdate", mock.Anything, mock.Anything)
}

func TestRekeningService_Setor_RekeningBeku(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	rek := buatRekeningAktif(100000)
	rek.Status = rekening.StatusBlokir
	jenis := buatJenisSimpananSukarela()

	repo.On("LockForUpdate", mock.Anything, rek.ID).Return(rek, nil)
	repo.On("GetJenisByID", mock.Anything, rek.JenisRekeningID).Return(jenis, nil)

	_, err := svc.Setor(context.Background(), rekening.SetoranInput{
		RekeningID: rek.ID,
		Nominal:    50000,
	})

	assert.ErrorIs(t, err, rekening.ErrRekeningBeku)
}

func TestRekeningService_Setor_NominalDibawahMinimum(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	rek := buatRekeningAktif(100000)
	jenis := buatJenisSimpananSukarela()
	jenis.SetoranMin = 50000

	repo.On("LockForUpdate", mock.Anything, rek.ID).Return(rek, nil)
	repo.On("GetJenisByID", mock.Anything, rek.JenisRekeningID).Return(jenis, nil)

	_, err := svc.Setor(context.Background(), rekening.SetoranInput{
		RekeningID: rek.ID,
		Nominal:    10000, // di bawah minimum 50000
	})

	assert.ErrorIs(t, err, rekening.ErrSetoranDibawahMin)
}

// ── Tests: Tarik ─────────────────────────────────────────────────────────────

func TestRekeningService_Tarik_Berhasil(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	rek := buatRekeningAktif(500000)
	jenis := buatJenisSimpananSukarela()

	repo.On("LockForUpdate", mock.Anything, rek.ID).Return(rek, nil)
	repo.On("GetJenisByID", mock.Anything, rek.JenisRekeningID).Return(jenis, nil)
	repo.On("UpdateSaldo", mock.Anything, rek.ID, mock.Anything).Return(nil)
	repo.On("CreateTransaksi", mock.Anything, mock.Anything).Return(nil)
	jurnal.On("PostJurnal", mock.Anything, mock.Anything).Return(nil)

	result, err := svc.Tarik(context.Background(), rekening.PenarikanInput{
		RekeningID: rek.ID,
		Nominal:    200000,
		Keterangan: "Penarikan tunai",
		CreatedBy:  uuid.New(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TARIK", result.Jenis)
	assert.Equal(t, int64(200000), result.Nominal)
	assert.Equal(t, int64(500000), result.SaldoSebelum)
	assert.Equal(t, int64(300000), result.SaldoSesudah)
	repo.AssertExpectations(t)
}

func TestRekeningService_Tarik_SaldoTidakCukup(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	rek := buatRekeningAktif(50000)
	jenis := buatJenisSimpananSukarela()

	repo.On("LockForUpdate", mock.Anything, rek.ID).Return(rek, nil)
	repo.On("GetJenisByID", mock.Anything, rek.JenisRekeningID).Return(jenis, nil)

	_, err := svc.Tarik(context.Background(), rekening.PenarikanInput{
		RekeningID: rek.ID,
		Nominal:    200000, // lebih dari saldo 50000
	})

	assert.ErrorIs(t, err, rekening.ErrSaldoTidakCukup)
}

func TestRekeningService_Tarik_JenisTidakBisaDitarik(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	rek := buatRekeningAktif(500000)
	jenis := buatJenisSimpananSukarela()
	jenis.BisaDitarik = false // deposito tidak bisa ditarik

	repo.On("LockForUpdate", mock.Anything, rek.ID).Return(rek, nil)
	repo.On("GetJenisByID", mock.Anything, rek.JenisRekeningID).Return(jenis, nil)

	_, err := svc.Tarik(context.Background(), rekening.PenarikanInput{
		RekeningID: rek.ID,
		Nominal:    100000,
	})

	assert.ErrorIs(t, err, rekening.ErrPenarikanTidakBisa)
}

// ── Tests: EksekusiAutodebetJadwal ───────────────────────────────────────────

// TestAutodebet_SaldoKurang_PartialDebitDanTunggakan adalah test wajib per CLAUDE.md.
// Memastikan autodebet partial: debit semampu saldo, sisa jadi tunggakan.
func TestAutodebet_SaldoKurang_PartialDebitDanTunggakan(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	rek := buatRekeningAktif(30000) // saldo 30rb
	jadwal := &autodebet.Jadwal{
		ID:            uuid.New(),
		BMTID:         rek.BMTID,
		RekeningID:    rek.ID,
		ConfigID:      uuid.New(),
		Jenis:         autodebet.JenisSimpananWajib,
		NominalTarget: money.New(100000), // target 100rb
		Status:        autodebet.StatusMenunggu,
		CreatedAt:     time.Now(),
	}

	repo.On("LockForUpdate", mock.Anything, rek.ID).Return(rek, nil)
	// Saldo 30rb didebit (sebagian dari target 100rb)
	repo.On("UpdateSaldo", mock.Anything, rek.ID, int64(0)).Return(nil) // saldo jadi 0
	repo.On("CreateTransaksi", mock.Anything, mock.Anything).Return(nil)
	// Tunggakan 70rb harus dibuat
	autodebetRepo.On("CreateTunggakan", mock.Anything, mock.MatchedBy(func(t *autodebet.Tunggakan) bool {
		return t.NominalSisa == money.New(70000) &&
			t.NominalTerbayar == money.New(30000) &&
			t.Status == "OUTSTANDING"
	})).Return(nil)
	autodebetRepo.On("UpdateJadwalStatus", mock.Anything, jadwal.ID, autodebet.StatusPartial).Return(nil)

	err := svc.EksekusiAutodebetJadwal(context.Background(), jadwal)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	autodebetRepo.AssertExpectations(t)
}

func TestAutodebet_SaldoCukup_TidakAdaTunggakan(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	rek := buatRekeningAktif(200000) // saldo 200rb
	jadwal := &autodebet.Jadwal{
		ID:            uuid.New(),
		BMTID:         rek.BMTID,
		RekeningID:    rek.ID,
		ConfigID:      uuid.New(),
		Jenis:         autodebet.JenisSimpananWajib,
		NominalTarget: money.New(100000), // target 100rb, saldo cukup
		Status:        autodebet.StatusMenunggu,
		CreatedAt:     time.Now(),
	}

	repo.On("LockForUpdate", mock.Anything, rek.ID).Return(rek, nil)
	repo.On("UpdateSaldo", mock.Anything, rek.ID, int64(100000)).Return(nil) // sisa 100rb
	repo.On("CreateTransaksi", mock.Anything, mock.Anything).Return(nil)
	autodebetRepo.On("UpdateJadwalStatus", mock.Anything, jadwal.ID, autodebet.StatusSukses).Return(nil)

	err := svc.EksekusiAutodebetJadwal(context.Background(), jadwal)

	assert.NoError(t, err)
	// CreateTunggakan tidak boleh dipanggil
	autodebetRepo.AssertNotCalled(t, "CreateTunggakan", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	autodebetRepo.AssertExpectations(t)
}

func TestAutodebet_SaldoNol_SemuaJadiTunggakan(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	rek := buatRekeningAktif(0) // saldo kosong
	jadwal := &autodebet.Jadwal{
		ID:            uuid.New(),
		BMTID:         rek.BMTID,
		RekeningID:    rek.ID,
		ConfigID:      uuid.New(),
		Jenis:         autodebet.JenisSimpananWajib,
		NominalTarget: money.New(100000),
		Status:        autodebet.StatusMenunggu,
		CreatedAt:     time.Now(),
	}

	repo.On("LockForUpdate", mock.Anything, rek.ID).Return(rek, nil)
	// NominalDidebit = 0 → UpdateSaldo dan CreateTransaksi tidak dipanggil
	autodebetRepo.On("CreateTunggakan", mock.Anything, mock.MatchedBy(func(t *autodebet.Tunggakan) bool {
		return t.NominalSisa == money.New(100000) &&
			t.NominalTerbayar == money.Zero
	})).Return(nil)
	autodebetRepo.On("UpdateJadwalStatus", mock.Anything, jadwal.ID, autodebet.StatusPartial).Return(nil)

	err := svc.EksekusiAutodebetJadwal(context.Background(), jadwal)

	assert.NoError(t, err)
	repo.AssertNotCalled(t, "UpdateSaldo", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "CreateTransaksi", mock.Anything, mock.Anything)
	autodebetRepo.AssertExpectations(t)
}

// ── Tests: GetSaldo ──────────────────────────────────────────────────────────

// TestCrossTenant_QueryTanpaBMTID_Dilarang adalah test wajib per CLAUDE.md.
// Memastikan tidak ada data yang bocor antar BMT.
func TestCrossTenant_QueryTanpaBMTID_Dilarang(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	bmtIDAsli := uuid.New()
	bmtIDLain := uuid.New() // BMT yang berbeda

	rek := buatRekeningAktif(500000)
	rek.BMTID = bmtIDAsli

	repo.On("GetByID", mock.Anything, rek.ID).Return(rek, nil)

	// Coba akses rekening BMT A dengan menggunakan konteks BMT B
	_, err := svc.GetSaldo(context.Background(), rek.ID, bmtIDLain)

	assert.ErrorIs(t, err, rekening.ErrRekeningNotFound,
		"akses rekening lintas BMT harus menghasilkan ErrRekeningNotFound")
}

func TestRekeningService_GetSaldo_TenantIsolationBenar(t *testing.T) {
	repo := new(MockRekeningRepo)
	autodebetRepo := new(MockAutodebetRepo)
	jurnal := new(MockJurnalService)
	svc := buatService(repo, autodebetRepo, jurnal)

	bmtID := uuid.New()
	rek := buatRekeningAktif(750000)
	rek.BMTID = bmtID

	repo.On("GetByID", mock.Anything, rek.ID).Return(rek, nil)

	saldo, err := svc.GetSaldo(context.Background(), rek.ID, bmtID)

	assert.NoError(t, err)
	assert.Equal(t, money.New(750000), saldo)
	repo.AssertExpectations(t)
}
