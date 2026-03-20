package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/pembiayaan"
	"github.com/google/uuid"
)

// KolektibilitasService mengklasifikasikan kualitas pembiayaan per standar OJK 5 level.
//
// Level kolektibilitas OJK:
//   1 = Lancar          : tunggak 0 hari
//   2 = Dalam Perhatian : tunggak 1–90 hari
//   3 = Kurang Lancar   : tunggak 91–120 hari
//   4 = Diragukan       : tunggak 121–180 hari
//   5 = Macet           : tunggak > 180 hari
type KolektibilitasService struct {
	pembiayaanRepo pembiayaan.Repository
}

func NewKolektibilitasService(pembiayaanRepo pembiayaan.Repository) *KolektibilitasService {
	return &KolektibilitasService{pembiayaanRepo: pembiayaanRepo}
}

// UpdateBMT memperbarui kolektibilitas semua pembiayaan aktif milik sebuah BMT.
func (s *KolektibilitasService) UpdateBMT(ctx context.Context, bmtID uuid.UUID) error {
	pembiayaans, err := s.pembiayaanRepo.ListAktifByBMT(ctx, bmtID)
	if err != nil {
		return fmt.Errorf("gagal ambil pembiayaan aktif BMT %s: %w", bmtID, err)
	}

	today := time.Now().Truncate(24 * time.Hour)

	for _, p := range pembiayaans {
		// Ambil angsuran yang belum dibayar
		angsurans, err := s.pembiayaanRepo.ListAngsuran(ctx, p.ID)
		if err != nil {
			fmt.Printf("gagal ambil angsuran pembiayaan %s: %v\n", p.ID, err)
			continue
		}

		hariTunggak := hitungHariTunggak(angsurans, today)
		kolektibilitas := hitungKolektibilitas(hariTunggak)

		if int16(hariTunggak) == p.Kolektibilitas && hariTunggak == p.HariTunggak {
			continue // tidak ada perubahan
		}

		if err := s.pembiayaanRepo.UpdateKolektibilitas(ctx, p.ID, kolektibilitas, hariTunggak); err != nil {
			fmt.Printf("gagal update kolektibilitas pembiayaan %s: %v\n", p.ID, err)
		}
	}

	return nil
}

// hitungHariTunggak menghitung jumlah hari tunggak berdasarkan angsuran tertua yang belum dibayar.
func hitungHariTunggak(angsurans []*pembiayaan.AngsuranPembiayaan, today time.Time) int {
	var tertua *time.Time
	for _, a := range angsurans {
		if a.Status == "MENUNGGU" || a.Status == "SEBAGIAN" {
			jatuhTempo := a.TanggalJatuhTempo.Truncate(24 * time.Hour)
			if jatuhTempo.Before(today) {
				if tertua == nil || jatuhTempo.Before(*tertua) {
					t := jatuhTempo
					tertua = &t
				}
			}
		}
	}
	if tertua == nil {
		return 0
	}
	return int(today.Sub(*tertua).Hours() / 24)
}

// hitungKolektibilitas mengklasifikasikan level kolektibilitas OJK berdasarkan hari tunggak.
func hitungKolektibilitas(hariTunggak int) int16 {
	switch {
	case hariTunggak == 0:
		return 1
	case hariTunggak <= 90:
		return 2
	case hariTunggak <= 120:
		return 3
	case hariTunggak <= 180:
		return 4
	default:
		return 5
	}
}
