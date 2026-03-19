CREATE TABLE nasabah (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    cabang_id           UUID NOT NULL REFERENCES cabang(id),
    nomor_nasabah       VARCHAR(30) UNIQUE NOT NULL,
    nik                 VARCHAR(16),
    nama_lengkap        VARCHAR(255) NOT NULL,
    tempat_lahir        VARCHAR(100),
    tanggal_lahir       DATE,
    jenis_kelamin       CHAR(1),
    alamat              TEXT,
    telepon             VARCHAR(20),
    email               VARCHAR(255),
    foto_url            TEXT,
    pekerjaan           VARCHAR(100),
    status              VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    -- Auth
    pin_hash            VARCHAR(255),
    password_hash       VARCHAR(255),
    last_login_at       TIMESTAMPTZ,
    -- Audit
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_nasabah_nomor ON nasabah(bmt_id, nomor_nasabah);
CREATE INDEX idx_nasabah_bmt ON nasabah(bmt_id, status);

CREATE TABLE kartu_nfc (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    nasabah_id          UUID NOT NULL REFERENCES nasabah(id),
    uid                 VARCHAR(50) UNIQUE NOT NULL,
    pin_hash            VARCHAR(255) NOT NULL,
    limit_per_transaksi BIGINT NOT NULL DEFAULT 500000,
    limit_harian        BIGINT NOT NULL DEFAULT 2000000,
    saldo_nfc           BIGINT NOT NULL DEFAULT 0,
    status              VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    expired_at          DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_kartu_nfc_nasabah ON kartu_nfc(nasabah_id);
