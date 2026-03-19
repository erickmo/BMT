-- Platform/BMT tables
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE bmt (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kode            VARCHAR(20) UNIQUE NOT NULL,
    nama            VARCHAR(255) NOT NULL,
    alamat          TEXT,
    telepon         VARCHAR(20),
    email           VARCHAR(255),
    logo_url        TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    whitelabel      JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE kontrak_bmt (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    tanggal_mulai   DATE NOT NULL,
    tanggal_selesai DATE NOT NULL,
    fitur           JSONB NOT NULL DEFAULT '{}',
    tarif           JSONB NOT NULL DEFAULT '{}',
    pic_nama        VARCHAR(255),
    pic_telepon     VARCHAR(20),
    pic_email       VARCHAR(255),
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cabang (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    kode            VARCHAR(20) NOT NULL,
    nama            VARCHAR(255) NOT NULL,
    alamat          TEXT,
    telepon         VARCHAR(20),
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (bmt_id, kode)
);

CREATE TABLE platform_settings (
    kunci           VARCHAR(150) PRIMARY KEY,
    nilai           TEXT NOT NULL,
    tipe            VARCHAR(20) NOT NULL DEFAULT 'string',
    deskripsi       TEXT,
    is_rahasia      BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      TEXT NOT NULL
);

CREATE TABLE bmt_settings (
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    kunci           VARCHAR(150) NOT NULL,
    nilai           TEXT NOT NULL,
    tipe            VARCHAR(20) NOT NULL DEFAULT 'string',
    is_locked       BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID NOT NULL,
    PRIMARY KEY (bmt_id, kunci)
);

CREATE TABLE cabang_settings (
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    kunci           VARCHAR(150) NOT NULL,
    nilai           TEXT NOT NULL,
    tipe            VARCHAR(20) NOT NULL DEFAULT 'string',
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID NOT NULL,
    PRIMARY KEY (cabang_id, kunci)
);

CREATE TABLE pecahan_uang (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nominal         BIGINT NOT NULL,
    jenis           VARCHAR(10) NOT NULL CHECK (jenis IN ('LOGAM','KERTAS')),
    label           VARCHAR(30) NOT NULL,
    is_aktif        BOOLEAN NOT NULL DEFAULT TRUE,
    urutan          SMALLINT NOT NULL,
    berlaku_sejak   DATE NOT NULL,
    ditarik_pada    DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (nominal, jenis)
);

-- Seed default pecahan uang
INSERT INTO pecahan_uang (nominal, jenis, label, urutan, berlaku_sejak) VALUES
(100000, 'KERTAS', 'Rp 100.000', 1, '2022-08-17'),
(75000,  'KERTAS', 'Rp 75.000',  2, '2022-08-17'),
(50000,  'KERTAS', 'Rp 50.000',  3, '2004-08-01'),
(20000,  'KERTAS', 'Rp 20.000',  4, '2004-08-01'),
(10000,  'KERTAS', 'Rp 10.000',  5, '2005-11-28'),
(5000,   'KERTAS', 'Rp 5.000',   6, '2016-12-19'),
(2000,   'KERTAS', 'Rp 2.000',   7, '2016-12-19'),
(1000,   'KERTAS', 'Rp 1.000',   8, '2016-12-19'),
(1000,   'LOGAM',  'Rp 1.000 (logam)', 9, '2010-01-01'),
(500,    'LOGAM',  'Rp 500',     10, '2003-01-01'),
(200,    'LOGAM',  'Rp 200',     11, '2003-01-01'),
(100,    'LOGAM',  'Rp 100',     12, '1999-01-01');

CREATE TABLE usage_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    jenis           VARCHAR(50) NOT NULL,
    referensi_id    UUID,
    nominal_fee     BIGINT NOT NULL DEFAULT 0,
    keterangan      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_usage_log_bmt ON usage_log(bmt_id, created_at);
