CREATE TABLE pengguna (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID REFERENCES cabang(id),
    username        VARCHAR(100) UNIQUE NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    nama_lengkap    VARCHAR(255) NOT NULL,
    telepon         VARCHAR(20),
    role            VARCHAR(50) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pengguna_pondok (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    username        VARCHAR(100) UNIQUE NOT NULL,
    email           VARCHAR(255),
    password_hash   VARCHAR(255) NOT NULL,
    nama_lengkap    VARCHAR(255) NOT NULL,
    telepon         VARCHAR(20),
    role            VARCHAR(50) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sesi_aktif (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subjek_id           UUID NOT NULL,
    subjek_tipe         VARCHAR(20) NOT NULL,
    refresh_token_hash  VARCHAR(255) UNIQUE NOT NULL,
    device_info         JSONB,
    ip_address          INET,
    last_active_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expired_at          TIMESTAMPTZ NOT NULL,
    is_aktif            BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sesi_aktif_subjek ON sesi_aktif(subjek_id, is_aktif);

CREATE TABLE otp_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tujuan          VARCHAR(255) NOT NULL,
    channel         VARCHAR(10) NOT NULL,
    kode_hash       VARCHAR(255) NOT NULL,
    tipe            VARCHAR(20) NOT NULL,
    referensi_id    UUID,
    is_digunakan    BOOLEAN NOT NULL DEFAULT FALSE,
    expired_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID REFERENCES bmt(id),
    subjek_id       UUID NOT NULL,
    subjek_tipe     VARCHAR(20) NOT NULL,
    aksi            VARCHAR(100) NOT NULL,
    resource_tipe   VARCHAR(50),
    resource_id     UUID,
    data_sebelum    JSONB,
    data_sesudah    JSONB,
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_log_subjek    ON audit_log(subjek_id, created_at DESC);
CREATE INDEX idx_audit_log_resource  ON audit_log(resource_tipe, resource_id, created_at DESC);
CREATE INDEX idx_audit_log_bmt       ON audit_log(bmt_id, created_at DESC);
