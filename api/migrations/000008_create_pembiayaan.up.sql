CREATE TABLE produk_pembiayaan (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    kode                VARCHAR(20) NOT NULL,
    nama                VARCHAR(100) NOT NULL,
    akad                VARCHAR(20) NOT NULL,
    margin_min          NUMERIC(5,2),
    margin_maks         NUMERIC(5,2),
    jangka_min          SMALLINT,
    jangka_maks         SMALLINT,
    plafon_min          BIGINT,
    plafon_maks         BIGINT,
    is_aktif            BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (bmt_id, kode)
);

CREATE TABLE pembiayaan (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id                  UUID NOT NULL REFERENCES bmt(id),
    cabang_id               UUID NOT NULL REFERENCES cabang(id),
    nasabah_id              UUID NOT NULL REFERENCES nasabah(id),
    produk_pembiayaan_id    UUID NOT NULL REFERENCES produk_pembiayaan(id),
    nomor_pembiayaan        VARCHAR(40) UNIQUE NOT NULL,
    akad                    VARCHAR(20) NOT NULL,
    pokok                   BIGINT NOT NULL,
    margin_persen           NUMERIC(5,2),
    nisbah_nasabah          SMALLINT,
    jangka_bulan            SMALLINT NOT NULL,
    angsuran_per_bulan      BIGINT NOT NULL,
    total_kewajiban         BIGINT NOT NULL,
    ada_beasiswa            BOOLEAN NOT NULL DEFAULT FALSE,
    beasiswa_persen         NUMERIC(5,2),
    beasiswa_nominal        BIGINT NOT NULL DEFAULT 0,
    beasiswa_sumber         VARCHAR(100),
    beasiswa_ditetapkan_oleh UUID REFERENCES pengguna_pondok(id),
    beasiswa_ditetapkan_at  TIMESTAMPTZ,
    status                  VARCHAR(20) NOT NULL DEFAULT 'PENGAJUAN',
    kolektibilitas          SMALLINT NOT NULL DEFAULT 1,
    hari_tunggak            INT NOT NULL DEFAULT 0,
    saldo_pokok             BIGINT NOT NULL,
    saldo_margin            BIGINT NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by              UUID NOT NULL REFERENCES pengguna(id),
    updated_by              UUID NOT NULL REFERENCES pengguna(id),
    is_voided               BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE angsuran_pembiayaan (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    pembiayaan_id       UUID NOT NULL REFERENCES pembiayaan(id),
    ke                  SMALLINT NOT NULL,
    tanggal_jatuh_tempo DATE NOT NULL,
    nominal_pokok       BIGINT NOT NULL,
    nominal_margin      BIGINT NOT NULL,
    nominal_total       BIGINT NOT NULL,
    nominal_terbayar    BIGINT NOT NULL DEFAULT 0,
    status              VARCHAR(20) NOT NULL DEFAULT 'BELUM_BAYAR',
    tanggal_bayar       DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE beasiswa_riwayat (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pembiayaan_id       UUID NOT NULL REFERENCES pembiayaan(id),
    persen_sebelum      NUMERIC(5,2),
    persen_sesudah      NUMERIC(5,2) NOT NULL,
    nominal_sebelum     BIGINT,
    nominal_sesudah     BIGINT NOT NULL,
    alasan              TEXT,
    ditetapkan_oleh     UUID NOT NULL REFERENCES pengguna_pondok(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pembiayaan_nasabah ON pembiayaan(nasabah_id, bmt_id);
CREATE INDEX idx_pembiayaan_bmt ON pembiayaan(bmt_id, cabang_id, status);
CREATE INDEX idx_angsuran_pembiayaan ON angsuran_pembiayaan(pembiayaan_id, ke);
