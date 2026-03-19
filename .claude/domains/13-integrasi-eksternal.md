# Domain: Integrasi Eksternal

## Prinsip

Semua integrasi eksternal bersifat **opsional per BMT** (dikonfigurasi di bmt_settings).
Jika tidak dikonfigurasi, fitur menggunakan data internal saja.

---

## Log Sinkronisasi

```sql
CREATE TABLE integrasi_log (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id      UUID NOT NULL REFERENCES bmt(id),
    provider    VARCHAR(30) NOT NULL,  -- DAPODIK | EMIS | PPDB
    arah        VARCHAR(10) NOT NULL,  -- PULL | PUSH
    status      VARCHAR(20) NOT NULL,  -- SUKSES | GAGAL | PARTIAL
    jumlah_record INT,
    error_detail TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## EMIS Kemenag

- **Tujuan:** Sinkronisasi data santri & lembaga dengan database Kemenag
- **Arah:** Pull dari EMIS (verifikasi NSM/NPSN) + Push (update kehadiran/data santri)
- **Konfigurasi:**
  ```
  integrasi.emis_aktif    → "true"
  integrasi.emis_nsm      → "21237..."
  integrasi.emis_token    → "..."   (is_rahasia = true)
  integrasi.sinkron_jadwal → "MINGGUAN"
  ```
- **Worker:** `SinkronEMIS` — Sabtu 02:00 WIB
- **Endpoint:** `POST /pondok/sinkron/emis`

---

## DAPODIK

- **Tujuan:** Sinkronisasi data siswa untuk pondok di bawah Kemendikbud
- **Arah:** Pull dari DAPODIK (import peserta didik) + Push (nilai/absensi)
- **Konfigurasi:**
  ```
  integrasi.dapodik_aktif    → "true"
  integrasi.dapodik_npsn     → "20572345"
  integrasi.dapodik_username → "..."
  integrasi.dapodik_password → "..."  (is_rahasia = true)
  ```
- **Worker:** `SinkronDAPODIK` — tgl 1, 03:00 WIB
- **Endpoint:** `POST /pondok/sinkron/dapodik`

---

## PPDB Online (Penerimaan Santri Baru)

```sql
CREATE TABLE ppdb_pendaftar (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    cabang_id           UUID NOT NULL REFERENCES cabang(id),
    tahun_ajaran        VARCHAR(10) NOT NULL,
    nama_lengkap        VARCHAR(255) NOT NULL,
    nik                 VARCHAR(16),
    tanggal_lahir       DATE,
    nama_wali           VARCHAR(255),
    telepon_wali        VARCHAR(20),
    email_wali          VARCHAR(255),
    pilihan_tingkat     VARCHAR(20),
    status              VARCHAR(20) NOT NULL DEFAULT 'DAFTAR',
    -- DAFTAR | SELEKSI | DITERIMA | DITOLAK | MUNDUR
    nomor_pendaftaran   VARCHAR(30) UNIQUE NOT NULL,
    dokumen             JSONB NOT NULL DEFAULT '{}',
    -- {"kk": "url", "akta": "url", "ijazah": "url"}
    catatan             TEXT,
    nasabah_id          UUID REFERENCES nasabah(id),
    -- Jika diterima → bisa langsung dibuat akun nasabah
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

PPDB dapat diakses publik (tanpa login) via URL pendaftaran per pondok.
Pembayaran biaya pendaftaran via Midtrans.
