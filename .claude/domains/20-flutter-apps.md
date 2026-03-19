# Domain: Flutter Apps (7 Aplikasi)

## 1. `apps-nasabah/` — Android + iOS

**Mode Nasabah Biasa:**
- E-banking: lihat saldo, setor online, riwayat transaksi, pembiayaan

**Mode Wali Santri (jika ada data santri):**
- Semua fitur nasabah biasa +
- Profil & informasi kesiswaan santri
- Saldo kartu NFC santri, riwayat transaksi NFC
- Top-up kartu NFC, kirim ke rekening santri
- Tagihan & bayar SPP pondok
- Raport digital santri
- Notifikasi absensi santri
- **Belanja di toko pondok** (e-commerce B2C)
- Riwayat pesanan & tracking pengiriman

**Mode Alumni:**
- Profil alumni, jaringan alumni
- Belanja OPOP

---

## 2. `apps-management/` — Web + Desktop + Mobile

**Pengguna:** Staf & management BMT
**Role:** `ADMIN_BMT`, `MANAJER_BMT`, `MANAJER_CABANG`, `ACCOUNT_OFFICER`, `KOMITE`, `FINANCE`, `AUDITOR_BMT`

**Fitur:**
- Dashboard & laporan (cabang + konsolidasi)
- Kelola form pengajuan (approve/tolak)
- View nasabah & rekening (edit via form)
- Kelola jenis rekening + biaya admin bulanan
- Kelola produk simpanan & pembiayaan
- Kelola jadwal & konfigurasi autodebet per rekening
- Pembiayaan (analisis, akad, pencairan, monitoring)
- Finance (jurnal manual, biaya operasional, approve slip gaji)
- Settings BMT (settings engine — tidak ada hardcode)
- Laporan RAT

---

## 3. `apps-developer/` — Web + Desktop + Mobile

**Pengguna:** Developer platform
**Auth:** `Developer-Token` header (bukan JWT)

**Fitur:**
- CRUD BMT, kontrak BMT (tarif, fitur aktif)
- **Kelola pecahan uang Rupiah** (data di DB, bukan konstanta)
- Tarif template, usage log
- Platform settings
- Health check, maintenance mode

---

## 4. `apps-teller/` — Desktop

**Pengguna:** Teller cabang BMT
**Aturan kritis:** Semua tombol transaksi **disabled** tanpa sesi kas aktif

**Fitur:**
- Buka sesi kas (redenominasi dari DB secara real-time)
- Tutup sesi (ditolak jika ada selisih — dikonfigurasi `sesi_teller.toleransi_selisih`)
- Transaksi tunai (setor, tarik, angsuran, bayar SPP)
- Cetak slip transaksi
- Buat form pengajuan nasabah

---

## 5. `apps-merchant/` — Android + iOS

**Mode Kasir:**
- Input nominal → tap NFC nasabah → input PIN 6 digit → konfirmasi → cetak struk

**Mode Owner:**
- Dashboard penjualan real-time
- Riwayat transaksi
- Laporan bulanan
- Export CSV

---

## 6. `apps-ceksaldo/` — Android / Kiosk

**Untuk:** Santri cek saldo mandiri

**Alur:** Tap kartu NFC → tampil nama + saldo + 5 transaksi terakhir → auto-reset 10 detik

**Keamanan:**
- Tidak ada PIN
- Tidak ada login
- IP whitelist terminal kiosk (dari DB/settings)
- Endpoint: `GET /nfc/ceksaldo/:uid`

---

## 7. `apps-pondok/` — Web + Mobile

**Pengguna:** Admin pondok (bukan staf BMT)
**Role:** `ADMIN_PONDOK`, `OPERATOR_PONDOK`, `BENDAHARA_PONDOK`, `PETUGAS_UKS`, `PUSTAKAWAN`, `PETUGAS_PPDB`, `BK`

**Fitur lengkap:**
- **Administrasi:** CRUD santri, pengajar, karyawan, alumni
- **Akademik:** Kelola mapel, silabus, RPP, materi ajar
- **Kurikulum:** Komponen penilaian, mapping ke raport
- **Jadwal:** Pelajaran, kegiatan, piket, shift karyawan, kalender akademik
- **Absensi:** Input manual, lihat rekap (metode dari settings BMT)
- **Penilaian:** Input nilai, tahfidz, akhlak, generate raport digital
- **Keuangan:** Tagihan SPP, beasiswa santri, pembiayaan pondok
- **Perpustakaan:** Kelola buku, peminjaman
- **Konsultasi:** Sesi konsultasi online santri/wali
- **Surat Izin:** Approve/tolak izin keluar santri
- **UKS:** Kunjungan kesehatan, health record santri
- **PPDB:** Proses pendaftaran santri baru
- **OPOP Toko:** Kelola produk & pesanan toko pondok
- **Laporan:** Kesiswaan, akademik, keuangan pondok
- **Integrasi:** Sinkronisasi DAPODIK & EMIS
