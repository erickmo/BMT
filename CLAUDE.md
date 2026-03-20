# CLAUDE.md — Platform Pesantren Terpadu (SaaS)

> **Terakhir diperbarui:** 20 Maret 2026
> **Versi:** 1.5

Platform digital terpadu untuk pondok pesantren:
**CBS (Core Banking) + ERP Pondok + E-commerce OPOP + Listing Stakeholder**
— satu monorepo, satu API Go, tujuh Flutter apps.

## Sprint Aktif

**Sprint 1 — 20 Mar s/d 3 Apr 2026** · [Rencana lengkap → docs/sprint-plan.md](docs/sprint-plan.md)

| # | Item | Status |
|---|------|--------|
| 1 | `RekeningService.Tarik()` + `Transfer()` | ✅ |
| 2 | `NasabahService` — GetByID, Search, ListRekening, GetMutasi | ✅ |
| 3 | `SesiTellerService` — buka/tutup sesi, validasi selisih | ✅ |
| 4 | Wire handler `/teller/*` + `/nasabah/*` → services | ✅ |
| 5 | `middleware/feature_gate.go` — `RequireFeature(kode)` | ✅ |
| 6 | `middleware/audit_log.go` — catat semua mutasi | ✅ |
| 7 | `PlatformFeatureChecker` + `AuditRepository` | ✅ |

---

## Dokumentasi → `docs/`

| File | Domain |
|------|--------|
| [docs/01-arsitektur.md](docs/01-arsitektur.md) | Hierarki tenant, prinsip global, struktur direktori |
| [docs/02-stack.md](docs/02-stack.md) | Tech stack, 7 Flutter apps, module-vernon-accounting |
| [docs/03-saas.md](docs/03-saas.md) | **SaaS: tier paket, add-on fitur, portal developer, listing stakeholder** |
| [docs/04-settings.md](docs/04-settings.md) | Settings engine 3-level, semua kunci konfigurasi |
| [docs/05-cbs.md](docs/05-cbs.md) | CBS: nasabah, rekening, transaksi, autodebet, teller, pembiayaan |
| [docs/06-pondok-akademik.md](docs/06-pondok-akademik.md) | Administrasi, akademik, kurikulum, jadwal, absensi, penilaian |
| [docs/07-pondok-ops.md](docs/07-pondok-ops.md) | Perpustakaan, UKS, asrama, inventaris, PPDB, konsultasi, izin |
| [docs/08-pondok-pengembangan.md](docs/08-pondok-pengembangan.md) | Portfolio, hafalan, ekstra, alumni, event, SDM, payroll |
| [docs/09-ecommerce.md](docs/09-ecommerce.md) | OPOP B2C+B2B, toko, produk, pesanan |
| [docs/10-keuangan.md](docs/10-keuangan.md) | Akuntansi, donasi, wakaf, zakat |
| [docs/11-keamanan.md](docs/11-keamanan.md) | 2FA, session, audit log, anti-fraud, rate limiting, offline mode |
| [docs/12-notifikasi.md](docs/12-notifikasi.md) | FCM, WhatsApp, SMS, Email, Bulletin Board, chat staf |
| [docs/13-integrasi.md](docs/13-integrasi.md) | Midtrans, DAPODIK, EMIS, white-label |
| [docs/14-api.md](docs/14-api.md) | Semua endpoint per domain |
| [docs/15-workers.md](docs/15-workers.md) | Background workers (asynq) |
| [docs/16-konvensi.md](docs/16-konvensi.md) | Konvensi Go & Flutter, error sentinel, testing |
| [docs/17-form-workflow.md](docs/17-form-workflow.md) | Form pengajuan + approval engine |
| [docs/18-syariah-glosarium.md](docs/18-syariah-glosarium.md) | Checklist syariah, role lengkap, glosarium |

---

## Stack Cepat
```
Backend  : Go 1.23 + Chi router    [api/]
Database : PostgreSQL 17 + Redis 7
Query    : sqlc (bukan ORM)
Queue    : asynq
Storage  : MinIO
Payment  : Midtrans
Akuntansi: module-vernon-accounting/ (internal Go)
HRM      : module-vernon-hrm/          (internal Go)
Apps     : 7 Flutter (→ docs/02-stack.md)
```

## 7 Prinsip Wajib
1. **TIDAK ADA HARDCODE** — semua konfigurasi dari `settings` DB
2. **TIDAK ADA UPDATE LANGSUNG** — nasabah/rekening wajib via form + approval
3. **SETIAP TRANSAKSI** = DB transaction + jurnal (module-vernon-accounting) + usage_log
4. **SETIAP QUERY** = di-scope `bmt_id` + `cabang_id`
5. **UANG** = `type Money int64` — **TIDAK PERNAH** float
6. **AUTODEBET GAGAL** = partial debit + INSERT tunggakan
7. **JURNAL** = Σ debit = Σ kredit, divalidasi sebelum persist

## Hierarki Tenant
```
PLATFORM (Developer)
└── BMT  ← dibuat developer, termasuk paket tier & fitur aktif
    └── CABANG
        ├── Nasabah / Rekening / Transaksi
        ├── Data Pondok (Santri, Akademik, SDM)
        └── Toko OPOP
```

## Struktur Direktori
```
bmt-saas/
├── api/                          # Go — satu API semua domain
├── app/
│   ├── nasabah/                  # Flutter Android + iOS
│   ├── management/               # Flutter Web + Desktop + Mobile
│   ├── developer/                # Flutter Web + Desktop + Mobile
│   ├── teller/                   # Flutter Desktop
│   ├── merchant/                 # Flutter Android + iOS
│   ├── ceksaldo/                 # Flutter Android / Kiosk
│   └── pondok/                   # Flutter Web + Mobile
├── module-vernon-accounting/
├── module-vernon-hrm/
└── docs/
```
