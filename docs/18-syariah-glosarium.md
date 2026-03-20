# 18 — Syariah, Role & Glosarium

> **Terakhir diperbarui:** 20 Maret 2026

## Checklist Syariah

- [ ] **Tidak ada riba** — margin/nisbah/ujrah transparan & disepakati sebelum akad
- [ ] **Ta'zir bukan pendapatan** — 100% masuk akun 211 (dana sosial)
- [ ] **Bagi hasil** — dari realisasi pendapatan BMT, bukan % nominal pokok
- [ ] **Bonus wadi'ah** — tidak diperjanjikan di akad, bersifat sukarela
- [ ] **Autodebet partial** — jurnal syariah tetap benar meski partial
- [ ] **Biaya admin rekening** — akad jelas saat pembukaan rekening
- [ ] **Donasi/infaq** — tidak ada iming-iming imbalan materi
- [ ] **Wakaf produktif** — dana tidak bercampur modal BMT, hasil ke mauquf alaih
- [ ] **Komisi OPOP** — akad wakalah bil ujrah, transparan di awal
- [ ] **Denda perpustakaan** — masuk dana sosial, bukan pendapatan
- [ ] **Payroll** — transfer gaji ke rekening, bukan pinjaman
- [ ] **Produk OPOP** — harga transparan, tidak ada gharar
- [ ] **Listing** — hanya direktori & kontak, bukan akad jual beli

## Role Lengkap

| Role | App | Scope |
|------|-----|-------|
| `DEVELOPER` | developer-app | Platform |
| `ADMIN_BMT` | management-app | Seluruh BMT |
| `MANAJER_BMT` | management-app | Seluruh BMT |
| `AUDITOR_BMT` | management-app | Seluruh BMT (read-only) |
| `MANAJER_CABANG` | management-app | 1 cabang |
| `KOMITE` | management-app | 1 cabang |
| `ACCOUNT_OFFICER` | management-app | 1 cabang |
| `FINANCE` | management-app | 1 cabang |
| `TELLER` | teller-app | 1 cabang |
| `NASABAH` | nasabah-app | Data milik sendiri |
| `ALUMNI` | nasabah-app | Data milik sendiri |
| `KASIR_MERCHANT` | merchant-app | 1 merchant |
| `OWNER_MERCHANT` | merchant-app | 1 merchant |
| `OWNER_LISTING` | — (JWT listing) | 1 listing |
| `ADMIN_PONDOK` | pondok-app | 1 cabang |
| `OPERATOR_PONDOK` | pondok-app | 1 cabang |
| `BENDAHARA_PONDOK` | pondok-app | 1 cabang |
| `PETUGAS_UKS` | pondok-app | 1 cabang |
| `PETUGAS_PERPUS` | pondok-app | 1 cabang |
| `KONSELOR` | pondok-app | 1 cabang |
| `PETUGAS_PPDB` | pondok-app | 1 cabang |

## Glosarium

| Istilah | Definisi |
|---------|----------|
| **Add-on** | Fitur à la carte yang dibeli terpisah di atas paket tier |
| **Anti-Fraud** | Rule-based engine deteksi transaksi mencurigakan |
| **Autodebet** | Debit otomatis terjadwal; tanggal per rekening dari DB |
| **Beasiswa** | Potongan biaya SPP/pembiayaan santri, ditetapkan admin pondok |
| **BMT** | Baitul Maal wa Tamwil — koperasi simpan pinjam syariah |
| **CBS** | Core Banking System |
| **DAPODIK** | Data Pokok Pendidikan — sistem data siswa Kemendikbud |
| **EMIS** | Education Management Information System — Kemenag |
| **Feature Gate** | Mekanisme cek apakah fitur aktif di tier/add-on BMT |
| **Finance** | Role staf yang mengelola jurnal manual & biaya operasional |
| **Form Pengajuan** | Mekanisme wajib untuk semua perubahan data nasabah/rekening |
| **GMV** | Gross Merchandise Value — total nilai transaksi OPOP |
| **Hardcode** | Nilai tertanam di kode — **dilarang**, semua dari settings DB |
| **Jenis Rekening** | Tipe rekening dengan aturan & tarif (CRUD management BMT) |
| **Kartu NFC** | Kartu fisik santri untuk transaksi di merchant pondok |
| **Kiosk** | Terminal cek saldo NFC tanpa login |
| **Komisi OPOP** | Biaya platform % dari GMV per pesanan selesai |
| **Listing** | Direktori layanan stakeholder sekitar pondok |
| **Listing Premium** | Listing berbayar — tampil di atas + badge verified |
| **Mauquf Alaih** | Penerima manfaat wakaf |
| **Modul Vernon** | `module-vernon-accounting` — mesin double-entry internal |
| **MRR** | Monthly Recurring Revenue — pendapatan langganan bulanan |
| **Nasabah** | Anggota BMT yang memiliki rekening |
| **Nazir** | Pengelola wakaf produktif (BMT sebagai nazir) |
| **NPSN** | Nomor Pokok Sekolah Nasional |
| **NSM** | Nomor Statistik Madrasah |
| **Offline Mode** | Mode operasi teller & absensi tanpa koneksi internet |
| **OPOP** | One Pondok One Product — marketplace produk pondok |
| **Paket Tier** | Bundle fitur dengan harga tetap (FREE/BASIC/PRO/ENTERPRISE) |
| **Partial Debit** | Debit sebesar saldo tersedia saat autodebet gagal |
| **Payroll** | Proses penggajian otomatis ke rekening BMT karyawan |
| **Pecahan Uang** | Data redenominasi Rupiah di DB — update tanpa deploy |
| **PPDB** | Penerimaan Peserta Didik Baru |
| **RAT** | Rapat Anggota Tahunan |
| **Redenominasi** | Rincian jumlah per pecahan uang saat hitung kas teller |
| **Santri** | Pelajar aktif pondok pesantren |
| **SaaS** | Software as a Service — model bisnis platform ini |
| **Sesi Teller** | Periode kerja teller satu hari dengan pembukuan kas |
| **Settings Engine** | Sistem 3-level resolusi konfigurasi dari DB |
| **SHU** | Sisa Hasil Usaha — "laba" koperasi |
| **SPP** | Sumbangan Pembinaan Pendidikan — iuran santri |
| **Ta'zir** | Denda keterlambatan — 100% masuk dana sosial |
| **Tahfidz** | Hafalan Al-Quran |
| **Tenant** | Satu BMT beserta seluruh cabangnya |
| **Tier** | Paket berlangganan SaaS (FREE/BASIC/PRO/ENTERPRISE) |
| **Tunggakan** | Sisa kewajiban autodebet belum terbayar |
| **Usage Log** | Catatan transaksi yang dikenai biaya admin platform |
| **Wadi'ah** | Akad titipan — BMT pemegang amanah simpanan |
| **Wakaf Produktif** | Aset wakaf dikelola BMT untuk usaha |
| **White-label** | Custom branding app per pondok |
| **Zakat Maal** | Zakat keuntungan BMT yang mencapai nisab & haul |
