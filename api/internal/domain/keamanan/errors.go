package keamanan

import "errors"

var (
	ErrOTPExpired            = errors.New("OTP sudah expired")
	ErrOTPSalah              = errors.New("OTP tidak sesuai")
	ErrOTPSudahDigunakan     = errors.New("OTP sudah digunakan")
	ErrTooManyAttempts       = errors.New("terlalu banyak percobaan, akun dikunci sementara")
	ErrKartuNFCTidakAktif    = errors.New("kartu NFC tidak aktif atau expired")
	ErrKartuNFCPINSalah      = errors.New("PIN kartu tidak sesuai")
	ErrIPKioskTidakDiizinkan = errors.New("IP tidak terdaftar sebagai terminal kiosk")
)
