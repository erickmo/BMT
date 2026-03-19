# Domain: Keamanan & Audit

## 2FA — OTP via SMS atau Email

```sql
CREATE TABLE otp_log (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tujuan       VARCHAR(255) NOT NULL,  -- no. HP atau email
    channel      VARCHAR(10) NOT NULL,  -- SMS | EMAIL
    kode_hash    VARCHAR(255) NOT NULL, -- bcrypt hash OTP 6 digit (TIDAK plaintext)
    tipe         VARCHAR(20) NOT NULL,
    -- LOGIN | RESET_PIN | KONFIRMASI_TRANSAKSI
    referensi_id UUID,
    is_digunakan BOOLEAN NOT NULL DEFAULT FALSE,
    expired_at   TIMESTAMPTZ NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Konfigurasi OTP dari settings:**
```
keamanan.otp_expired_menit     → "5"
keamanan.otp_maks_percobaan    → "3"
keamanan.2fa_wajib_staf        → "true"
keamanan.2fa_opsional_nasabah  → "true"
keamanan.lockout_menit         → "15"
keamanan.maks_gagal_login      → "5"
```

---

## Session Management

```sql
CREATE TABLE sesi_aktif (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subjek_id           UUID NOT NULL,
    subjek_tipe         VARCHAR(20) NOT NULL,
    -- NASABAH | PENGGUNA | PENGGUNA_PONDOK
    refresh_token_hash  VARCHAR(255) UNIQUE NOT NULL,
    device_info         JSONB,
    -- {platform, os_version, app_version, device_id}
    ip_address          INET,
    last_active_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expired_at          TIMESTAMPTZ NOT NULL,
    is_aktif            BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

- JWT: access token 15 menit + refresh token 7 hari (dari settings)
- Paksa logout semua device: `UPDATE sesi_aktif SET is_aktif = false WHERE subjek_id = ?`
- Max sesi aktif per subjek: `keamanan.max_sesi_aktif` (default: 3)
- Cleanup expired: Worker `CleanupSesiExpired` setiap jam

---

## Audit Log Lengkap

```sql
CREATE TABLE audit_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID REFERENCES bmt(id),
    -- NULL jika aksi level platform (developer)
    subjek_id       UUID NOT NULL,
    subjek_tipe     VARCHAR(20) NOT NULL,
    -- NASABAH | PENGGUNA | PENGGUNA_PONDOK | DEVELOPER
    aksi            VARCHAR(100) NOT NULL,
    -- "LOGIN", "UPDATE_NASABAH", "APPROVE_FORM", "POST_JURNAL", dll.
    resource_tipe   VARCHAR(50),
    resource_id     UUID,
    data_sebelum    JSONB,         -- snapshot sebelum perubahan
    data_sesudah    JSONB,         -- snapshot sesudah perubahan
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_log_subjek    ON audit_log(subjek_id, created_at DESC);
CREATE INDEX idx_audit_log_resource  ON audit_log(resource_tipe, resource_id, created_at DESC);
CREATE INDEX idx_audit_log_bmt       ON audit_log(bmt_id, created_at DESC);
```

Retensi: `keamanan.audit_log_retensi_hari` (default: 365 hari).
Worker `CleanupAuditLog` setiap Minggu 03:00 hapus log melebihi retensi.

---

## Anti-Fraud Detection

Rule-based, threshold dikonfigurasi per BMT di settings atau via tabel:

```sql
CREATE TABLE fraud_rule (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id      UUID REFERENCES bmt(id),   -- NULL = berlaku semua BMT
    nama        VARCHAR(100) NOT NULL,
    tipe        VARCHAR(30) NOT NULL,
    -- FREKUENSI | NOMINAL | LOKASI | WAKTU | VELOCITY
    kondisi     JSONB NOT NULL,
    -- {"max_transaksi_per_jam": 10, "nominal_min": 5000000}
    aksi        VARCHAR(20) NOT NULL,
    -- LOG | NOTIFIKASI | BLOKIR_SEMENTARA | REQUIRE_OTP
    is_aktif    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fraud_alert (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    rule_id         UUID NOT NULL REFERENCES fraud_rule(id),
    nasabah_id      UUID REFERENCES nasabah(id),
    transaksi_id    UUID,
    deskripsi       TEXT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    -- OPEN | REVIEWED | FALSE_POSITIVE | CONFIRMED
    direview_oleh   UUID REFERENCES pengguna(id),
    direview_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Contoh rules:**
- Lebih dari 10 transaksi dalam 1 jam dari rekening yang sama
- Transaksi NFC di luar jam operasional pondok
- Nominal transaksi > 5× rata-rata 30 hari terakhir
- Login dari IP/device baru langsung bertransaksi besar

Worker `FraudDetection` berjalan real-time (event-driven) setiap transaksi masuk.
