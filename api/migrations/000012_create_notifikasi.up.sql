CREATE TABLE notifikasi_template (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id      UUID REFERENCES bmt(id),
    kode        VARCHAR(50) NOT NULL,
    channel     VARCHAR(20) NOT NULL,
    judul       VARCHAR(255),
    isi         TEXT NOT NULL,
    is_aktif    BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (bmt_id, kode, channel)
);

CREATE TABLE notifikasi_antrian (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    channel         VARCHAR(20) NOT NULL,
    tujuan          VARCHAR(255) NOT NULL,
    subjek          VARCHAR(255),
    pesan           TEXT NOT NULL,
    data_ekstra     JSONB,
    status          VARCHAR(20) NOT NULL DEFAULT 'MENUNGGU',
    percobaan       SMALLINT NOT NULL DEFAULT 0,
    error_terakhir  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    dikirim_at      TIMESTAMPTZ
);

CREATE TABLE notifikasi_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    template_kode   VARCHAR(50) NOT NULL,
    channel         VARCHAR(20) NOT NULL,
    tujuan          VARCHAR(255) NOT NULL,
    isi_terkirim    TEXT NOT NULL,
    status          VARCHAR(20) NOT NULL,
    error_message   TEXT,
    referensi_id    UUID,
    referensi_tipe  VARCHAR(30),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pengumuman (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    judul           VARCHAR(255) NOT NULL,
    isi             TEXT NOT NULL,
    tipe            VARCHAR(20) NOT NULL DEFAULT 'SEMUA',
    target_id       UUID,
    target_asrama   VARCHAR(100),
    file_url        TEXT,
    is_pinned       BOOLEAN NOT NULL DEFAULT FALSE,
    tanggal_mulai   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    tanggal_selesai TIMESTAMPTZ,
    dibuat_oleh     UUID NOT NULL REFERENCES pengguna_pondok(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifikasi_antrian_status ON notifikasi_antrian(status, created_at);
CREATE INDEX idx_pengumuman_bmt ON pengumuman(bmt_id, cabang_id, tanggal_mulai DESC);
