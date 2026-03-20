# 11 — Keamanan, Rate Limiting & Offline Mode

> **Terakhir diperbarui:** 20 Maret 2026

## 2FA (Feature: `KEAMANAN_2FA`)
```sql
otp_log (tujuan, channel,  -- SMS|EMAIL
         kode_hash VARCHAR(255),
         konteks,          -- LOGIN|RESET_PIN|TRANSAKSI_BESAR
         expired_at TIMESTAMPTZ,
         digunakan BOOLEAN, ip_address INET)
```
Semua threshold dari settings: `keamanan.otp_expired_menit`, `keamanan.lockout_menit`, dll.

## Session Management
```sql
sesi_aktif (pengguna_id | nasabah_id,
            device_id, device_info JSONB,
            ip_address INET,
            refresh_token_hash VARCHAR(255),
            expired_at TIMESTAMPTZ, last_active TIMESTAMPTZ)
```
Paksa logout semua device: `DELETE FROM sesi_aktif WHERE pengguna_id = $1`

## Audit Log Lengkap
```sql
audit_log (bmt_id,
           aktor_id, aktor_tipe,    -- STAF|NASABAH|SISTEM|DEVELOPER
           aksi VARCHAR(50),        -- CREATE|READ|UPDATE|DELETE|LOGIN|LOGOUT|APPROVE|REJECT
           entitas_tipe, entitas_id UUID,
           ip_address INET, user_agent TEXT,
           data_sebelum JSONB, data_sesudah JSONB,
           created_at TIMESTAMPTZ)

-- Index wajib
idx_audit_log_bmt_tgl ON audit_log(bmt_id, created_at)
idx_audit_log_aktor   ON audit_log(aktor_id, created_at)
idx_audit_log_entitas ON audit_log(entitas_tipe, entitas_id)
```

## Anti-Fraud Detection
```sql
fraud_rule (bmt_id,  -- NULL = global platform
            kode, nama,
            kondisi JSONB,  -- {"jenis_transaksi":"SETOR_TUNAI","nominal_min":50000000}
            aksi,           -- ALERT|BLOCK|REQUIRE_APPROVAL
            is_aktif BOOLEAN)

fraud_alert (bmt_id, rule_id, nasabah_id, transaksi_id,
             detail JSONB,
             status,         -- OPEN|REVIEWED|DISMISSED|ESCALATED
             ditangani_oleh, ditangani_at)
```
Worker `FraudDetection` berjalan **real-time per event transaksi**.

## Rate Limiting per BMT
```
# Redis sliding window per bmt_id
# Key: "rl:{bmt_id}:{menit}" → INCR + EXPIRE
# Default dari platform_settings.rate_limit_rpm
# Override per BMT: bmt_settings."rate_limit.rpm"
```
Middleware `ratelimit.go` membaca limit dari settings sebelum inject ke handler.

## Keamanan Umum
- JWT: access_token 15 menit, refresh_token 7 hari
- Password/PIN: bcrypt cost ≥ 12
- `is_rahasia = true` → ter-mask di log & response API
- HTTPS/TLS 1.3
- `/dev/*`: `Developer-Token` via env
- NFC kiosk: IP whitelist per terminal
- Biometrik template: disimpan encrypted

## Offline Mode (Teller & Absensi)

### Cakupan
- `app/teller` — transaksi tunai saat koneksi terputus
- Absensi NFC & manual saat koneksi terputus

### Strategi
```
Online  : request langsung ke API
Offline : simpan ke SQLite lokal (Drift/Isar) dengan status PENDING_SYNC
          → saat online → kirim ke API via SyncOfflineTransaksi worker
          → server validasi: saldo, idempotency_key, timestamp
          → konflik (saldo berubah saat offline) → flag CONFLICT → teller review
```

```sql
-- Tabel lokal SQLite (app/teller)
offline_queue (
    id TEXT PRIMARY KEY,          -- UUID lokal
    endpoint TEXT,                -- "/teller/rekening/:id/setor"
    payload TEXT,                 -- JSON body
    idempotency_key TEXT UNIQUE,
    status TEXT,                  -- PENDING|SYNCED|CONFLICT|FAILED
    error_detail TEXT,
    created_at INTEGER,           -- unix timestamp
    synced_at INTEGER
)
```

**Aturan offline:**
- Teller wajib login online minimal sekali per sesi
- Saldo offline = saldo terakhir sync + perhitungan lokal
- Konflik muncul di dashboard manajer untuk review

## Disaster Recovery

### Backup
```
PostgreSQL harian (02:00 WIB): pg_dump → gzip → MinIO bucket backup/
Retensi: 7 hari harian, 4 minggu mingguan, 12 bulan bulanan
Redis: RDB snapshot tiap 1 jam + AOF enabled
```

### RPO/RTO
- Transaksi keuangan: RPO maks 1 jam (backup + WAL)
- Data pondok: RPO maks 24 jam
- Target RTO: < 4 jam restore full database

### Maintenance Mode
```
POST /dev/maintenance {"aktif": true, "pesan": "..."}
→ /api/* dan /platform/* → 503
→ /dev/* dan /health → tetap jalan
→ Workers tetap berjalan
```
