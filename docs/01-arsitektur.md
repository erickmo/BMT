# 01 — Arsitektur

> **Terakhir diperbarui:** 20 Maret 2026

## Domain Utama
1. **CBS** — Core Banking System BMT syariah, multi-tenant, multi-cabang
2. **ERP Pondok** — Administrasi, akademik, SDM, operasional pesantren
3. **E-commerce OPOP** — Marketplace produk pondok (B2C wali + B2B antar pondok)
4. **Listing Stakeholder** — Direktori layanan sekitar pondok (guru les, kursus, antar jemput, dll.)
5. **SaaS Platform** — Tier paket + add-on, portal developer, billing

## Prinsip Global
- **Monorepo** — satu repo, satu API, semua domain di `api/internal/domain/`
- **Settings over hardcode** — semua nilai konfigurasi dari DB (→ docs/04-settings.md)
- **Form workflow** — nasabah/rekening tidak pernah diubah langsung (→ docs/17-form-workflow.md)
- **Tenant isolation** — semua query wajib di-scope `bmt_id` + `cabang_id`
- **Double-entry** — semua transaksi keuangan catat jurnal via module-vernon-accounting
- **HRM terpusat** — semua proses kepegawaian via module-vernon-hrm
- **Autodebet partial** — gagal → debit semampu saldo, sisa jadi tunggakan
- **Feature gating** — setiap fitur dicek via `featureGate.IsEnabled(ctx, bmtID, "kode_fitur")`

## Hierarki Tenant
```
PLATFORM (Developer)
└── BMT  ← dibuat developer via portal dev
    ├── PAKET_TIER (FREE|BASIC|PRO|ENTERPRISE)
    ├── ADD_ON_FITUR (fitur à la carte di atas tier)
    ├── KONTRAK (harga custom, PIC, dokumen)
    ├── SETTINGS (konfigurasi operasional)
    ├── JENIS_REKENING (CRUD management BMT)
    ├── PRODUK_SIMPANAN & PEMBIAYAAN (CRUD management BMT)
    └── CABANG
        ├── PENGGUNA STAF
        ├── PENGGUNA PONDOK
        ├── NASABAH → REKENING → TRANSAKSI
        ├── DATA SANTRI → AKADEMIK
        └── TOKO OPOP
```

## Keputusan Desain Dikonfirmasi

| Topik | Keputusan |
|-------|-----------|
| Arsitektur | Monorepo, satu API |
| Model SaaS | Paket tier + add-on per fitur |
| Pembuatan BMT | Developer via portal dev (`app/developer`) |
| Listing stakeholder | Self-register via form publik → developer approve |
| Listing akses | Wali santri bisa lihat & kontak di `app/nasabah` |
| Listing biaya | Langganan bulanan/tahunan + fitur premium |
| Platform melayani | BMT pondok pesantren |
| Nomor nasabah | 1 per BMT, berlaku lintas cabang |
| Produk & jenis rekening | CRUD management BMT |
| Perubahan data | Wajib via form + approval |
| Approver | Dikonfigurasi settings BMT per form |
| Pecahan uang | Data di DB (bukan konstanta) |
| Autodebet gagal | Partial debit → tunggakan |
| Relasi nasabah ↔ santri | 1 nasabah = 1 santri |
| Kartu NFC | Tap → PIN 6 digit → debit rekening |
| Absensi | Manual + NFC + biometrik |
| Beasiswa | Ditetapkan admin pondok |
| Sistem pondok | Dibangun dari nol dalam proyek ini |
| E-commerce | B2C wali + B2B antar pondok |
| Monetisasi OPOP | Komisi % + iklan premium |
| Integrasi | DAPODIK + EMIS Kemenag |
| Offline mode | Teller (transaksi) + absensi |
| White-label | Custom branding app per pondok |

## Struktur Direktori API
```
api/internal/domain/
├── platform/          # BMT, Cabang, Kontrak
├── saas/              # Tier, AddOn, FeatureGate, Billing
├── listing/           # Listing stakeholder, kategori, review
├── settings/          # Settings engine 3-level
├── nasabah/
├── form/              # Form workflow + approval engine
├── rekening/
├── transaksi/
├── autodebet/
├── sesi_teller/
├── pembiayaan/
├── finance/
├── nfc/
├── merchant/
├── notifikasi/
├── keamanan/
├── analytics/
├── donasi_wakaf/
├── sdm/
├── inventaris/
├── ppdb/
├── integrasi/
├── monetisasi/
├── payment/
├── ecommerce/
│   ├── toko/ produk/ pesanan/ pembayaran/ ulasan/ opop/
└── pondok/
    ├── administrasi/ akademik/ jadwal/ absensi/ penilaian/
    ├── keuangan/ sdm/ perpustakaan/ konsultasi/ surat_izin/
    ├── health_record/ alumni/ portfolio/ hafalan/ ekstra/ event/
```
