# 04 — Settings Engine

> **Terakhir diperbarui:** 20 Maret 2026

## Prinsip
**Semua nilai yang bisa berbeda antar BMT/pondok/cabang wajib disimpan sebagai settings di DB.**
Tidak ada konstanta konfigurasi di kode Go maupun Flutter.

## 3-Level Inheritance
```
PLATFORM SETTINGS (Developer)      ← global, tidak bisa di-override tenant
       ↓
BMT SETTINGS (ADMIN_BMT)           ← berlaku semua cabang
       ↓ (jika is_locked = false)
CABANG SETTINGS (MANAJER_CABANG)   ← override lokal
```

## Resolver
```go
func Resolve(ctx, bmtID, cabangID, kunci string) string
func ResolveInt(...) int
func ResolveBool(...) bool
func ResolveJSON(...) []byte

// WAJIB digunakan — BUKAN konstanta
jam := settings.Resolve(ctx, bmtID, cabangID, "operasional.jam_buka")
```

## Tabel
```sql
CREATE TABLE platform_settings (
    kunci VARCHAR(150) PRIMARY KEY, nilai TEXT,
    tipe VARCHAR(20) DEFAULT 'string',
    deskripsi TEXT, is_rahasia BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMPTZ, updated_by TEXT
);
CREATE TABLE bmt_settings (
    bmt_id UUID REFERENCES bmt(id), kunci VARCHAR(150),
    nilai TEXT, tipe VARCHAR(20) DEFAULT 'string',
    is_locked BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMPTZ, updated_by UUID REFERENCES pengguna(id),
    PRIMARY KEY (bmt_id, kunci)
);
CREATE TABLE cabang_settings (
    cabang_id UUID REFERENCES cabang(id),
    bmt_id UUID REFERENCES bmt(id), kunci VARCHAR(150),
    nilai TEXT, tipe VARCHAR(20) DEFAULT 'string',
    updated_at TIMESTAMPTZ, updated_by UUID REFERENCES pengguna(id),
    PRIMARY KEY (cabang_id, kunci)
);
```

## Kunci Settings Lengkap

### Platform (Developer)
| Kunci | Tipe | Default |
|-------|------|---------|
| `platform.midtrans_server_key` | string* | — |
| `platform.midtrans_client_key` | string* | — |
| `platform.midtrans_env` | string | `"sandbox"` |
| `platform.maintenance_mode` | bool | `"false"` |
| `platform.min_app_version.nasabah` | string | `"2.0.0"` |
| `platform.min_app_version.teller` | string | `"1.5.0"` |
| `platform.rate_limit_rpm` | int | `"300"` |
| `listing.harga_basic_bulanan` | int | `"99000"` |
| `listing.harga_basic_tahunan` | int | `"990000"` |
| `listing.harga_premium_bulanan` | int | `"249000"` |
| `listing.harga_premium_tahunan` | int | `"2490000"` |
| `listing.maks_foto` | int | `"5"` |
| `listing.radius_default_km` | int | `"10"` |
| `backup.retensi_harian_hari` | int | `"7"` |
| `backup.retensi_mingguan_minggu` | int | `"4"` |
| `backup.retensi_bulanan_bulan` | int | `"12"` |

### BMT (ADMIN_BMT)
| Kunci | Tipe | Contoh |
|-------|------|--------|
| `operasional.jam_buka` | time | `"08:00"` |
| `operasional.jam_tutup` | time | `"16:00"` |
| `operasional.hari_kerja` | json | `"[1,2,3,4,5]"` |
| `operasional.zona_waktu` | string | `"Asia/Jakarta"` |
| `autodebet.jam_eksekusi` | time | `"07:00"` |
| `autodebet.tanggal_simpanan_wajib` | int | `"1"` |
| `sesi_teller.toleransi_selisih` | int | `"0"` |
| `approval.FORM_DAFTAR_NASABAH` | json | `"[\"TELLER\"]"` |
| `approval.FORM_UBAH_NASABAH` | json | `"[\"MANAJER_CABANG\"]"` |
| `approval.FORM_BUKA_REKENING` | json | `"[\"TELLER\"]"` |
| `approval.FORM_BLOKIR_REKENING` | json | `"[\"MANAJER_CABANG\"]"` |
| `approval.FORM_TUTUP_REKENING` | json | `"[\"MANAJER_CABANG\"]"` |
| `approval.FORM_BUKA_PEMBIAYAAN` | json | `"[\"KOMITE\"]"` |
| `notifikasi.wa_provider` | string | `"fonnte"` |
| `notifikasi.wa_token` | string* | — |
| `notifikasi.sms_provider` | string | `"zenziva"` |
| `notifikasi.sms_token` | string* | — |
| `notifikasi.reminder_angsuran_hari` | int | `"3"` |
| `notifikasi.reminder_spp_hari` | int | `"3"` |
| `notifikasi.email_smtp_host` | string | — |
| `notifikasi.email_smtp_port` | int | `"587"` |
| `notifikasi.email_smtp_user` | string | — |
| `notifikasi.email_smtp_pass` | string* | — |
| `midtrans.server_key` | string* | — |
| `midtrans.client_key` | string | — |
| `midtrans.env` | string | `"production"` |
| `midtrans.enabled_methods` | json | `"[\"gopay\",\"qris\"]"` |
| `nfc.limit_default_per_transaksi` | int | `"500000"` |
| `nfc.limit_default_harian` | int | `"2000000"` |
| `keamanan.2fa_wajib_staf` | bool | `"true"` |
| `keamanan.lockout_menit` | int | `"15"` |
| `keamanan.maks_gagal_login` | int | `"5"` |
| `pondok.absensi_metode` | json | `"[\"MANUAL\",\"NFC\"]"` |
| `ecommerce.komisi_persen` | float | `"2.5"` |
| `sdm.tanggal_gajian` | int | `"25"` |
| `integrasi.dapodik_aktif` | bool | `"false"` |
| `integrasi.dapodik_npsn` | string | — |
| `integrasi.dapodik_username` | string | — |
| `integrasi.dapodik_password` | string* | — |
| `integrasi.emis_aktif` | bool | `"false"` |
| `integrasi.emis_nsm` | string | — |
| `integrasi.emis_token` | string* | — |
| `whitelabel.nama_app` | string | — |
| `whitelabel.primary_color` | string | — |
| `whitelabel.logo_url` | string | — |

`*` = `is_rahasia = true` — nilai ter-mask di log & response API
