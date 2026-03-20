# Platform Pesantren Terpadu — BMT SaaS

Platform ekosistem digital terpadu untuk pondok pesantren. Satu monorepo yang menggabungkan **Core Banking System (CBS)** syariah, **ERP Pondok**, **E-commerce OPOP**, dan **7 aplikasi Flutter** dalam satu API Go.

---

## Daftar Isi

- [Gambaran Sistem](#gambaran-sistem)
- [Aplikasi](#aplikasi)
- [Stack Teknologi](#stack-teknologi)
- [Struktur Direktori](#struktur-direktori)
- [Prinsip Utama](#prinsip-utama)
- [Quick Start — Development](#quick-start--development)
- [Deployment Docker](#deployment-docker)
- [Migrasi Database](#migrasi-database)
- [Testing](#testing)
- [Role & Akses](#role--akses)

---

## Gambaran Sistem

```
┌─────────────────────────────────────────────────────────┐
│                    Satu API Go (REST)                   │
│  /auth  /platform  /teller  /nasabah  /pondok  /shop   │
│  /nfc   /finance   /opop    /merchant /dev             │
└─────────────────────────┬───────────────────────────────┘
                          │
         ┌────────────────┼────────────────┐
         ▼                ▼                ▼
   PostgreSQL 16       Redis 7          MinIO
  (external/self)    (Queue+Cache)   (File Storage)
```

**Multi-tenant:** Setiap BMT adalah tenant independen. Semua data di-scope dengan `bmt_id + cabang_id` — tidak ada data yang bocor antar BMT.

---

## Aplikasi

| App | Platform | Pengguna | Deskripsi |
|-----|----------|----------|-----------|
| `apps-nasabah` | Android + iOS | Nasabah, Wali Santri, Alumni | E-banking, top-up NFC, belanja OPOP, raport digital |
| `apps-management` | Web + Desktop + Mobile | Staf & Management BMT | Dashboard, pembiayaan, laporan, settings |
| `apps-developer` | Web + Desktop + Mobile | Developer Platform | CRUD BMT, kontrak, pecahan uang, platform settings |
| `apps-teller` | Desktop | Teller Cabang | Sesi kas, transaksi tunai, cetak slip |
| `apps-merchant` | Android + iOS | Kasir & Owner Toko Pondok | Transaksi NFC, laporan penjualan |
| `apps-ceksaldo` | Android / Kiosk | Santri | Tap NFC → lihat saldo + 5 transaksi terakhir |
| `apps-pondok` | Web + Mobile | Admin Pondok | Akademik, absensi, raport, SPP, perpustakaan, PPDB |

---

## Stack Teknologi

**Backend**
- Go 1.24 · Chi Router · zerolog
- PostgreSQL 16 · pgx/v5 (tanpa ORM) · golang-migrate
- Redis 7 · asynq (background workers)
- MinIO (object storage)
- Midtrans (payment gateway)
- JWT (access 15m + refresh 7d)

**Frontend**
- Flutter (7 apps — Android, iOS, Web, Desktop)

**Akuntansi**
- `module-vernon-accounting` — double-entry accounting engine internal

---

## Struktur Direktori

```
bmt-saas/
├── api/                          # Backend Go — satu API semua domain
│   ├── cmd/server/main.go        # Entry point
│   ├── internal/
│   │   ├── config/               # Config loader (env vars + .env)
│   │   ├── domain/               # Business logic & domain errors
│   │   │   ├── rekening/
│   │   │   ├── autodebet/
│   │   │   ├── nfc/
│   │   │   ├── form/
│   │   │   ├── sesi_teller/
│   │   │   ├── pembiayaan/
│   │   │   ├── ecommerce/
│   │   │   └── pondok/
│   │   ├── handler/              # HTTP handlers (satu folder per domain)
│   │   ├── middleware/           # Auth JWT, tenant, idempotency, developer
│   │   ├── repository/postgres/  # Query PostgreSQL (parameterized, bukan ORM)
│   │   ├── service/              # Business logic layer
│   │   └── worker/               # Background workers (asynq)
│   ├── pkg/
│   │   ├── jwt/                  # JWT manager
│   │   ├── money/                # Type Money (int64) — tidak pernah float
│   │   ├── response/             # HTTP response helper
│   │   └── settings/             # Settings resolver 3-level
│   ├── migrations/               # SQL migrations (golang-migrate)
│   ├── Dockerfile
│   ├── Makefile
│   └── docker-compose.yml        # Development (dengan PostgreSQL lokal)
│
├── apps-nasabah/                 # Flutter Android + iOS
├── apps-management/              # Flutter Web + Desktop + Mobile
├── apps-developer/               # Flutter Web + Desktop + Mobile
├── apps-teller/                  # Flutter Desktop
├── apps-merchant/                # Flutter Android + iOS
├── apps-ceksaldo/                # Flutter Android / Kiosk
├── apps-pondok/                  # Flutter Web + Mobile
│
├── module-vernon-accounting/     # Modul akuntansi double-entry
│
├── docker-compose.yml            # Production (PostgreSQL external)
├── .env.example                  # Template environment variables
└── docs/
    ├── audit/                    # Laporan audit code
    └── requirements/             # PRD dan requirement fitur
```

---

## Prinsip Utama

Tujuh aturan yang **wajib dipatuhi** di seluruh codebase:

```
1. TIDAK ADA HARDCODE          → semua konfigurasi dari settings DB
2. TIDAK ADA UPDATE LANGSUNG   → data nasabah/rekening via form + approval
3. SETIAP TRANSAKSI KEUANGAN   → DB transaction + jurnal + usage_log
4. SETIAP QUERY                → di-scope bmt_id + cabang_id (tenant isolation)
5. UANG = int64 (Money)        → tidak pernah float
6. AUTODEBET GAGAL             → partial debit + INSERT tunggakan (bukan skip)
7. JURNAL                      → Σ debit = Σ kredit, divalidasi sebelum persist
```

---

## Quick Start — Development

### Prasyarat

- Go 1.24+
- Flutter 3.x
- Docker & Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate)

### 1. Clone dan masuk ke direktori

```bash
git clone https://github.com/erickmo/BMT.git
cd BMT
```

### 2. Jalankan infrastruktur development (PostgreSQL + Redis + MinIO)

```bash
cd api
docker compose up -d
```

### 3. Setup environment

```bash
cp .env.example api/.env
# Edit api/.env — isi DATABASE_URL, JWT secrets, dll.
```

### 4. Jalankan migrasi

```bash
cd api
make migrate-up
```

### 5. Jalankan API

```bash
cd api
make dev
# API berjalan di http://localhost:8080
```

### 6. Jalankan Flutter app (contoh: apps-management)

```bash
cd apps-management
flutter pub get
flutter run -d chrome --dart-define=API_BASE_URL=http://localhost:8080
```

---

## Deployment Docker

PostgreSQL dikelola di luar Docker (external). Docker hanya menjalankan API, Redis, MinIO, dan 3 web app.

### 1. Siapkan environment

```bash
cp .env.example .env
nano .env
```

Variabel wajib diisi:

```env
# Koneksi ke PostgreSQL external
DATABASE_URL=postgres://user:password@host:5432/bmt_saas_db?sslmode=disable

# JWT secrets (random string >= 32 karakter)
JWT_ACCESS_SECRET=...
JWT_REFRESH_SECRET=...

# URL API yang dapat diakses dari browser (bukan nama Docker service)
API_BASE_URL=http://your-server-ip-or-domain:8080

# Developer token untuk endpoint /dev/*
DEVELOPER_TOKEN=...
```

### 2. Build dan jalankan

```bash
docker compose up -d --build
```

### 3. Cek status

```bash
docker compose ps
docker compose logs api --follow
```

### Port Default

| Service | Port | Keterangan |
|---------|------|------------|
| API | `8080` | Go backend |
| Web Management | `3001` | Staff & management BMT |
| Web Developer | `3002` | Developer platform |
| Web Pondok | `3003` | Admin pondok |
| MinIO API | `9000` | Object storage |
| MinIO Console | `9001` | Web UI MinIO |

> Port dapat diubah via env vars: `API_PORT`, `WEB_MANAGEMENT_PORT`, dll.

### Update deployment

```bash
git pull
docker compose up -d --build
```

### Catatan: API_BASE_URL di Flutter Web

`API_BASE_URL` diinjeksikan saat container start — **tidak perlu rebuild image** saat ganti URL. Cukup update `.env` dan restart container:

```bash
docker compose restart web-management web-developer web-pondok
```

---

## Migrasi Database

```bash
cd api

# Jalankan semua migrasi
make migrate-up

# Rollback 1 step
make migrate-down

# Buat migration baru
make migrate-create name=nama_migration
```

File migrasi ada di `api/migrations/`.

---

## Testing

```bash
cd api

# Semua unit test
make test

# Dengan race detector (wajib sebelum merge)
make test-race

# Integration test (memerlukan Docker untuk testcontainers)
make test-integration
```

### Test Wajib

Test berikut harus selalu pass:

```
TestSettings_TidakAdaHardcode_SelaluDariDB
TestPecahanUang_DariDB_BukanKonstanta
TestAutodebet_SaldoKurang_PartialDebitDanTunggakan
TestJurnal_SemuaTransaksi_DoubleEntryBalance
TestCrossTenant_QueryTanpaBMTID_Dilarang
```

---

## Role & Akses

| Role | App | Kapabilitas |
|------|-----|-------------|
| `DEVELOPER` | apps-developer | CRUD BMT, kontrak, pecahan uang, platform settings |
| `ADMIN_BMT` | apps-management | Settings, jenis rekening, pengguna |
| `MANAJER_BMT` | apps-management | Laporan konsolidasi, approval besar |
| `MANAJER_CABANG` | apps-management | Approval form, laporan cabang |
| `TELLER` | apps-teller | Transaksi tunai, sesi kas |
| `FINANCE` | apps-management | Jurnal manual, biaya operasional |
| `ACCOUNT_OFFICER` | apps-management | Pengajuan & monitoring pembiayaan |
| `NASABAH` | apps-nasabah | E-banking, NFC, belanja OPOP |
| `KASIR_MERCHANT` | apps-merchant | Transaksi NFC kasir |
| `OWNER_MERCHANT` | apps-merchant | Laporan penjualan |
| `ADMIN_PONDOK` | apps-pondok | Semua fitur pondok |
| `OPERATOR_PONDOK` | apps-pondok | Input santri, absensi, nilai |
| `BENDAHARA_PONDOK` | apps-pondok | Tagihan, beasiswa, laporan keuangan |
| `ACCOUNT_OFFICER` | apps-management | Pengajuan & monitoring pembiayaan |
| `AUDITOR_BMT` | apps-management | Read-only semua laporan keuangan |
| `WALI_SANTRI` | apps-nasabah | Monitoring akademik anak, bayar SPP |
| `SANTRI` | apps-ceksaldo + apps-nasabah | Cek saldo NFC, raport digital |

**Aturan khusus:**
- Teller: semua tombol transaksi **disabled** tanpa sesi kas aktif
- Kiosk (`apps-ceksaldo`): tidak ada login, tidak ada PIN — hanya IP whitelist terminal
- Transaksi NFC: wajib header `X-Idempotency-Key`
- Developer: akses via header `Developer-Token` (bukan JWT)

---

## Background Workers

Workers berjalan otomatis via asynq (Redis-backed). Jadwal dalam WIB.

| Worker | Jadwal | Fungsi |
|--------|--------|--------|
| `AutodebetHarian` | 07:00 daily | Eksekusi jadwal autodebet |
| `AutodebetBulanan` | Tgl 1, 07:30 | Biaya admin rekening bulanan |
| `GenerateTagihanSPP` | Tgl 25, 08:30 | Generate tagihan SPP periode berikutnya |
| `DistribusiBagiHasil` | Akhir bulan | Posting bagi hasil deposito |
| `GenerateSlipGaji` | Tgl 25, 09:00 | Generate slip gaji pegawai |
| `EksekusiPayroll` | Tgl 1, 08:00 | Transfer gaji ke rekening pegawai |
| `KirimNotifikasi` | Setiap 1 menit | Proses antrian FCM/WA/SMS/Email |
| `BackupDatabase` | 02:00 daily | Dump PostgreSQL → MinIO |
| `HitungZakat` | 31 Des, 23:00 | Hitung zakat mal akhir tahun per nasabah |

---

## Checklist Syariah

Prinsip syariah yang diterapkan di seluruh sistem:

- **Ta'zir (denda):** 100% masuk akun dana sosial — bukan pendapatan BMT
- **Bagi hasil:** dari realisasi pendapatan, bukan % nominal pokok
- **Margin/nisbah/ujrah:** transparan & disepakati sebelum akad
- **Autodebet partial:** menghasilkan jurnal syariah yang benar meski saldo kurang
- **Wakaf produktif:** hasil usaha ke mauquf alaih, bukan ke BMT
- **Dana sosial:** tidak bercampur dengan dana operasional BMT

---

## Status Pengembangan

| Sprint | Scope | Status |
|--------|-------|--------|
| Sprint 1 | CBS Core (Teller, Nasabah, Rekening, Autodebet) | ✅ |
| Sprint 2 | Keamanan + Middleware (Auth JWT, OTP, Session, Feature Gate, Audit Log) | ✅ |
| Sprint 3 | Workers CBS (Autodebet bulanan, Kolektibilitas OJK, Distribusi Bagi Hasil, Reminder Angsuran) | ✅ |
| Sprint 4 | Notifikasi + Midtrans (FCM, WhatsApp, SMS, Email, Webhook Midtrans) | ✅ |
| Sprint 5 | Pembiayaan + Akuntansi (State machine pembiayaan, Laporan Neraca/SHU/Arus Kas, Zakat) | ✅ |
| Sprint 6 | Pondok (Santri, Akademik, SPP, Raport) | upcoming |
| Sprint 7 | OPOP E-commerce + NFC | upcoming |
| Sprint 8 | SaaS Portal + Fraud + Integrasi DAPODIK/EMIS | upcoming |

---

## Lisensi

Proprietary — hak cipta milik pengembang. Tidak untuk didistribusikan.
