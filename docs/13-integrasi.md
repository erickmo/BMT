# 13 — Integrasi & White-label

> **Terakhir diperbarui:** 20 Maret 2026

## Midtrans

### Konfigurasi (per BMT, fallback ke platform)
```
bmt_settings: midtrans.server_key, midtrans.client_key, midtrans.env
platform_settings: platform.midtrans_server_key (fallback)
```

### Alur
```
POST /api/.../setor-online (X-Idempotency-Key)
→ Backend buat Snap (order_id = idempotency_key)
→ Return snap_token
→ Flutter tampilkan Snap UI (WebView)
→ POST /webhook/midtrans
    → Verifikasi SHA512 signature
    → SETTLEMENT → posting transaksi + jurnal + usage_log
    → Duplikat → 200 tanpa re-proses
```

```sql
midtrans_transaksi (bmt_id, cabang_id,
                    order_id UNIQUE,
                    referensi_id, referensi_tipe,
                    -- SIMPANAN|ANGSURAN|SPP|DONASI|PPDB|PESANAN|TIKET_EVENT|LISTING
                    nominal, status, payment_type,
                    snap_token, raw_notifikasi JSONB, settled_at)
```

Worker `CekMidtransPending` setiap 15 menit — poll Midtrans untuk transaksi PENDING > 30 menit.

## DAPODIK (Feature: `INTEGRASI_DAPODIK`)
Konfigurasi: `integrasi.dapodik_aktif`, `dapodik_npsn`, `dapodik_username`, `dapodik_password*`
Worker `SinkronDAPODIK` — jadwal dari `integrasi.sinkron_jadwal`

## EMIS Kemenag (Feature: `INTEGRASI_EMIS`)
Konfigurasi: `integrasi.emis_aktif`, `emis_nsm`, `emis_token*`
Worker `SinkronEMIS` — jadwal dari `integrasi.sinkron_jadwal`

```sql
sinkronisasi_eksternal (bmt_id, sumber,  -- DAPODIK|EMIS
                         jenis, status, jumlah_record, berhasil, gagal,
                         error_detail JSONB, dijalankan_oleh)
```

## White-label (Feature: `WHITELABEL`)

Konfigurasi di `bmt_settings`:
```
whitelabel.nama_app           → "Santri Pay"
whitelabel.bundle_id_android  → "com.annur.santripay"
whitelabel.bundle_id_ios      → "com.annur.santripay"
whitelabel.primary_color      → "#1B5E20"
whitelabel.logo_url           → "https://..."
whitelabel.splash_url         → "https://..."
```

**Implementasi:** CI/CD parameterized build per BMT.
Setiap BMT yang aktifkan white-label → build Flutter terpisah yang ambil asset dari settings saat build time (GitHub Actions matrix atau Codemagic multi-app).
