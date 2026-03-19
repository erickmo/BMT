package postgres

import (
	"context"
	"time"

	"github.com/bmt-saas/api/internal/domain/autodebet"
	"github.com/bmt-saas/api/pkg/money"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AutodebetRepository struct {
	db *pgxpool.Pool
}

func NewAutodebetRepository(db *pgxpool.Pool) *AutodebetRepository {
	return &AutodebetRepository{db: db}
}

func (r *AutodebetRepository) CreateConfig(ctx context.Context, c *autodebet.Config) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO rekening_autodebet_config (id, bmt_id, rekening_id, jenis, tanggal_debet, is_aktif, referensi_id, created_at, updated_at, updated_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, c.ID, c.BMTID, c.RekeningID, c.Jenis, c.TanggalDebet, c.IsAktif, c.ReferensiID, c.CreatedAt, c.UpdatedAt, c.UpdatedBy)
	return err
}

func (r *AutodebetRepository) GetConfig(ctx context.Context, id uuid.UUID) (*autodebet.Config, error) {
	c := &autodebet.Config{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, rekening_id, jenis, tanggal_debet, is_aktif, referensi_id, created_at, updated_at, updated_by
		FROM rekening_autodebet_config WHERE id = $1
	`, id).Scan(&c.ID, &c.BMTID, &c.RekeningID, &c.Jenis, &c.TanggalDebet, &c.IsAktif,
		&c.ReferensiID, &c.CreatedAt, &c.UpdatedAt, &c.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *AutodebetRepository) ListConfigByRekening(ctx context.Context, rekeningID uuid.UUID) ([]*autodebet.Config, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, rekening_id, jenis, tanggal_debet, is_aktif, referensi_id, created_at, updated_at, updated_by
		FROM rekening_autodebet_config WHERE rekening_id = $1 AND is_aktif = true
	`, rekeningID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*autodebet.Config
	for rows.Next() {
		c := &autodebet.Config{}
		err := rows.Scan(&c.ID, &c.BMTID, &c.RekeningID, &c.Jenis, &c.TanggalDebet, &c.IsAktif,
			&c.ReferensiID, &c.CreatedAt, &c.UpdatedAt, &c.UpdatedBy)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *AutodebetRepository) ListConfigAktifByBMT(ctx context.Context, bmtID uuid.UUID) ([]*autodebet.Config, error) {
	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.bmt_id, c.rekening_id, c.jenis, c.tanggal_debet, c.is_aktif, c.referensi_id, c.created_at, c.updated_at, c.updated_by
		FROM rekening_autodebet_config c
		JOIN rekening r ON r.id = c.rekening_id
		WHERE c.bmt_id = $1 AND c.is_aktif = true AND r.status = 'AKTIF'
	`, bmtID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*autodebet.Config
	for rows.Next() {
		c := &autodebet.Config{}
		err := rows.Scan(&c.ID, &c.BMTID, &c.RekeningID, &c.Jenis, &c.TanggalDebet, &c.IsAktif,
			&c.ReferensiID, &c.CreatedAt, &c.UpdatedAt, &c.UpdatedBy)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *AutodebetRepository) UpdateConfig(ctx context.Context, c *autodebet.Config) error {
	_, err := r.db.Exec(ctx, `
		UPDATE rekening_autodebet_config SET tanggal_debet=$1, is_aktif=$2, updated_at=NOW(), updated_by=$3
		WHERE id=$4
	`, c.TanggalDebet, c.IsAktif, c.UpdatedBy, c.ID)
	return err
}

func (r *AutodebetRepository) CreateJadwal(ctx context.Context, j *autodebet.Jadwal) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO jadwal_autodebet (id, bmt_id, rekening_id, config_id, jenis, nominal_target, tanggal_jatuh_tempo, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, j.ID, j.BMTID, j.RekeningID, j.ConfigID, j.Jenis, j.NominalTarget.Int64(),
		j.TanggalJatuhTempo, j.Status, j.CreatedAt)
	return err
}

func (r *AutodebetRepository) GetJadwalByID(ctx context.Context, id uuid.UUID) (*autodebet.Jadwal, error) {
	j := &autodebet.Jadwal{}
	var nominal int64
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, rekening_id, config_id, jenis, nominal_target, tanggal_jatuh_tempo, status, created_at
		FROM jadwal_autodebet WHERE id = $1
	`, id).Scan(&j.ID, &j.BMTID, &j.RekeningID, &j.ConfigID, &j.Jenis, &nominal,
		&j.TanggalJatuhTempo, &j.Status, &j.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, autodebet.ErrJadwalNotFound
		}
		return nil, err
	}
	j.NominalTarget = money.New(nominal)
	return j, nil
}

func (r *AutodebetRepository) ListJadwalByTanggal(ctx context.Context, bmtID uuid.UUID, tanggal time.Time) ([]*autodebet.Jadwal, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, rekening_id, config_id, jenis, nominal_target, tanggal_jatuh_tempo, status, created_at
		FROM jadwal_autodebet
		WHERE bmt_id = $1 AND tanggal_jatuh_tempo::date = $2::date AND status = 'MENUNGGU'
	`, bmtID, tanggal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*autodebet.Jadwal
	for rows.Next() {
		j := &autodebet.Jadwal{}
		var nominal int64
		err := rows.Scan(&j.ID, &j.BMTID, &j.RekeningID, &j.ConfigID, &j.Jenis, &nominal,
			&j.TanggalJatuhTempo, &j.Status, &j.CreatedAt)
		if err != nil {
			return nil, err
		}
		j.NominalTarget = money.New(nominal)
		result = append(result, j)
	}
	return result, nil
}

func (r *AutodebetRepository) UpdateJadwalStatus(ctx context.Context, id uuid.UUID, status autodebet.StatusJadwal) error {
	_, err := r.db.Exec(ctx, `UPDATE jadwal_autodebet SET status=$1 WHERE id=$2`, status, id)
	return err
}

func (r *AutodebetRepository) CreateTunggakan(ctx context.Context, t *autodebet.Tunggakan) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO tunggakan_autodebet (id, bmt_id, rekening_id, jadwal_id, jenis, nominal_target, nominal_terbayar, nominal_sisa, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, t.ID, t.BMTID, t.RekeningID, t.JadwalID, t.Jenis, t.NominalTarget.Int64(),
		t.NominalTerbayar.Int64(), t.NominalSisa.Int64(), t.Status, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *AutodebetRepository) ListTunggakanByRekening(ctx context.Context, rekeningID uuid.UUID) ([]*autodebet.Tunggakan, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, rekening_id, jadwal_id, jenis, nominal_target, nominal_terbayar, nominal_sisa, status, created_at, updated_at
		FROM tunggakan_autodebet WHERE rekening_id = $1 AND status = 'OUTSTANDING'
		ORDER BY created_at
	`, rekeningID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*autodebet.Tunggakan
	for rows.Next() {
		t := &autodebet.Tunggakan{}
		var target, terbayar, sisa int64
		err := rows.Scan(&t.ID, &t.BMTID, &t.RekeningID, &t.JadwalID, &t.Jenis,
			&target, &terbayar, &sisa, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		t.NominalTarget = money.New(target)
		t.NominalTerbayar = money.New(terbayar)
		t.NominalSisa = money.New(sisa)
		result = append(result, t)
	}
	return result, nil
}

func (r *AutodebetRepository) UpdateTunggakan(ctx context.Context, t *autodebet.Tunggakan) error {
	_, err := r.db.Exec(ctx, `
		UPDATE tunggakan_autodebet SET nominal_terbayar=$1, nominal_sisa=$2, status=$3, updated_at=NOW()
		WHERE id=$4
	`, t.NominalTerbayar.Int64(), t.NominalSisa.Int64(), t.Status, t.ID)
	return err
}
