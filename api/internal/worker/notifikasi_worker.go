package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bmt-saas/api/internal/service"
	"github.com/hibiken/asynq"
)

// NotifikasiWorker menangani task pengiriman notifikasi dan cek Midtrans pending.
type NotifikasiWorker struct {
	notifikasiSvc *service.NotifikasiService
	midtransSvc   *service.MidtransService
}

func NewNotifikasiWorker(notifikasiSvc *service.NotifikasiService, midtransSvc *service.MidtransService) *NotifikasiWorker {
	return &NotifikasiWorker{
		notifikasiSvc: notifikasiSvc,
		midtransSvc:   midtransSvc,
	}
}

// HandleKirimNotifikasi memproses antrian notifikasi MENUNGGU dalam batch.
// Setiap antrian dikirim via provider yang sesuai dengan retry max 3x.
func (w *NotifikasiWorker) HandleKirimNotifikasi(ctx context.Context, t *asynq.Task) error {
	_ = json.RawMessage(t.Payload()) // payload kosong untuk task ini

	antrians, err := w.notifikasiSvc.GetPendingAntrian(ctx, 100)
	if err != nil {
		return fmt.Errorf("gagal ambil antrian notifikasi: %w", err)
	}

	var errs []error
	for _, a := range antrians {
		if err := w.notifikasiSvc.DeliverAntrian(ctx, a); err != nil {
			errs = append(errs, fmt.Errorf("antrian %s: %w", a.ID, err))
		}
	}

	if len(errs) > 0 {
		fmt.Printf("[NotifikasiWorker] %d antrian gagal dari %d\n", len(errs), len(antrians))
	}

	return nil
}

// HandleCekMidtransPending memverifikasi status transaksi PENDING > 30 menit ke Midtrans.
func (w *NotifikasiWorker) HandleCekMidtransPending(ctx context.Context, t *asynq.Task) error {
	if err := w.midtransSvc.CekPending(ctx); err != nil {
		return fmt.Errorf("gagal cek Midtrans pending: %w", err)
	}
	return nil
}
