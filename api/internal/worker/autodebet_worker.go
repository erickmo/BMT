package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/service"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TaskAutodebetHarian      = "autodebet:harian"
	TaskAutodebetBulanan     = "autodebet:bulanan"
	TaskGenerateJadwal       = "autodebet:generate_jadwal"
	TaskGenerateTagihanSPP   = "spp:generate_tagihan"
	TaskKirimNotifikasi      = "notifikasi:kirim"
	TaskPayrollBulanan       = "sdm:payroll"
	TaskGenerateSlipGaji     = "sdm:generate_slip"
	TaskHitungKomisiOPOP     = "opop:hitung_komisi"
	TaskSinkronEMIS          = "integrasi:emis"
	TaskSinkronDAPODIK       = "integrasi:dapodik"
	TaskAnalyticsHarian      = "analytics:snapshot"
	TaskCleanupSesi          = "keamanan:cleanup_sesi"
	TaskUpdateKolektibilitas = "pembiayaan:update_kolektibilitas"
	TaskDistribusiBagiHasil  = "deposito:distribusi_bagi_hasil"
	TaskReminderAngsuran     = "pembiayaan:reminder_angsuran"
	TaskReminderSPP          = "spp:reminder"
	TaskCekMidtransPending   = "midtrans:cek_pending"
	TaskCekKontrakExpiry     = "platform:cek_kontrak"
	TaskExpiredKartuNFC      = "nfc:expired_kartu"
	TaskBackupDatabase       = "platform:backup_db"
	TaskReminderPerpus       = "perpus:reminder"
	TaskGenerateRaport       = "akademik:generate_raport"
	TaskHitungZakat          = "zakat:hitung_mal"
)

type PayloadBMT struct {
	BMTID uuid.UUID `json:"bmt_id"`
}

type PayloadAutodebetHarian struct {
	BMTID   uuid.UUID `json:"bmt_id"`
	Tanggal time.Time `json:"tanggal"`
}

// AutodebetWorker menangani task autodebet harian, bulanan, dan generate jadwal.
type AutodebetWorker struct {
	autodebetService *service.AutodebetService
}

func NewAutodebetWorker(autodebetService *service.AutodebetService) *AutodebetWorker {
	return &AutodebetWorker{autodebetService: autodebetService}
}

func (w *AutodebetWorker) HandleAutodebetHarian(ctx context.Context, t *asynq.Task) error {
	var payload PayloadAutodebetHarian
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("gagal unmarshal payload: %w", err)
	}
	return w.autodebetService.EksekusiHarian(ctx, payload.BMTID, payload.Tanggal)
}

func (w *AutodebetWorker) HandleAutodebetBulanan(ctx context.Context, t *asynq.Task) error {
	var payload PayloadBMT
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("gagal unmarshal payload: %w", err)
	}
	return w.autodebetService.EksekusiBulanan(ctx, payload.BMTID)
}

func (w *AutodebetWorker) HandleGenerateJadwal(ctx context.Context, t *asynq.Task) error {
	var payload PayloadBMT
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("gagal unmarshal payload: %w", err)
	}
	return w.autodebetService.GenerateJadwalBulanDepan(ctx, payload.BMTID)
}

// CBSWorker menangani task CBS: kolektibilitas, distribusi bagi hasil, reminder angsuran, dan zakat.
type CBSWorker struct {
	kolektibilitasSvc *service.KolektibilitasService
	distribusiSvc     *service.DistribusiService
	reminderSvc       *service.ReminderService
	zakatSvc          *service.ZakatService
}

func NewCBSWorker(
	kolektibilitasSvc *service.KolektibilitasService,
	distribusiSvc *service.DistribusiService,
	reminderSvc *service.ReminderService,
	zakatSvc *service.ZakatService,
) *CBSWorker {
	return &CBSWorker{
		kolektibilitasSvc: kolektibilitasSvc,
		distribusiSvc:     distribusiSvc,
		reminderSvc:       reminderSvc,
		zakatSvc:          zakatSvc,
	}
}

// HandleUpdateKolektibilitas mengklasifikasikan ulang kolektibilitas semua pembiayaan aktif BMT.
// Payload: { "bmt_id": "..." }
func (w *CBSWorker) HandleUpdateKolektibilitas(ctx context.Context, t *asynq.Task) error {
	var payload PayloadBMT
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("gagal unmarshal payload: %w", err)
	}
	return w.kolektibilitasSvc.UpdateBMT(ctx, payload.BMTID)
}

// HandleDistribusiBagiHasil mendistribusikan bagi hasil deposito akhir bulan.
// Dijalankan scheduler pada 28–31 bulan, worker mengecek apakah hari ini adalah hari terakhir bulan.
func (w *CBSWorker) HandleDistribusiBagiHasil(ctx context.Context, t *asynq.Task) error {
	var payload PayloadBMT
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("gagal unmarshal payload: %w", err)
	}

	now := time.Now()
	// Hanya jalankan pada hari terakhir bulan
	besok := now.AddDate(0, 0, 1)
	if besok.Month() == now.Month() {
		return nil // bukan hari terakhir bulan
	}

	hasil, err := w.distribusiSvc.DistribusiBagiHasil(ctx, payload.BMTID, now)
	if err != nil {
		return err
	}
	fmt.Printf("[DistribusiBagiHasil] BMT %s: %d rekening diproses\n", payload.BMTID, len(hasil))
	return nil
}

// HandleReminderAngsuran mengirim reminder H-N ke nasabah yang angsurannya akan jatuh tempo.
func (w *CBSWorker) HandleReminderAngsuran(ctx context.Context, t *asynq.Task) error {
	var payload PayloadBMT
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("gagal unmarshal payload: %w", err)
	}
	return w.reminderSvc.ReminderAngsuranBMT(ctx, payload.BMTID)
}

// HandleHitungZakat menghitung kewajiban zakat mal akhir tahun berdasarkan nisab dari settings.
// Payload: { "bmt_id": "..." }
func (w *CBSWorker) HandleHitungZakat(ctx context.Context, t *asynq.Task) error {
	var payload PayloadBMT
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("gagal unmarshal payload: %w", err)
	}
	count, err := w.zakatSvc.HitungZakatBMT(ctx, payload.BMTID)
	if err != nil {
		return err
	}
	fmt.Printf("[HandleHitungZakat] BMT %s: %d rekening memenuhi nisab zakat\n", payload.BMTID, count)
	return nil
}

// RegisterWorkers mendaftarkan semua task handlers ke ServeMux.
func RegisterWorkers(mux *asynq.ServeMux, autodebetWorker *AutodebetWorker, cbsWorker *CBSWorker, notifikasiWorker *NotifikasiWorker) {
	mux.HandleFunc(TaskAutodebetHarian, autodebetWorker.HandleAutodebetHarian)
	mux.HandleFunc(TaskAutodebetBulanan, autodebetWorker.HandleAutodebetBulanan)
	mux.HandleFunc(TaskGenerateJadwal, autodebetWorker.HandleGenerateJadwal)
	mux.HandleFunc(TaskUpdateKolektibilitas, cbsWorker.HandleUpdateKolektibilitas)
	mux.HandleFunc(TaskDistribusiBagiHasil, cbsWorker.HandleDistribusiBagiHasil)
	mux.HandleFunc(TaskReminderAngsuran, cbsWorker.HandleReminderAngsuran)
	mux.HandleFunc(TaskHitungZakat, cbsWorker.HandleHitungZakat)
	mux.HandleFunc(TaskKirimNotifikasi, notifikasiWorker.HandleKirimNotifikasi)
	mux.HandleFunc(TaskCekMidtransPending, notifikasiWorker.HandleCekMidtransPending)
}

// SchedulePeriodicTasks mendaftarkan semua cron jobs ke Scheduler.
// Semua waktu ditulis dalam UTC; jam operasional WIB = UTC+7.
func SchedulePeriodicTasks(scheduler *asynq.Scheduler) {
	// Autodebet harian: 07:00 WIB = 00:00 UTC
	scheduler.Register("0 0 * * *", asynq.NewTask(TaskAutodebetHarian, nil))

	// Generate jadwal autodebet: tgl 25, 08:00 WIB = 01:00 UTC
	scheduler.Register("0 1 25 * *", asynq.NewTask(TaskGenerateJadwal, nil))

	// Generate tagihan SPP: tgl 25, 08:30 WIB = 01:30 UTC
	scheduler.Register("30 1 25 * *", asynq.NewTask(TaskGenerateTagihanSPP, nil))

	// Update kolektibilitas: 00:05 WIB = 17:05 UTC sehari sebelumnya
	scheduler.Register("5 17 * * *", asynq.NewTask(TaskUpdateKolektibilitas, nil))

	// Distribusi bagi hasil: akhir bulan (hari terakhir) 22:00 WIB = 15:00 UTC
	scheduler.Register("0 15 28-31 * *", asynq.NewTask(TaskDistribusiBagiHasil, nil))

	// Cek Midtrans pending: setiap 15 menit
	scheduler.Register("*/15 * * * *", asynq.NewTask(TaskCekMidtransPending, nil))

	// Analytics harian: 23:00 WIB = 16:00 UTC
	scheduler.Register("0 16 * * *", asynq.NewTask(TaskAnalyticsHarian, nil))

	// Cleanup sesi expired: setiap jam
	scheduler.Register("0 * * * *", asynq.NewTask(TaskCleanupSesi, nil))

	// Kirim notifikasi: setiap 1 menit
	scheduler.Register("* * * * *", asynq.NewTask(TaskKirimNotifikasi, nil))

	// Cek kontrak expiry: 08:00 WIB = 01:00 UTC
	scheduler.Register("0 1 * * *", asynq.NewTask(TaskCekKontrakExpiry, nil))

	// Expired kartu NFC: 08:00 WIB = 01:00 UTC
	scheduler.Register("5 1 * * *", asynq.NewTask(TaskExpiredKartuNFC, nil))

	// Backup database: 02:00 WIB = 19:00 UTC sehari sebelumnya
	scheduler.Register("0 19 * * *", asynq.NewTask(TaskBackupDatabase, nil))

	// Hitung komisi OPOP: tgl 1, 07:00 WIB = 00:00 UTC
	scheduler.Register("0 0 1 * *", asynq.NewTask(TaskHitungKomisiOPOP, nil))

	// Generate slip gaji: tgl 25, 09:00 WIB = 02:00 UTC
	scheduler.Register("0 2 25 * *", asynq.NewTask(TaskGenerateSlipGaji, nil))

	// Eksekusi payroll: tgl 1, 08:00 WIB = 01:00 UTC
	scheduler.Register("0 1 1 * *", asynq.NewTask(TaskPayrollBulanan, nil))

	// Sinkron EMIS: Sabtu 02:00 WIB = Jumat 19:00 UTC
	scheduler.Register("0 19 * * 5", asynq.NewTask(TaskSinkronEMIS, nil))

	// Sinkron DAPODIK: tgl 1, 03:00 WIB = 20:00 UTC sehari sebelumnya
	scheduler.Register("0 20 1 * *", asynq.NewTask(TaskSinkronDAPODIK, nil))

	// Reminder perpustakaan: 08:00 WIB = 01:00 UTC
	scheduler.Register("10 1 * * *", asynq.NewTask(TaskReminderPerpus, nil))

	// Reminder SPP: terkonfigurasi dari settings BMT
	scheduler.Register("0 2 * * *", asynq.NewTask(TaskReminderSPP, nil))

	// Reminder angsuran pembiayaan
	scheduler.Register("15 2 * * *", asynq.NewTask(TaskReminderAngsuran, nil))

	// Generate raport: dikonfigurasi per BMT, worker cek setiap hari
	scheduler.Register("0 3 * * *", asynq.NewTask(TaskGenerateRaport, nil))

	// Hitung zakat mal: 1 Januari, 06:00 WIB = 23:00 UTC sehari sebelumnya
	// Worker membaca ZAKAT_NISAB_RUPIAH dari settings; skip jika 0.
	scheduler.Register("0 23 31 12 *", asynq.NewTask(TaskHitungZakat, nil))
}
