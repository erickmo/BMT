# Domain: E-commerce OPOP

## Model Bisnis

- **B2C:** Wali santri beli produk dari toko pondok anaknya
- **B2B (OPOP):** Pondok A beli produk dari Pondok B (antar pondok saling support)
- **Seller:** Pondok pesantren & BMT/koperasi pondok
- **Buyer:** Wali santri (B2C) atau pondok lain (B2B)
- **Pembayaran:** Midtrans + potong saldo rekening BMT + kartu NFC

---

## Toko & Produk

```sql
CREATE TABLE toko (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    nama            VARCHAR(255) NOT NULL,
    slug            VARCHAR(100) UNIQUE NOT NULL,
    deskripsi       TEXT,
    logo_url        TEXT,
    banner_url      TEXT,
    kategori_toko   VARCHAR(30) NOT NULL,  -- PONDOK | BMT_KOPERASI
    is_opop         BOOLEAN NOT NULL DEFAULT FALSE,
    -- Toko ini berpartisipasi di marketplace OPOP lintas pondok
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
    harga_b2b       BIGINT,                -- Harga khusus antar pondok (B2B)
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
```

---

## Pesanan

```sql
CREATE TABLE pesanan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_tipe      VARCHAR(20) NOT NULL,  -- WALI_SANTRI | PONDOK
    nasabah_id      UUID REFERENCES nasabah(id),        -- jika WALI_SANTRI
    bmt_buyer_id    UUID REFERENCES bmt(id),             -- jika PONDOK (B2B)
    toko_id         UUID NOT NULL REFERENCES toko(id),
    bmt_seller_id   UUID NOT NULL REFERENCES bmt(id),
    nomor_pesanan   VARCHAR(40) UNIQUE NOT NULL,
    status          VARCHAR(30) NOT NULL DEFAULT 'MENUNGGU_PEMBAYARAN',
    -- MENUNGGU_PEMBAYARAN | DIBAYAR | DIPROSES | DIKIRIM | SELESAI | DIBATALKAN
    subtotal        BIGINT NOT NULL,
    ongkir          BIGINT NOT NULL DEFAULT 0,
    total           BIGINT NOT NULL,
    alamat_kirim    JSONB NOT NULL,
    kurir           VARCHAR(50),
    nomor_resi      VARCHAR(100),
    metode_bayar    VARCHAR(30),
    -- MIDTRANS | REKENING_BMT | NFC
    catatan         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pesanan_item (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pesanan_id  UUID NOT NULL REFERENCES pesanan(id),
    produk_id   UUID NOT NULL REFERENCES produk(id),
    nama_produk VARCHAR(255) NOT NULL,   -- snapshot saat pesan
    harga       BIGINT NOT NULL,          -- snapshot saat pesan
    jumlah      INT NOT NULL CHECK (jumlah > 0),
    subtotal    BIGINT NOT NULL
);

CREATE TABLE pembayaran_pesanan (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pesanan_id          UUID NOT NULL REFERENCES pesanan(id),
    metode              VARCHAR(30) NOT NULL,
    nominal             BIGINT NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    -- PENDING | SETTLEMENT | EXPIRE | CANCEL
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
```

---

## Analytics OPOP

```sql
-- Snapshot harian (dihitung worker malam)
CREATE TABLE opop_analytics_harian (
    tanggal             DATE NOT NULL,
    toko_id             UUID NOT NULL REFERENCES toko(id),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    total_pesanan       INT NOT NULL DEFAULT 0,
    total_pendapatan    BIGINT NOT NULL DEFAULT 0,
    total_item_terjual  INT NOT NULL DEFAULT 0,
    pengunjung_unik     INT NOT NULL DEFAULT 0,
    PRIMARY KEY (tanggal, toko_id)
);
```

---

## Checklist Syariah E-commerce

- [ ] Harga produk OPOP: transparan, tidak ada penipuan (gharar)
- [ ] Komisi platform dari e-commerce: akad jelas (ujrah/wakalah)
