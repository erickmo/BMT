package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/bmt-saas/api/internal/domain/notifikasi"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotifikasiRepository struct {
	db *pgxpool.Pool
}

func NewNotifikasiRepository(db *pgxpool.Pool) *NotifikasiRepository {
	return &NotifikasiRepository{db: db}
}

// ── Template ──────────────────────────────────────────────────────────────────

func (r *NotifikasiRepository) GetTemplate(ctx context.Context, bmtID *uuid.UUID, kode string, channel notifikasi.Channel) (*notifikasi.NotifikasiTemplate, error) {
	t := &notifikasi.NotifikasiTemplate{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, kode, channel, judul, isi, is_aktif, created_at, updated_at
		FROM notifikasi_template
		WHERE (bmt_id = $1 OR bmt_id IS NULL) AND kode = $2 AND channel = $3 AND is_aktif = true
		ORDER BY bmt_id NULLS LAST
		LIMIT 1
	`, bmtID, kode, channel).Scan(
		&t.ID, &t.BMTID, &t.Kode, &t.Channel, &t.Judul, &t.Isi, &t.IsAktif, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, notifikasi.ErrTemplateNotFound
		}
		return nil, fmt.Errorf("get template notifikasi: %w", err)
	}
	return t, nil
}

func (r *NotifikasiRepository) ListTemplate(ctx context.Context, bmtID uuid.UUID) ([]*notifikasi.NotifikasiTemplate, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, kode, channel, judul, isi, is_aktif, created_at, updated_at
		FROM notifikasi_template
		WHERE bmt_id = $1 OR bmt_id IS NULL
		ORDER BY kode, channel
	`, bmtID)
	if err != nil {
		return nil, fmt.Errorf("list template notifikasi: %w", err)
	}
	defer rows.Close()

	var result []*notifikasi.NotifikasiTemplate
	for rows.Next() {
		t := &notifikasi.NotifikasiTemplate{}
		err := rows.Scan(&t.ID, &t.BMTID, &t.Kode, &t.Channel, &t.Judul, &t.Isi, &t.IsAktif, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan template notifikasi: %w", err)
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *NotifikasiRepository) UpsertTemplate(ctx context.Context, t *notifikasi.NotifikasiTemplate) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO notifikasi_template (id, bmt_id, kode, channel, judul, isi, is_aktif, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (bmt_id, kode, channel) DO UPDATE
		SET judul = EXCLUDED.judul,
		    isi = EXCLUDED.isi,
		    is_aktif = EXCLUDED.is_aktif,
		    updated_at = NOW()
	`, t.ID, t.BMTID, t.Kode, t.Channel, t.Judul, t.Isi, t.IsAktif, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert template notifikasi: %w", err)
	}
	return nil
}

// ── Antrian ───────────────────────────────────────────────────────────────────

func (r *NotifikasiRepository) CreateAntrian(ctx context.Context, a *notifikasi.NotifikasiAntrian) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO notifikasi_antrian (id, bmt_id, channel, tujuan, subjek, pesan, data_ekstra, status, percobaan, error_terakhir, created_at, dikirim_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, a.ID, a.BMTID, a.Channel, a.Tujuan, a.Subjek, a.Pesan, a.DataEkstra,
		a.Status, a.Percobaan, a.ErrorTerakhir, a.CreatedAt, a.DikirimAt)
	if err != nil {
		return fmt.Errorf("create antrian notifikasi: %w", err)
	}
	return nil
}

func (r *NotifikasiRepository) GetAntrianMenunggu(ctx context.Context, limit int) ([]*notifikasi.NotifikasiAntrian, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, channel, tujuan, subjek, pesan, data_ekstra, status, percobaan, error_terakhir, created_at, dikirim_at
		FROM notifikasi_antrian
		WHERE status = 'MENUNGGU'
		ORDER BY created_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("get antrian menunggu: %w", err)
	}
	defer rows.Close()

	var result []*notifikasi.NotifikasiAntrian
	for rows.Next() {
		a := &notifikasi.NotifikasiAntrian{}
		err := rows.Scan(&a.ID, &a.BMTID, &a.Channel, &a.Tujuan, &a.Subjek, &a.Pesan,
			&a.DataEkstra, &a.Status, &a.Percobaan, &a.ErrorTerakhir, &a.CreatedAt, &a.DikirimAt)
		if err != nil {
			return nil, fmt.Errorf("scan antrian notifikasi: %w", err)
		}
		result = append(result, a)
	}
	return result, nil
}

func (r *NotifikasiRepository) UpdateStatusAntrian(ctx context.Context, id uuid.UUID, status notifikasi.StatusAntrian, errorMsg string) error {
	var dikirimAt *time.Time
	if status == notifikasi.StatusTerkirim {
		now := time.Now()
		dikirimAt = &now
	}
	_, err := r.db.Exec(ctx, `
		UPDATE notifikasi_antrian
		SET status = $1, error_terakhir = $2, dikirim_at = $3
		WHERE id = $4
	`, status, errorMsg, dikirimAt, id)
	if err != nil {
		return fmt.Errorf("update status antrian notifikasi: %w", err)
	}
	return nil
}

func (r *NotifikasiRepository) IncrementPercobaan(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifikasi_antrian SET percobaan = percobaan + 1 WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("increment percobaan notifikasi: %w", err)
	}
	return nil
}

// ── Log ───────────────────────────────────────────────────────────────────────

func (r *NotifikasiRepository) CreateLog(ctx context.Context, l *notifikasi.NotifikasiLog) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO notifikasi_log (id, bmt_id, template_kode, channel, tujuan, isi_terkirim, status, error_message, referensi_id, referensi_tipe, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, l.ID, l.BMTID, l.TemplateKode, l.Channel, l.Tujuan, l.IsiTerkirim,
		l.Status, l.ErrorMessage, l.ReferensiID, l.ReferensiTipe, l.CreatedAt)
	if err != nil {
		return fmt.Errorf("create log notifikasi: %w", err)
	}
	return nil
}

// ── Pengumuman ────────────────────────────────────────────────────────────────

func (r *NotifikasiRepository) CreatePengumuman(ctx context.Context, p *notifikasi.Pengumuman) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO pengumuman (id, bmt_id, cabang_id, judul, isi, tipe, target_id, target_asrama, file_url, is_pinned, tanggal_mulai, tanggal_selesai, dibuat_oleh, is_aktif, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`, p.ID, p.BMTID, p.CabangID, p.Judul, p.Isi, p.Tipe, p.TargetID, p.TargetAsrama,
		p.FileURL, p.IsPinned, p.TanggalMulai, p.TanggalSelesai, p.DibuatOleh, p.IsAktif, p.CreatedAt)
	if err != nil {
		return fmt.Errorf("create pengumuman: %w", err)
	}
	return nil
}

func (r *NotifikasiRepository) GetPengumumanByID(ctx context.Context, id uuid.UUID) (*notifikasi.Pengumuman, error) {
	p := &notifikasi.Pengumuman{}
	err := r.db.QueryRow(ctx, `
		SELECT id, bmt_id, cabang_id, judul, isi, tipe, target_id, target_asrama, file_url, is_pinned, tanggal_mulai, tanggal_selesai, dibuat_oleh, is_aktif, created_at
		FROM pengumuman WHERE id = $1
	`, id).Scan(&p.ID, &p.BMTID, &p.CabangID, &p.Judul, &p.Isi, &p.Tipe, &p.TargetID, &p.TargetAsrama,
		&p.FileURL, &p.IsPinned, &p.TanggalMulai, &p.TanggalSelesai, &p.DibuatOleh, &p.IsAktif, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, notifikasi.ErrPengumumanNotFound
		}
		return nil, fmt.Errorf("get pengumuman by id: %w", err)
	}
	return p, nil
}

func (r *NotifikasiRepository) ListPengumuman(ctx context.Context, bmtID, cabangID uuid.UUID, tipe notifikasi.TargetPengumuman, page, perPage int) ([]*notifikasi.Pengumuman, int64, error) {
	offset := (page - 1) * perPage
	rows, err := r.db.Query(ctx, `
		SELECT id, bmt_id, cabang_id, judul, isi, tipe, target_id, target_asrama, file_url, is_pinned, tanggal_mulai, tanggal_selesai, dibuat_oleh, is_aktif, created_at
		FROM pengumuman
		WHERE bmt_id = $1 AND cabang_id = $2 AND tipe = $3 AND is_aktif = true
		ORDER BY is_pinned DESC, created_at DESC
		LIMIT $4 OFFSET $5
	`, bmtID, cabangID, tipe, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list pengumuman: %w", err)
	}
	defer rows.Close()

	var result []*notifikasi.Pengumuman
	for rows.Next() {
		p := &notifikasi.Pengumuman{}
		err := rows.Scan(&p.ID, &p.BMTID, &p.CabangID, &p.Judul, &p.Isi, &p.Tipe, &p.TargetID, &p.TargetAsrama,
			&p.FileURL, &p.IsPinned, &p.TanggalMulai, &p.TanggalSelesai, &p.DibuatOleh, &p.IsAktif, &p.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan pengumuman: %w", err)
		}
		result = append(result, p)
	}

	var total int64
	r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM pengumuman WHERE bmt_id = $1 AND cabang_id = $2 AND tipe = $3 AND is_aktif = true
	`, bmtID, cabangID, tipe).Scan(&total)

	return result, total, nil
}

func (r *NotifikasiRepository) UpdatePengumuman(ctx context.Context, p *notifikasi.Pengumuman) error {
	_, err := r.db.Exec(ctx, `
		UPDATE pengumuman
		SET judul = $1, isi = $2, tipe = $3, target_id = $4, target_asrama = $5,
		    file_url = $6, is_pinned = $7, tanggal_mulai = $8, tanggal_selesai = $9, is_aktif = $10
		WHERE id = $11
	`, p.Judul, p.Isi, p.Tipe, p.TargetID, p.TargetAsrama, p.FileURL,
		p.IsPinned, p.TanggalMulai, p.TanggalSelesai, p.IsAktif, p.ID)
	if err != nil {
		return fmt.Errorf("update pengumuman: %w", err)
	}
	return nil
}

func (r *NotifikasiRepository) MarkBaca(ctx context.Context, b *notifikasi.PengumumanBaca) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO pengumuman_baca (pengumuman_id, nasabah_id, pengguna_id, dibaca_at)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (pengumuman_id, COALESCE(nasabah_id, '00000000-0000-0000-0000-000000000000'), COALESCE(pengguna_id, '00000000-0000-0000-0000-000000000000'))
		DO NOTHING
	`, b.PengumumanID, b.NasabahID, b.PenggunaID, b.DibacaAt)
	if err != nil {
		return fmt.Errorf("mark baca pengumuman: %w", err)
	}
	return nil
}
