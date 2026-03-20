package notif

import "context"

// Provider adalah interface untuk mengirim notifikasi ke satu channel
type Provider interface {
	Kirim(ctx context.Context, tujuan, subjek, pesan string) error
}
