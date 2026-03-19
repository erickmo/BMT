package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bmt-saas/api/internal/domain/sesi_teller"
	"github.com/bmt-saas/api/pkg/money"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SesiTellerRepository struct {
	db *pgxpool.Pool
}

func NewSesiTellerRepository(db *pgxpool.Pool) *SesiTellerRepository {
	return &SesiTellerRepository{db: db}
}

func (r *SesiTellerRepository) Create(ctx context.Context, s *sesi_teller.SesiTeller) error {
	redenominasi, _ := json.Marshal(s.Redenominasi)
	_, err := r.db.Exec(ctx, `
		INSERT INTO sesi_teller (id, bmt_id, cabang_id, teller_id, tanggal, saldo_awal, redenominasi,
		status, toleransi_selisih, dibuka_pada, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, s.ID, s.BMTID, s.CabangID, s.TellerID, s.Tanggal, s.SaldoAwal.Int64(), redenominasi,
		s.Status, s.ToleransiSelisih.Int64(), s.DibukaPada, s.DibukaPada)
	return err
}

func (r *SesiTellerRepository) GetAktifByTeller(ctx context.Context, tellerID uuid.UUID) (*sesi_teller.SesiTeller, error) {
	return r.scanSesi(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, teller_id, tanggal, saldo_awal, redenominasi,
		saldo_akhir, redenominasi_akhir, status, toleransi_selisih, selisih, dibuka_pada, ditutup_pada
		FROM sesi_teller WHERE teller_id = $1 AND status = 'AKTIF'
		ORDER BY dibuka_pada DESC LIMIT 1
	`, tellerID))
}

func (r *SesiTellerRepository) GetByID(ctx context.Context, id uuid.UUID) (*sesi_teller.SesiTeller, error) {
	return r.scanSesi(r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, teller_id, tanggal, saldo_awal, redenominasi,
		saldo_akhir, redenominasi_akhir, status, toleransi_selisih, selisih, dibuka_pada, ditutup_pada
		FROM sesi_teller WHERE id = $1
	`, id))
}

func (r *SesiTellerRepository) Update(ctx context.Context, s *sesi_teller.SesiTeller) error {
	redenominasiAkhir, _ := json.Marshal(s.RedenominasiAkhir)
	var saldoAkhir *int64
	if s.SaldoAkhir != nil {
		v := s.SaldoAkhir.Int64()
		saldoAkhir = &v
	}
	var selisih *int64
	if s.Selisih != nil {
		v := s.Selisih.Int64()
		selisih = &v
	}
	_, err := r.db.Exec(ctx, `
		UPDATE sesi_teller SET saldo_akhir=$1, redenominasi_akhir=$2, status=$3, selisih=$4, ditutup_pada=$5
		WHERE id=$6
	`, saldoAkhir, redenominasiAkhir, s.Status, selisih, s.DitutupPada, s.ID)
	return err
}

func (r *SesiTellerRepository) ListByBMT(ctx context.Context, bmtID, cabangID uuid.UUID, tanggal time.Time) ([]*sesi_teller.SesiTeller, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, teller_id, tanggal, saldo_awal, redenominasi,
		saldo_akhir, redenominasi_akhir, status, toleransi_selisih, selisih, dibuka_pada, ditutup_pada
		FROM sesi_teller WHERE bmt_id = $1 AND cabang_id = $2 AND tanggal::date = $3::date
		ORDER BY dibuka_pada DESC
	`, bmtID, cabangID, tanggal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*sesi_teller.SesiTeller
	for rows.Next() {
		s, err := r.scanSesi(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *SesiTellerRepository) scanSesi(s scanner) (*sesi_teller.SesiTeller, error) {
	sesi := &sesi_teller.SesiTeller{}
	var saldoAwal int64
	var saldoAkhir *int64
	var selisih *int64
	var redenominasi, redenominasiAkhir []byte
	var toleransi int64

	err := s.Scan(&sesi.ID, &sesi.BMTID, &sesi.CabangID, &sesi.TellerID, &sesi.Tanggal,
		&saldoAwal, &redenominasi, &saldoAkhir, &redenominasiAkhir,
		&sesi.Status, &toleransi, &selisih, &sesi.DibukaPada, &sesi.DitutupPada)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, sesi_teller.ErrSesiTidakAktif
		}
		return nil, err
	}

	sesi.SaldoAwal = money.New(saldoAwal)
	sesi.ToleransiSelisih = money.New(toleransi)

	if saldoAkhir != nil {
		v := money.New(*saldoAkhir)
		sesi.SaldoAkhir = &v
	}
	if selisih != nil {
		v := money.New(*selisih)
		sesi.Selisih = &v
	}

	json.Unmarshal(redenominasi, &sesi.Redenominasi)
	if redenominasiAkhir != nil {
		json.Unmarshal(redenominasiAkhir, &sesi.RedenominasiAkhir)
	}

	return sesi, nil
}
