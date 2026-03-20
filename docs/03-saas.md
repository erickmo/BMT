# 03 — SaaS: Tier, Add-on, Portal Developer & Listing

> **Terakhir diperbarui:** 20 Maret 2026

## Model Bisnis SaaS

Platform menggunakan model **paket tier + add-on per fitur**:
- **Paket tier** — bundle fitur dasar dengan harga bulanan/tahunan
- **Add-on** — fitur à la carte di atas tier, beli per fitur
- **Listing stakeholder** — pendapatan terpisah dari fitur platform

## Paket Tier

```sql
CREATE TABLE saas_paket (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kode            VARCHAR(20) UNIQUE NOT NULL,
    -- FREE | BASIC | PRO | ENTERPRISE
    nama            VARCHAR(100) NOT NULL,
    deskripsi       TEXT,
    harga_bulanan   BIGINT NOT NULL DEFAULT 0,   -- 0 = gratis
    harga_tahunan   BIGINT NOT NULL DEFAULT 0,   -- diskon vs bulanan
    -- Batas kapasitas
    maks_cabang     SMALLINT NOT NULL DEFAULT 1,
    maks_nasabah    INT NOT NULL DEFAULT 100,     -- -1 = unlimited
    maks_storage_gb SMALLINT NOT NULL DEFAULT 1,
    -- Fitur yang termasuk (JSON array kode fitur)
    fitur_termasuk  JSONB NOT NULL DEFAULT '[]',
    -- ["CBS_DASAR","PONDOK_ADMINISTRASI","NOTIFIKASI_EMAIL"]
    is_aktif        BOOLEAN NOT NULL DEFAULT TRUE,
    urutan_tampil   SMALLINT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Contoh Tier (Dikonfigurasi Developer)

| Tier | Harga/bln | Cabang | Nasabah | Fitur Utama |
|------|-----------|--------|---------|-------------|
| FREE | Rp 0 | 1 | 100 | CBS dasar, admin santri |
| BASIC | Rp 299k | 1 | 500 | + Autodebet, NFC, OPOP |
| PRO | Rp 799k | 5 | 2000 | + Payroll, Raport, Listing |
| ENTERPRISE | Custom | Unlimited | Unlimited | Semua fitur + white-label |

## Add-on Fitur

Fitur yang bisa dibeli secara terpisah di atas paket tier:

```sql
CREATE TABLE saas_addon (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kode_fitur      VARCHAR(50) UNIQUE NOT NULL,
    nama            VARCHAR(100) NOT NULL,
    deskripsi       TEXT,
    kategori        VARCHAR(30) NOT NULL,
    -- CBS | PONDOK | ECOMMERCE | KEAMANAN | INTEGRASI | KOMUNIKASI
    harga_bulanan   BIGINT NOT NULL,
    harga_tahunan   BIGINT NOT NULL,
    -- Tier minimum yang bisa membeli add-on ini
    tier_minimum    VARCHAR(20) NOT NULL DEFAULT 'FREE',
    is_aktif        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Contoh Add-on

| Kode Fitur | Nama | Harga/bln | Tier Min |
|-----------|------|-----------|---------|
| `CBS_PEMBIAYAAN` | Modul Pembiayaan Lengkap | Rp 99k | FREE |
| `CBS_DEPOSITO` | Deposito Mudharabah | Rp 49k | BASIC |
| `PONDOK_RAPORT` | Raport Digital | Rp 49k | BASIC |
| `PONDOK_TAHFIDZ` | Manajemen Tahfidz | Rp 29k | FREE |
| `PONDOK_PPDB` | PPDB Online | Rp 79k | BASIC |
| `PONDOK_UKS` | Health Record UKS | Rp 29k | FREE |
| `PONDOK_PERPUS` | Perpustakaan Digital | Rp 39k | FREE |
| `ECOMMERCE_OPOP` | Marketplace OPOP | Rp 99k | BASIC |
| `KOMUNIKASI_WA` | WhatsApp Notifikasi | Rp 49k | FREE |
| `KOMUNIKASI_SMS` | SMS OTP & Reminder | Rp 29k | FREE |
| `INTEGRASI_DAPODIK` | Sinkron DAPODIK | Rp 79k | PRO |
| `INTEGRASI_EMIS` | Sinkron EMIS Kemenag | Rp 79k | PRO |
| `KEAMANAN_2FA` | 2FA untuk semua user | Rp 29k | FREE |
| `WHITELABEL` | Custom Branding App | Rp 499k | ENTERPRISE |
| `SDM_PAYROLL` | Payroll & Slip Gaji | Rp 99k | PRO |
| `LISTING_AKSES` | Akses direktori listing | Rp 49k | BASIC |

## Feature Gate

Setiap fitur ber-kode dicek via middleware sebelum handler:

```go
// pkg/featuregate/gate.go
func (g *FeatureGate) IsEnabled(ctx context.Context, bmtID uuid.UUID, kodefitur string) bool {
    // 1. Ambil paket tier BMT
    // 2. Cek apakah kode_fitur ada di saas_paket.fitur_termasuk
    // 3. Jika tidak → cek saas_bmt_addon (add-on yang dibeli BMT ini)
    // 4. Jika tidak → return false
}

// Middleware
func RequireFeature(kode string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            bmtID := TenantFromCtx(r.Context()).BMTID
            if !featureGate.IsEnabled(r.Context(), bmtID, kode) {
                response.Error(w, 403, "FITUR_TIDAK_AKTIF",
                    "Fitur ini tidak tersedia di paket Anda")
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// Penggunaan di router
r.With(RequireFeature("CBS_DEPOSITO")).Post("/simpanan/deposito", ...)
r.With(RequireFeature("ECOMMERCE_OPOP")).Get("/opop/toko", ...)
```

## Langganan BMT

```sql
CREATE TABLE saas_bmt_langganan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    paket_id        UUID NOT NULL REFERENCES saas_paket(id),
    siklus          VARCHAR(10) NOT NULL DEFAULT 'BULANAN',
    -- BULANAN | TAHUNAN
    harga_disepakati BIGINT NOT NULL, -- bisa berbeda dari harga paket (custom)
    tanggal_mulai   DATE NOT NULL,
    tanggal_selesai DATE,             -- NULL = tidak terbatas
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    -- AKTIF | SUSPENDED | EXPIRED | CANCELLED
    auto_renew      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add-on yang dibeli per BMT
CREATE TABLE saas_bmt_addon (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    addon_id        UUID NOT NULL REFERENCES saas_addon(id),
    siklus          VARCHAR(10) NOT NULL DEFAULT 'BULANAN',
    harga_disepakati BIGINT NOT NULL,
    tanggal_mulai   DATE NOT NULL,
    tanggal_selesai DATE,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    UNIQUE (bmt_id, addon_id)
);

-- Invoice tagihan SaaS ke BMT
CREATE TABLE saas_invoice (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    nomor_invoice   VARCHAR(30) UNIQUE NOT NULL,
    periode         CHAR(7) NOT NULL,   -- "2025-01"
    -- Item tagihan
    item            JSONB NOT NULL,
    -- [{nama, kode, harga, siklus}, ...]
    subtotal        BIGINT NOT NULL,
    total           BIGINT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'UNPAID',
    -- UNPAID | PAID | OVERDUE | CANCELLED
    jatuh_tempo     DATE NOT NULL,
    dibayar_at      TIMESTAMPTZ,
    midtrans_order_id VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## Portal Developer — Manajemen BMT & SaaS

### Alur Setup BMT Baru
```
1. POST /dev/bmt         → input data BMT
2. POST /dev/bmt/:id/langganan  → pilih paket tier + add-on
3. POST /dev/bmt/:id/cabang     → buat minimal 1 cabang
4. POST /dev/bmt/:id/seed       → buat akun ADMIN_BMT pertama
   → sistem kirim email undangan ke management BMT
5. (Opsional) PUT /dev/bmt/:id/custom-price → override harga per kontrak
```

### Endpoint Developer SaaS
```
# Paket & Add-on
GET|POST|PUT /dev/saas/paket / /:id
GET|POST|PUT /dev/saas/addon / /:id

# Manajemen BMT
GET|POST|PUT /dev/bmt / /:id / /:id/status
POST         /dev/bmt/:id/langganan
PUT          /dev/bmt/:id/langganan/:id   # Update/upgrade/downgrade
POST         /dev/bmt/:id/addon          # Tambah add-on
DELETE       /dev/bmt/:id/addon/:id      # Hapus add-on
POST         /dev/bmt/:id/cabang
POST         /dev/bmt/:id/seed

# Invoice & Billing
GET  /dev/invoice                        # List semua invoice
GET  /dev/invoice/:bmt_id               # Invoice per BMT
POST /dev/invoice/generate              # Generate invoice bulanan manual
GET  /dev/revenue                       # Dashboard revenue MRR/ARR

# Listing
GET  /dev/listing/pendaftar             # Antrian yang menunggu approve
POST /dev/listing/pendaftar/:id/approve
POST /dev/listing/pendaftar/:id/reject
GET|POST|PUT /dev/listing/kategori / /:id
```

---

## Listing Stakeholder

### Konsep
Direktori layanan sekitar pondok pesantren. Stakeholder mendaftar sendiri via form publik, developer mereview dan approve.

Wali santri bisa **lihat dan kontak** listing di `app/nasabah` — tidak ada transaksi dalam app (kontak via telepon/WA).

### Kategori Listing (Dikonfigurasi Developer, Open-ended)
```sql
CREATE TABLE listing_kategori (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kode        VARCHAR(30) UNIQUE NOT NULL,
    nama        VARCHAR(100) NOT NULL,
    icon_url    TEXT,
    urutan      SMALLINT NOT NULL DEFAULT 0,
    is_aktif    BOOLEAN NOT NULL DEFAULT TRUE
);
-- Contoh data: GURU_LES, BIMBEL, ANTAR_JEMPUT, CATERING,
--              EKSTRA_SANGGAR, KESEHATAN, LAINNYA
```

### Pendaftaran Listing (Self-register)
```sql
CREATE TABLE listing_pendaftaran (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Identitas pendaftar
    nama_usaha          VARCHAR(255) NOT NULL,
    nama_pic            VARCHAR(255) NOT NULL,
    telepon             VARCHAR(20) NOT NULL,
    email               VARCHAR(255),
    kategori_id         UUID NOT NULL REFERENCES listing_kategori(id),
    -- Lokasi
    alamat              TEXT NOT NULL,
    kota                VARCHAR(100) NOT NULL,
    provinsi            VARCHAR(100) NOT NULL,
    lat                 NUMERIC(10,7),
    lng                 NUMERIC(10,7),
    -- Konten listing
    deskripsi           TEXT NOT NULL,
    foto_urls           JSONB DEFAULT '[]',
    -- Maksimal 5 foto
    -- Target pondok (bisa lebih dari 1)
    target_bmt_ids      JSONB DEFAULT '[]',
    -- UUID array BMT yang mau di-listing-kan
    -- Status review
    status              VARCHAR(20) NOT NULL DEFAULT 'MENUNGGU',
    -- MENUNGGU | DISETUJUI | DITOLAK | SUSPENDED
    catatan_developer   TEXT,
    diproses_oleh       TEXT,            -- identifier developer
    diproses_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Listing Aktif
```sql
CREATE TABLE listing (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pendaftaran_id      UUID NOT NULL REFERENCES listing_pendaftaran(id),
    kategori_id         UUID NOT NULL REFERENCES listing_kategori(id),
    nama_usaha          VARCHAR(255) NOT NULL,
    nama_pic            VARCHAR(255) NOT NULL,
    telepon             VARCHAR(20) NOT NULL,
    email               VARCHAR(255),
    deskripsi           TEXT NOT NULL,
    foto_urls           JSONB DEFAULT '[]',
    alamat              TEXT NOT NULL,
    kota                VARCHAR(100) NOT NULL,
    lat                 NUMERIC(10,7),
    lng                 NUMERIC(10,7),
    -- Visibilitas per BMT
    -- Listing bisa tampil di beberapa BMT sekaligus
    bmt_ids             JSONB NOT NULL DEFAULT '[]',
    -- Paket langganan listing
    paket_listing       VARCHAR(20) NOT NULL DEFAULT 'BASIC',
    -- BASIC | PREMIUM
    -- PREMIUM: tampil di atas + badge verified + highlight
    langganan_aktif_sampai DATE,
    is_verified         BOOLEAN NOT NULL DEFAULT FALSE,
    status              VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    -- AKTIF | SUSPEND | EXPIRED
    rating_avg          NUMERIC(3,2),
    total_ulasan        INT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Ulasan listing dari wali santri
CREATE TABLE listing_ulasan (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id  UUID NOT NULL REFERENCES listing(id),
    nasabah_id  UUID NOT NULL REFERENCES nasabah(id),
    rating      SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    komentar    TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (listing_id, nasabah_id)
);

-- Langganan & invoice listing
CREATE TABLE listing_langganan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id      UUID NOT NULL REFERENCES listing(id),
    paket           VARCHAR(20) NOT NULL,   -- BASIC | PREMIUM
    siklus          VARCHAR(10) NOT NULL,   -- BULANAN | TAHUNAN
    harga           BIGINT NOT NULL,
    tanggal_mulai   DATE NOT NULL,
    tanggal_selesai DATE NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    midtrans_order_id VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Harga Listing (Dikonfigurasi Developer di platform_settings)
```
listing.harga_basic_bulanan    → "99000"   -- Rp 99k/bulan
listing.harga_basic_tahunan    → "990000"  -- Rp 990k/tahun (2 bln gratis)
listing.harga_premium_bulanan  → "249000"
listing.harga_premium_tahunan  → "2490000"
listing.maks_foto              → "5"
listing.radius_default_km      → "10"      -- radius cari dari lokasi pondok
```

### Endpoint Listing
```
# Publik (tanpa auth) — form pendaftaran
POST /listing/daftar             # Form self-register stakeholder
GET  /listing/status/:id         # Cek status pendaftaran

# Nasabah App (JWT nasabah)
GET  /nasabah/listing            # ?kategori=&bmt_id=&radius=
GET  /nasabah/listing/:id        # Detail listing
POST /nasabah/listing/:id/ulasan # Beri ulasan
GET  /nasabah/listing/kategori   # List kategori aktif

# Listing owner (JWT listing)
POST /auth/listing/login|refresh
GET|PUT /listing/profil           # Update info listing
POST    /listing/langganan        # Upgrade ke premium
GET     /listing/ulasan           # Lihat ulasan dari wali

# Developer
GET  /dev/listing/pendaftar       # Antrian approve
POST /dev/listing/pendaftar/:id/approve|reject
GET|POST|PUT /dev/listing/kategori / /:id
GET  /dev/listing/aktif           # Semua listing aktif
PUT  /dev/listing/:id/suspend     # Suspend listing
GET  /dev/listing/revenue         # Revenue dari listing
```

### Tampilan di `app/nasabah`
- Tab "Layanan" di dashboard wali santri
- Filter by kategori (ikon grid)
- Cari by nama/lokasi
- Sort: Terdekat / Rating tertinggi / Premium dulu
- Card listing: nama, kategori, rating, jarak dari pondok, badge PREMIUM jika applicable
- Detail: deskripsi, foto, kontak (klik → buka WA/telepon), ulasan
- Listing PREMIUM tampil di posisi atas dan memiliki highlight visual
