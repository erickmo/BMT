CREATE TABLE rekening_autodebet_config (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    rekening_id         UUID NOT NULL REFERENCES rekening(id),
    jenis               VARCHAR(30) NOT NULL,
    tanggal_debet       SMALLINT NOT NULL CHECK (tanggal_debet BETWEEN 1 AND 28),
    is_aktif            BOOLEAN NOT NULL DEFAULT TRUE,
    referensi_id        UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID NOT NULL REFERENCES pengguna(id)
);

CREATE TABLE jadwal_autodebet (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    rekening_id         UUID NOT NULL REFERENCES rekening(id),
    config_id           UUID NOT NULL REFERENCES rekening_autodebet_config(id),
    jenis               VARCHAR(30) NOT NULL,
    nominal_target      BIGINT NOT NULL,
    tanggal_jatuh_tempo DATE NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'MENUNGGU',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tunggakan_autodebet (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    rekening_id         UUID NOT NULL REFERENCES rekening(id),
    jadwal_id           UUID NOT NULL REFERENCES jadwal_autodebet(id),
    jenis               VARCHAR(30) NOT NULL,
    nominal_target      BIGINT NOT NULL,
    nominal_terbayar    BIGINT NOT NULL DEFAULT 0,
    nominal_sisa        BIGINT NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'OUTSTANDING',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_jadwal_autodebet_bmt ON jadwal_autodebet(bmt_id, tanggal_jatuh_tempo, status);
CREATE INDEX idx_tunggakan_rekening ON tunggakan_autodebet(rekening_id, status);
