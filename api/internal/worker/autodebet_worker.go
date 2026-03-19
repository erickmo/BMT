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
)

type PayloadBMT struct {
	BMTID uuid.UUID `json:"bmt_id"`
}

type PayloadAutodebetHarian struct {
	BMTID   uuid.UUID `json:"bmt_id"`
	Tanggal time.Time `json:"tanggal"`
}

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

func (w *AutodebetWorker) HandleGenerateJadwal(ctx context.Context, t *asynq.Task) error {
	var payload PayloadBMT
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("gagal unmarshal payload: %w", err)
	}

	// Generate jadwal untuk bulan depan.
	// RekeningIDs kosong karena GenerateJadwalBulanan perlu di-populate oleh caller
	// dengan semua rekening aktif milik BMT tersebut.
	bulanDepan := time.Now().AddDate(0, 1, 0)
	return w.autodebetService.GenerateJadwalBulanan(ctx, []uuid.UUID{}, payload.BMTID, bulanDepan)
}

// RegisterWorkers mendaftarkan semua task handlers ke ServeMux.
func RegisterWorkers(mux *asynq.ServeMux, autodebetWorker *AutodebetWorker) {
	mux.HandleFunc(TaskAutodebetHarian, autodebetWorker.HandleAutodebetHarian)
	mux.HandleFunc(TaskGenerateJadwal, autodebetWorker.HandleGenerateJadwal)
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
}
