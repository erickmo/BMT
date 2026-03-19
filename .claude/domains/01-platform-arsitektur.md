# Domain: Platform & Arsitektur

## Prinsip Arsitektur Global

- **Monorepo** — satu repo, satu API, semua domain di dalam `api/internal/domain/`
- **Settings over hardcode** — semua nilai yang bisa berbeda antar BMT/pondok **wajib** disimpan sebagai settings di database, bukan konstanta kode
- **Form workflow** — data nasabah & rekening tidak pernah diubah langsung, selalu via form ber-approval
- **Tenant isolation** — semua query operasional wajib di-scope `bmt_id` + `cabang_id`
- **Double-entry accounting** — semua transaksi keuangan wajib catat jurnal via `module-vernon-accounting`
- **Autodebet partial** — jika saldo tidak cukup, debit semampu saldo, sisanya jadi tunggakan
- **Data sebagai konfigurasi** — pecahan uang, jenis rekening, produk, tarif, jadwal autodebet semuanya dari DB

---

## Keputusan Desain yang Sudah Dikonfirmasi

### Platform & Infrastruktur
| Topik | Keputusan |
|-------|-----------|
| Arsitektur | Monorepo, satu API besar (`api/`) |
| Semua frontend | Flutter (platform sesuai per app) |
| Pembuatan BMT & cabang | Developer — via `/dev/*` |
| Pengawasan eksternal | Tidak ada — laporan internal + RAT |

### CBS (Banking)
| Topik | Keputusan |
|-------|-----------|
| Nomor nasabah | 1 per BMT, berlaku lintas cabang |
| Produk & jenis rekening | CRUD management BMT |
| Perubahan data | Wajib via form + approval |
| Approver form | Dikonfigurasi settings BMT per jenis form |
| Pecahan uang Rupiah | **Data di DB** (bukan konstanta) |
| Biaya admin rekening | Per jenis rekening, dikonfigurasi management BMT |
| Tanggal autodebet | Per rekening, diset management BMT di settings |
| Autodebet gagal | Partial debit — sisa jadi tunggakan |
| Revenue platform | Biaya admin per transaksi, dikonfigurasi developer per kontrak |
| Modul akuntansi | `module-vernon-accounting` — internal Go |

### ERP Pondok
| Topik | Keputusan |
|-------|-----------|
| Sistem pondok | Dibangun dari nol dalam proyek ini |
| Relasi nasabah ↔ santri | 1 nasabah = 1 santri |
| Kartu NFC | Tap → PIN 6 digit → debit rekening |
| Absensi | Manual guru + scan NFC + biometrik (sidik jari/wajah) |
| Jadwal | Pelajaran, piket, kegiatan, karyawan, kalender akademik |
| Kurikulum | Mapel, RPP, silabus, jadwal, mapping ke penilaian |
| Penilaian | Harian/UTS/UAS, tahfidz, akhlak, raport digital, peringkat |
| Beasiswa | Ditetapkan admin pondok (sebagian atau seluruh biaya) |
| app/pondok platform | Flutter Web + Mobile |

### E-commerce OPOP
| Topik | Keputusan |
|-------|-----------|
| Model bisnis | B2C (wali → pondok) + B2B (pondok ↔ pondok) |
| Seller | Pondok & BMT/koperasi pondok |
| Buyer | Wali santri (B2C), pondok lain (B2B) |
| Pembayaran | Midtrans + potong saldo rekening BMT + kartu NFC |

---

## Hierarki Tenant

```
PLATFORM (Developer)
└── BMT  ← dibuat & dikonfigurasi developer
    ├── KONTRAK_BMT (tarif, fitur aktif, PIC)
    ├── SETTINGS_BMT (semua konfigurasi — tidak ada hardcode)
    ├── PECAHAN_UANG (data per platform — dikelola developer)
    ├── JENIS_REKENING (CRUD management BMT)
    ├── PRODUK_SIMPANAN & PEMBIAYAAN (CRUD management BMT)
    ├── TARIF_AUTODEBET (jenis, tanggal, retry — settings BMT)
    └── CABANG
        ├── PENGGUNA STAF (Teller, AO, Komite, Finance, Manajer)
        ├── PENGGUNA PONDOK (Admin, Operator, Bendahara Pondok)
        ├── NASABAH (1 nomor per BMT, lintas cabang)
        │   ├── DATA_SANTRI (1:1)
        │   │   ├── ABSENSI
        │   │   ├── PENILAIAN & RAPORT
        │   │   └── KARTU_NFC
        │   └── REKENING
        │       ├── TRANSAKSI_REKENING
        │       └── TUNGGAKAN_AUTODEBET
        ├── AKADEMIK (kurikulum, jadwal, penilaian)
        ├── SDM (guru, karyawan, shift)
        ├── PEMBIAYAAN (termasuk beasiswa)
        ├── SESI_TELLER
        ├── JURNAL_AKUNTANSI
        └── TOKO_OPOP (produk pondok untuk marketplace)
```

---

## Struktur Direktori Proyek

```
bmt-saas/
├── api/                              # Backend Go — satu API semua domain
├── apps-nasabah/                     # Flutter Android + iOS
├── apps-management/                  # Flutter Web + Desktop + Mobile
├── apps-developer/                   # Flutter Web + Desktop + Mobile
├── apps-teller/                      # Flutter Desktop
├── apps-merchant/                    # Flutter Android + iOS
├── apps-ceksaldo/                    # Flutter Android / Kiosk
├── apps-pondok/                      # Flutter Web + Mobile
├── module-vernon-accounting/         # Modul akuntansi double-entry internal
└── docs/
```

---

## Stack Teknologi

```
Backend API              : Go 1.23 (net/http + Chi router)
Akuntansi                : module-vernon-accounting (Go internal)
Database                 : PostgreSQL 16
Cache                    : Redis 7
Queue/Worker             : asynq (Redis-backed)
Auth                     : JWT — access token 15m + refresh token 7d
Storage                  : MinIO (self-hosted S3-compatible)
PDF Generator            : chromedp (Go → headless Chrome)
DB Migration             : golang-migrate
Query                    : sqlc (type-safe SQL → Go, BUKAN ORM)
Payment Gateway          : Midtrans (Snap + Core API)
Email                    : SMTP — dikonfigurasi per BMT di settings
Biometrik                : Fingerprint SDK (integrasi perangkat via REST)
Search (OPOP)            : PostgreSQL full-text search (fase 1), Elasticsearch (fase 2)

── Flutter Apps ──────────────────────────────────────────────────────────────
apps-nasabah             : Flutter → Android, iOS
apps-management          : Flutter → Web, Desktop (Win/macOS/Linux), Mobile
apps-developer           : Flutter → Web, Desktop (Win/macOS/Linux), Mobile
apps-teller              : Flutter → Desktop (Win/macOS/Linux)
apps-merchant            : Flutter → Android, iOS
apps-ceksaldo            : Flutter → Android / Kiosk
apps-pondok              : Flutter → Web, Mobile (Android/iOS)
──────────────────────────────────────────────────────────────────────────────

Testing                  : Go testing + testify + testcontainers-go
                           Flutter widget test + integration test
Container                : Docker + Docker Compose
```

---

## Struktur Direktori API

```
api/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── platform/              # BMT, Cabang, Kontrak, PlatformSettings
│   │   ├── settings/              # Settings engine (platform/bmt/cabang resolver)
│   │   ├── nasabah/               # Nasabah, KartuNFC
│   │   ├── form/                  # Form workflow + approval engine
│   │   ├── rekening/              # Rekening, JenisRekening
│   │   ├── transaksi/             # Transaksi rekening
│   │   ├── autodebet/             # Jadwal, eksekusi partial, tunggakan
│   │   ├── sesi_teller/           # Sesi kas + redenominasi (pecahan dari DB)
│   │   ├── pembiayaan/            # Pembiayaan, Angsuran, Beasiswa
│   │   ├── finance/               # Jurnal manual, vendor, operasional
│   │   ├── nfc/                   # Kartu NFC, terminal, transaksi
│   │   ├── merchant/              # Merchant pondok (NFC)
│   │   │
│   │   ├── pondok/
│   │   │   ├── administrasi/      # Santri, pengajar, karyawan, alumni
│   │   │   ├── akademik/          # Kurikulum, mapel, RPP, silabus
│   │   │   ├── jadwal/            # Jadwal pelajaran, kegiatan, piket, shift, kalender
│   │   │   ├── absensi/           # Absensi santri & karyawan (manual/NFC/biometrik)
│   │   │   ├── penilaian/         # Nilai, tahfidz, akhlak, raport, peringkat
│   │   │   ├── keuangan/          # Tagihan SPP, beasiswa, pembiayaan pondok
│   │   │   ├── sdm/               # Data karyawan, kontrak, shift
│   │   │   ├── perpustakaan/
│   │   │   ├── konsultasi/
│   │   │   ├── surat_izin/
│   │   │   ├── health_record/
│   │   │   └── alumni/
│   │   │
│   │   ├── ecommerce/
│   │   │   ├── toko/              # Toko per pondok, profil seller
│   │   │   ├── produk/            # Produk, kategori, stok, varian
│   │   │   ├── pesanan/           # Keranjang, order, checkout, status
│   │   │   ├── pengiriman/        # Alamat, kurir, tracking
│   │   │   ├── pembayaran/        # Midtrans, rekening BMT, NFC
│   │   │   ├── ulasan/            # Rating & review produk
│   │   │   └── opop/              # OPOP marketplace lintas pondok (B2B)
│   │   │
│   │   ├── notifikasi/
│   │   ├── keamanan/
│   │   ├── analytics/
│   │   ├── sosial_keuangan/       # Donasi, wakaf, infaq
│   │   ├── sdm/                   # Kontrak, slip gaji, payroll
│   │   ├── inventaris/
│   │   ├── integrasi/             # DAPODIK, EMIS, PPDB
│   │   ├── monetisasi/
│   │   └── payment/               # Midtrans, UsageLog
│   │
│   ├── handler/
│   │   ├── developer/             # /dev/*
│   │   ├── platform/              # /platform/*
│   │   ├── teller/                # /teller/*
│   │   ├── finance/               # /finance/*
│   │   ├── nfc/                   # /nfc/*
│   │   ├── merchant/              # /merchant/*
│   │   ├── pondok/                # /pondok/*
│   │   ├── ecommerce/             # /shop/*
│   │   └── nasabah/               # /nasabah/*
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── tenant.go
│   │   ├── developer.go
│   │   ├── idempotency.go
│   │   └── ratelimit.go
│   ├── repository/postgres/
│   │   ├── query/                 # .sql per domain
│   │   ├── sqlc.yaml
│   │   └── *.go
│   ├── service/
│   ├── worker/
│   └── config/
├── migrations/
└── pkg/
    ├── money/                     # type Money int64
    ├── syariah/                   # Kalkulasi akad
    ├── midtrans/
    ├── pdfgen/
    ├── audit/
    └── response/
```
