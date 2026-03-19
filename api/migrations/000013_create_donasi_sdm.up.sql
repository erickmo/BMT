-- Donasi & Wakaf
CREATE TABLE program_donasi (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    nama            VARCHAR(255) NOT NULL,
    deskripsi       TEXT,
    tipe            VARCHAR(20) NOT NULL,
    target_nominal  BIGINT,
    terkumpul       BIGINT NOT NULL DEFAULT 0,
    tanggal_mulai   DATE NOT NULL,
    tanggal_selesai DATE,
    foto_url        TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    rekening_id     UUID NOT NULL REFERENCES rekening(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE transaksi_donasi (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    program_id      UUID NOT NULL REFERENCES program_donasi(id),
    nasabah_id      UUID REFERENCES nasabah(id),
    nominal         BIGINT NOT NULL CHECK (nominal > 0),
    is_anonim       BOOLEAN NOT NULL DEFAULT FALSE,
    pesan           TEXT,
    metode          VARCHAR(30) NOT NULL,
    midtrans_order_id VARCHAR(100) UNIQUE,
    rekening_id     UUID REFERENCES rekening(id),
    idempotency_key UUID UNIQUE,
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- SDM & Payroll
CREATE TABLE sdm_kontrak (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    pegawai_id      UUID NOT NULL,
    tipe_pegawai    VARCHAR(20) NOT NULL,
    nomor_kontrak   VARCHAR(40) UNIQUE NOT NULL,
    tipe_kontrak    VARCHAR(20) NOT NULL,
    tanggal_mulai   DATE NOT NULL,
    tanggal_selesai DATE,
    gaji_pokok      BIGINT NOT NULL,
    tunjangan       JSONB NOT NULL DEFAULT '{}',
    potongan_tetap  JSONB NOT NULL DEFAULT '{}',
    rekening_gaji_id UUID REFERENCES rekening(id),
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    dokumen_url     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sdm_slip_gaji (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    kontrak_id      UUID NOT NULL REFERENCES sdm_kontrak(id),
    periode         CHAR(7) NOT NULL,
    gaji_pokok      BIGINT NOT NULL,
    tunjangan_total BIGINT NOT NULL DEFAULT 0,
    tunjangan_detail JSONB NOT NULL DEFAULT '{}',
    potongan_absensi BIGINT NOT NULL DEFAULT 0,
    potongan_tetap  BIGINT NOT NULL DEFAULT 0,
    potongan_lain   BIGINT NOT NULL DEFAULT 0,
    gaji_bersih     BIGINT NOT NULL,
    hari_kerja      SMALLINT NOT NULL,
    hari_hadir      SMALLINT NOT NULL,
    hari_sakit      SMALLINT NOT NULL,
    hari_izin       SMALLINT NOT NULL,
    hari_alfa       SMALLINT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    dibayar_at      TIMESTAMPTZ,
    transaksi_id    UUID REFERENCES transaksi_rekening(id),
    file_url        TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (kontrak_id, periode)
);

-- Fraud detection
CREATE TABLE fraud_rule (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id      UUID REFERENCES bmt(id),
    nama        VARCHAR(100) NOT NULL,
    tipe        VARCHAR(30) NOT NULL,
    kondisi     JSONB NOT NULL,
    aksi        VARCHAR(20) NOT NULL,
    is_aktif    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fraud_alert (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    rule_id         UUID NOT NULL REFERENCES fraud_rule(id),
    nasabah_id      UUID REFERENCES nasabah(id),
    transaksi_id    UUID,
    deskripsi       TEXT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    direview_oleh   UUID REFERENCES pengguna(id),
    direview_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sdm_kontrak_pegawai ON sdm_kontrak(pegawai_id, status);
CREATE INDEX idx_fraud_alert_bmt ON fraud_alert(bmt_id, status, created_at DESC);
