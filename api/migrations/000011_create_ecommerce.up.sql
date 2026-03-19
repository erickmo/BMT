CREATE TABLE toko (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    nama            VARCHAR(255) NOT NULL,
    slug            VARCHAR(100) UNIQUE NOT NULL,
    deskripsi       TEXT,
    logo_url        TEXT,
    banner_url      TEXT,
    kategori_toko   VARCHAR(30) NOT NULL,
    is_opop         BOOLEAN NOT NULL DEFAULT FALSE,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    rating          NUMERIC(3,2),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE produk (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    toko_id         UUID NOT NULL REFERENCES toko(id),
    nama            VARCHAR(255) NOT NULL,
    slug            VARCHAR(100) NOT NULL,
    deskripsi       TEXT,
    kategori        VARCHAR(50) NOT NULL,
    harga           BIGINT NOT NULL CHECK (harga > 0),
    harga_b2b       BIGINT,
    stok            INT NOT NULL DEFAULT 0,
    satuan          VARCHAR(20) NOT NULL DEFAULT 'pcs',
    berat_gram      INT,
    foto_urls       JSONB NOT NULL DEFAULT '[]',
    is_opop         BOOLEAN NOT NULL DEFAULT FALSE,
    is_aktif        BOOLEAN NOT NULL DEFAULT TRUE,
    rating          NUMERIC(3,2),
    total_terjual   INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (toko_id, slug)
);

CREATE TABLE pesanan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_tipe      VARCHAR(20) NOT NULL,
    nasabah_id      UUID REFERENCES nasabah(id),
    bmt_buyer_id    UUID REFERENCES bmt(id),
    toko_id         UUID NOT NULL REFERENCES toko(id),
    bmt_seller_id   UUID NOT NULL REFERENCES bmt(id),
    nomor_pesanan   VARCHAR(40) UNIQUE NOT NULL,
    status          VARCHAR(30) NOT NULL DEFAULT 'MENUNGGU_PEMBAYARAN',
    subtotal        BIGINT NOT NULL,
    ongkir          BIGINT NOT NULL DEFAULT 0,
    total           BIGINT NOT NULL,
    alamat_kirim    JSONB NOT NULL,
    kurir           VARCHAR(50),
    nomor_resi      VARCHAR(100),
    metode_bayar    VARCHAR(30),
    catatan         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pesanan_item (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pesanan_id  UUID NOT NULL REFERENCES pesanan(id),
    produk_id   UUID NOT NULL REFERENCES produk(id),
    nama_produk VARCHAR(255) NOT NULL,
    harga       BIGINT NOT NULL,
    jumlah      INT NOT NULL CHECK (jumlah > 0),
    subtotal    BIGINT NOT NULL
);

CREATE TABLE pembayaran_pesanan (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pesanan_id          UUID NOT NULL REFERENCES pesanan(id),
    metode              VARCHAR(30) NOT NULL,
    nominal             BIGINT NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    midtrans_order_id   VARCHAR(100) UNIQUE,
    rekening_id         UUID REFERENCES rekening(id),
    kartu_nfc_id        UUID REFERENCES kartu_nfc(id),
    idempotency_key     UUID UNIQUE,
    settled_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE ulasan_produk (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    produk_id   UUID NOT NULL REFERENCES produk(id),
    pesanan_id  UUID NOT NULL REFERENCES pesanan(id),
    nasabah_id  UUID NOT NULL REFERENCES nasabah(id),
    rating      SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    komentar    TEXT,
    foto_urls   JSONB DEFAULT '[]',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (pesanan_id, produk_id)
);

CREATE INDEX idx_produk_toko ON produk(toko_id, is_aktif);
CREATE INDEX idx_produk_opop ON produk(is_opop, is_aktif);
CREATE INDEX idx_pesanan_nasabah ON pesanan(nasabah_id, created_at DESC);
CREATE INDEX idx_pesanan_toko ON pesanan(toko_id, status);
