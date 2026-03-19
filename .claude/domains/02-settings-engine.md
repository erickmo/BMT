# Domain: Settings Engine

## Prinsip

**Semua nilai yang bisa berbeda antar BMT/pondok/cabang wajib disimpan sebagai settings.**
Tidak ada nilai konfigurasi yang hardcode di kode Go maupun Flutter.

---

## Tabel Settings

```sql
-- Settings platform (developer)
CREATE TABLE platform_settings (
    kunci           VARCHAR(150) PRIMARY KEY,
    nilai           TEXT NOT NULL,
    tipe            VARCHAR(20) NOT NULL DEFAULT 'string',
    -- string | int | bool | json | float | date | time
    deskripsi       TEXT,
    is_rahasia      BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      TEXT NOT NULL
);

-- Settings BMT (management BMT)
CREATE TABLE bmt_settings (
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    kunci           VARCHAR(150) NOT NULL,
    nilai           TEXT NOT NULL,
    tipe            VARCHAR(20) NOT NULL DEFAULT 'string',
    is_locked       BOOLEAN NOT NULL DEFAULT FALSE,
    -- Jika true, cabang tidak bisa override
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID NOT NULL REFERENCES pengguna(id),
    PRIMARY KEY (bmt_id, kunci)
);

-- Settings cabang (manajer cabang)
CREATE TABLE cabang_settings (
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    kunci           VARCHAR(150) NOT NULL,
    nilai           TEXT NOT NULL,
    tipe            VARCHAR(20) NOT NULL DEFAULT 'string',
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID NOT NULL REFERENCES pengguna(id),
    PRIMARY KEY (cabang_id, kunci)
);
```

---

## Settings Resolver

```go
// pkg/settings/resolver.go
// Prioritas: cabang_settings > bmt_settings > platform_settings
func Resolve(ctx context.Context, bmtID, cabangID uuid.UUID, kunci string) string
func ResolveJSON(ctx context.Context, bmtID, cabangID uuid.UUID, kunci string) []string

// Semua service menggunakan ini — TIDAK PERNAH membaca hardcode
jam := settings.Resolve(ctx, bmtID, cabangID, "operasional.jam_buka")
// → "08:00" (dari DB, bisa beda per cabang)

metodeDiizinkan := settings.ResolveJSON(ctx, bmtID, cabangID, "pondok.absensi_metode")
// → ["MANUAL", "NFC", "BIOMETRIK"]
```

---

## Contoh Settings Lengkap

```jsonc
// ── platform_settings (dikelola developer) ──────────────────────────────────
"pecahan_uang.sumber"               : "DB"
"platform.midtrans_server_key"      : "..."
"platform.midtrans_env"             : "sandbox"
"platform.maintenance_mode"         : "false"
"platform.min_app_version.nasabah"  : "2.0.0"
"platform.min_app_version.teller"   : "1.5.0"
"platform.rate_limit_rpm"           : "300"
"analytics.dashboard_cache_ttl_detik" : "30"
"laporan.default_format"            : "PDF"

// ── bmt_settings (dikelola management BMT) ──────────────────────────────────
"operasional.jam_buka"              : "08:00"
"operasional.jam_tutup"             : "16:00"
"operasional.hari_kerja"            : "[1,2,3,4,5]"
"operasional.zona_waktu"            : "Asia/Jakarta"

"autodebet.retry_hari"              : "3"
"autodebet.jam_eksekusi"            : "07:00"
"autodebet.tanggal_simpanan_wajib"  : "1"

"sesi_teller.toleransi_selisih"     : "0"

"approval.FORM_DAFTAR_NASABAH"      : "[\"TELLER\",\"MANAJER_CABANG\"]"
"approval.FORM_BUKA_REKENING"       : "[\"TELLER\",\"MANAJER_CABANG\"]"
"approval.FORM_TUTUP_REKENING"      : "[\"MANAJER_CABANG\"]"
"approval.FORM_BLOKIR_REKENING"     : "[\"MANAJER_CABANG\"]"
"approval.FORM_BUKA_PEMBIAYAAN"     : "[\"KOMITE\"]"

"midtrans.server_key"               : "..."          // is_rahasia = true
"midtrans.client_key"               : "..."
"midtrans.env"                      : "production"
"midtrans.enabled_methods"          : "[\"gopay\",\"qris\",\"va_bni\"]"

"notifikasi.email_smtp_host"        : "smtp.example.com"
"notifikasi.email_smtp_port"        : "587"
"notifikasi.email_smtp_user"        : "bmt@example.com"
"notifikasi.email_smtp_pass"        : "..."          // is_rahasia = true
"notifikasi.wa_provider"            : "fonnte"       // fonnte | wablas | whacenter
"notifikasi.wa_token"               : "..."          // is_rahasia = true
"notifikasi.sms_provider"           : "zenziva"
"notifikasi.sms_apikey"             : "..."          // is_rahasia = true
"notifikasi.fcm_server_key"         : "..."          // is_rahasia = true (di platform_settings)

"nfc.limit_default_per_transaksi"   : "500000"
"nfc.limit_default_harian"          : "2000000"

"pondok.absensi_metode"             : "[\"MANUAL\",\"NFC\",\"BIOMETRIK\"]"
"pondok.reminder_spp_hari_sebelum"  : "3"

"ecommerce.komisi_persen"           : "2.5"
"ecommerce.opop_aktif"              : "true"

"keamanan.otp_expired_menit"        : "5"
"keamanan.otp_maks_percobaan"       : "3"
"keamanan.2fa_wajib_staf"           : "true"
"keamanan.2fa_opsional_nasabah"     : "true"
"keamanan.lockout_menit"            : "15"
"keamanan.maks_gagal_login"         : "5"
"keamanan.max_sesi_aktif"           : "3"
"keamanan.refresh_token_expired_hari" : "7"
"keamanan.audit_log_retensi_hari"   : "365"

"sdm.tanggal_gajian"                : "25"

"integrasi.dapodik_aktif"           : "true"
"integrasi.dapodik_npsn"            : "20572345"
"integrasi.dapodik_username"        : "..."
"integrasi.dapodik_password"        : "..."          // is_rahasia = true
"integrasi.emis_aktif"              : "true"
"integrasi.emis_nsm"                : "21237..."
"integrasi.emis_token"              : "..."          // is_rahasia = true
"integrasi.sinkron_jadwal"          : "MINGGUAN"

// ── whitelabel (per BMT) ─────────────────────────────────────────────────────
"whitelabel.nama_app"               : "Santri Pay"
"whitelabel.bundle_id_android"      : "com.annur.santripay"
"whitelabel.bundle_id_ios"          : "com.annur.santripay"
"whitelabel.primary_color"          : "#1B5E20"
"whitelabel.logo_url"               : "https://..."
"whitelabel.splash_url"             : "https://..."
```

---

## Pecahan Uang Rupiah (Data, Bukan Konstanta)

```sql
CREATE TABLE pecahan_uang (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nominal       BIGINT NOT NULL,
    jenis         VARCHAR(10) NOT NULL,       -- LOGAM | KERTAS
    label         VARCHAR(30) NOT NULL,       -- "Rp 1.000 (logam)"
    is_aktif      BOOLEAN NOT NULL DEFAULT TRUE,
    urutan        SMALLINT NOT NULL,
    berlaku_sejak DATE NOT NULL,
    ditarik_pada  DATE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (nominal, jenis)
);
```

Teller app mengambil pecahan aktif dari DB saat membuka sesi — **tidak ada array hardcode di kode**.

```go
// ✅ BENAR
pecahans, _ := pecahanRepo.GetAktif(ctx)

// ❌ SALAH
pecahans := []Pecahan{{100, "LOGAM"}, {200, "LOGAM"}, ...}
```
