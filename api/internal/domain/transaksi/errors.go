package transaksi

import "errors"

var (
	ErrIdempotencyDuplikat = errors.New("transaksi dengan idempotency key ini sudah diproses")
	ErrTidakAdaSesiTeller  = errors.New("tidak ada sesi teller aktif")
	ErrSesiTellerSelisih   = errors.New("saldo fisik tidak sesuai, sesi ditolak")
)
