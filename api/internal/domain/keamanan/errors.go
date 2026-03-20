package keamanan

import "errors"

var (
	ErrOTPNotFound           = errors.New("kode OTP tidak ditemukan atau sudah digunakan")
	ErrOTPExpired            = errors.New("OTP sudah expired")
	ErrOTPSalah              = errors.New("OTP tidak sesuai")
	ErrOTPBlokir             = errors.New("terlalu banyak percobaan OTP, coba lagi nanti")
	ErrOTPSudahDigunakan     = errors.New("OTP sudah digunakan")
	ErrTooManyAttempts       = errors.New("terlalu banyak percobaan, akun dikunci sementara")
	ErrSesiNotFound          = errors.New("sesi tidak ditemukan atau sudah tidak aktif")
	ErrKartuNFCTidakAktif    = errors.New("kartu NFC tidak aktif atau expired")
	ErrKartuNFCPINSalah      = errors.New("PIN kartu tidak sesuai")
	ErrIPKioskTidakDiizinkan = errors.New("IP tidak terdaftar sebagai terminal kiosk")
)
