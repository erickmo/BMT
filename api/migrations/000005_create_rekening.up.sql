CREATE TABLE jenis_rekening (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id                  UUID NOT NULL REFERENCES bmt(id),
    kode                    VARCHAR(20) NOT NULL,
    nama                    VARCHAR(100) NOT NULL,
    tipe_dasar              VARCHAR(30) NOT NULL,
    akad                    VARCHAR(30) NOT NULL,
    deskripsi               TEXT,
    setoran_awal_min        BIGINT NOT NULL DEFAULT 0,
    setoran_min             BIGINT NOT NULL DEFAULT 0,
    bisa_ditarik            BOOLEAN NOT NULL DEFAULT TRUE,
    syarat_penarikan        TEXT,
    nisbah_nasabah          SMALLINT,
    jangka_hari             SMALLINT,
    biaya_admin_bulanan     BIGINT NOT NULL DEFAULT 0,
    bisa_nfc                BOOLEAN NOT NULL DEFAULT FALSE,
    bisa_autodebet          BOOLEAN NOT NULL DEFAULT TRUE,
    biaya_admin_buka        BIGINT NOT NULL DEFAULT 0,
    is_aktif                BOOLEAN NOT NULL DEFAULT TRUE,
    urutan_tampil           SMALLINT NOT NULL DEFAULT 0,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by              UUID NOT NULL REFERENCES pengguna(id),
    updated_by              UUID NOT NULL REFERENCES pengguna(id),
    UNIQUE (bmt_id, kode)
);

CREATE TABLE rekening (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    cabang_id           UUID NOT NULL REFERENCES cabang(id),
    nasabah_id          UUID NOT NULL REFERENCES nasabah(id),
    jenis_rekening_id   UUID NOT NULL REFERENCES jenis_rekening(id),
    nomor_rekening      VARCHAR(40) UNIQUE NOT NULL,
    saldo               BIGINT NOT NULL DEFAULT 0 CHECK (saldo >= 0),
    status              VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    alasan_blokir       TEXT,
    biaya_admin_bulanan BIGINT NOT NULL DEFAULT 0,
    nominal_deposito    BIGINT,
    nisbah_nasabah      SMALLINT,
    tanggal_buka        DATE NOT NULL,
    tanggal_jatuh_tempo DATE,
    tanggal_tutup       DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by_form_id  UUID REFERENCES form_pengajuan(id),
    updated_by_form_id  UUID REFERENCES form_pengajuan(id)
);

CREATE INDEX idx_rekening_nasabah ON rekening(nasabah_id, bmt_id);
CREATE INDEX idx_rekening_bmt ON rekening(bmt_id, cabang_id, status);

CREATE TABLE transaksi_rekening (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    rekening_id     UUID NOT NULL REFERENCES rekening(id),
    jenis           VARCHAR(30) NOT NULL,
    posisi          VARCHAR(10) NOT NULL CHECK (posisi IN ('DEBIT','KREDIT')),
    nominal         BIGINT NOT NULL CHECK (nominal > 0),
    saldo_sebelum   BIGINT NOT NULL,
    saldo_sesudah   BIGINT NOT NULL,
    keterangan      TEXT,
    referensi_id    UUID,
    referensi_tipe  VARCHAR(30),
    idempotency_key UUID UNIQUE,
    created_by      UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transaksi_rekening ON transaksi_rekening(rekening_id, created_at DESC);
CREATE INDEX idx_transaksi_bmt ON transaksi_rekening(bmt_id, created_at DESC);
