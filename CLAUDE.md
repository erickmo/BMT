# CLAUDE.md — Platform Pesantren Terpadu (SaaS)

## Gambaran Proyek

Platform ekosistem digital terpadu untuk pondok pesantren dalam **satu monorepo**, **satu API Go**, dan **tujuh aplikasi Flutter**. Sistem wajib mematuhi prinsip syariah Islam di seluruh domain keuangan.

---

## Prinsip Utama (Selalu Ikuti)

```
1. TIDAK ADA HARDCODE — semua konfigurasi dari settings DB
2. TIDAK ADA UPDATE LANGSUNG — data nasabah/rekening via form + approval
3. SETIAP TRANSAKSI KEUANGAN = DB transaction + jurnal + usage_log
4. SETIAP QUERY = di-scope bmt_id + cabang_id (tenant isolation)
5. UANG = int64 (Money), TIDAK PERNAH float
6. AUTODEBET GAGAL = partial debit + INSERT tunggakan
```

---

## Stack Teknologi

| Layer | Teknologi |
|-------|-----------|
| Backend API | Go 1.23 (net/http + Chi router) |
| Database | PostgreSQL 16 + sqlc (BUKAN ORM) |
| Cache / Queue | Redis 7 + asynq |
| Auth | JWT — access 15m + refresh 7d |
| Storage | MinIO (self-hosted S3) |
| Payment | Midtrans (Snap + Core API) |
| PDF | chromedp |
| Migration | golang-migrate |
| Frontend | Flutter (7 apps — lihat domain 20) |
| Akuntansi | module-vernon-accounting (Go internal) |

---

## Struktur Direktori

```
bmt-saas/
├── api/                          # Backend Go — satu API semua domain
├── apps-nasabah/                 # Flutter Android + iOS
├── apps-management/              # Flutter Web + Desktop + Mobile
├── apps-developer/               # Flutter Web + Desktop + Mobile
├── apps-teller/                  # Flutter Desktop
├── apps-merchant/                # Flutter Android + iOS
├── apps-ceksaldo/                # Flutter Android / Kiosk
├── apps-pondok/                  # Flutter Web + Mobile
├── module-vernon-accounting/     # Modul akuntansi double-entry
└── docs/
```

---

## Domain Files — Baca Sesuai Konteks Pekerjaan

Detail lengkap tersedia di `.claude/domains/`. Baca domain yang relevan sebelum bekerja.

| File | Domain | Kapan Dibaca |
|------|--------|--------------|
| `01-platform-arsitektur.md` | Prinsip, tenant hierarchy, struktur API | Setup, onboarding, arsitektur baru |
| `02-settings-engine.md` | Settings tables, resolver, pecahan uang | Saat membuat fitur yang butuh konfigurasi |
| `03-cbs-banking.md` | Rekening, autodebet, pembiayaan, akuntansi | Domain CBS / banking |
| `04-erp-pondok.md` | Santri, kurikulum, jadwal, absensi, raport | Domain ERP pondok |
| `05-ecommerce-opop.md` | Toko, produk, pesanan, OPOP | Domain e-commerce |
| `06-notifikasi.md` | FCM, WA, SMS, Email, Bulletin Board | Domain notifikasi |
| `07-keamanan-audit.md` | 2FA, session, audit log, anti-fraud | Domain keamanan |
| `08-analytics-laporan.md` | Dashboard RT, laporan custom | Domain analytics |
| `09-donasi-wakaf.md` | Donasi, wakaf produktif, infaq | Domain sosial keuangan |
| `10-sdm-payroll.md` | Kontrak, slip gaji, payroll | Domain SDM |
| `11-sosial-pondok.md` | Perpustakaan, konsultasi, izin, UKS, alumni | Domain sosial pondok |
| `12-inventaris-aset.md` | Aset tetap, peminjaman ruang | Domain inventaris |
| `13-integrasi-eksternal.md` | DAPODIK, EMIS, PPDB | Domain integrasi |
| `14-monetisasi.md` | Komisi OPOP, iklan, white-label | Domain monetisasi |
| `15-api-endpoints.md` | Semua endpoint API | Saat membuat/mengubah handler |
| `16-workers.md` | Background workers (jadwal & deskripsi) | Saat membuat/mengubah worker |
| `17-roles-permissions.md` | Semua role & aturan akses | Saat mengatur middleware auth |
| `18-testing.md` | Strategi testing, test wajib | Saat membuat test |
| `19-konvensi-koding.md` | Pola Go wajib, error sentinel, glosarium | **Selalu dibaca** |
| `20-flutter-apps.md` | Deskripsi 7 Flutter apps | Saat bekerja di Flutter app |
