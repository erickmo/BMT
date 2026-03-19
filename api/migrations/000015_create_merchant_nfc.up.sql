CREATE TABLE merchant (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    nama            VARCHAR(255) NOT NULL,
    kategori        VARCHAR(50),
    alamat          TEXT,
    logo_url        TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE terminal_kiosk (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    nama            VARCHAR(100) NOT NULL,
    ip_address      INET NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE transaksi_nfc (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    kartu_nfc_id    UUID NOT NULL REFERENCES kartu_nfc(id),
    merchant_id     UUID REFERENCES merchant(id),
    nominal         BIGINT NOT NULL CHECK (nominal > 0),
    saldo_sebelum   BIGINT NOT NULL,
    saldo_sesudah   BIGINT NOT NULL,
    jenis           VARCHAR(20) NOT NULL,
    keterangan      TEXT,
    idempotency_key UUID UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transaksi_nfc_kartu ON transaksi_nfc(kartu_nfc_id, created_at DESC);
CREATE INDEX idx_terminal_kiosk_ip ON terminal_kiosk(ip_address);

-- Analytics
CREATE TABLE analytics_snapshot (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID REFERENCES cabang(id),
    tanggal         DATE NOT NULL,
    metrik          VARCHAR(50) NOT NULL,
    nilai           NUMERIC NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (bmt_id, cabang_id, tanggal, metrik)
);

CREATE INDEX idx_analytics_snapshot ON analytics_snapshot(bmt_id, tanggal DESC);

-- Monetisasi
CREATE TABLE opop_komisi (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pesanan_id      UUID NOT NULL REFERENCES pesanan(id) UNIQUE,
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    nilai_pesanan   BIGINT NOT NULL,
    persen_komisi   NUMERIC(5,2) NOT NULL,
    nominal_komisi  BIGINT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    periode         CHAR(7) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE konsultasi_sesi (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    penanya_id      UUID NOT NULL,
    penanya_tipe    VARCHAR(20) NOT NULL,
    penjawab_id     UUID REFERENCES pengguna_pondok(id),
    topik           VARCHAR(30) NOT NULL,
    judul           VARCHAR(255) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE konsultasi_pesan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sesi_id         UUID NOT NULL REFERENCES konsultasi_sesi(id),
    pengirim_id     UUID NOT NULL,
    pengirim_tipe   VARCHAR(20) NOT NULL,
    pesan           TEXT NOT NULL,
    file_url        TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_konsultasi_sesi_bmt ON konsultasi_sesi(bmt_id, status);
