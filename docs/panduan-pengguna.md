# Panduan Pengguna — Platform Pesantren Terpadu

> **Versi:** 1.0 · **Terakhir diperbarui:** 20 Maret 2026
> Platform digital terpadu untuk pondok pesantren: CBS Syariah · ERP Pondok · E-commerce OPOP

---

## Daftar Isi

1. [Pendahuluan](#1-pendahuluan)
2. [Cara Login & Keamanan Akun](#2-cara-login--keamanan-akun)
3. [Aplikasi Nasabah](#3-aplikasi-nasabah-androidios)
   - [3.1 Beranda & Saldo](#31-beranda--saldo)
   - [3.2 Riwayat Transaksi (Mutasi)](#32-riwayat-transaksi-mutasi)
   - [3.3 Transfer Antar Rekening](#33-transfer-antar-rekening)
   - [3.4 Kartu NFC Virtual — Top-up & Cek Saldo](#34-kartu-nfc-virtual--top-up--cek-saldo)
   - [3.5 Belanja di OPOP](#35-belanja-di-opop-e-commerce-pondok)
   - [3.6 Melihat Raport & Akademik Anak](#36-melihat-raport--akademik-anak-wali-santri)
   - [3.7 Pembayaran SPP](#37-pembayaran-spp)
4. [Aplikasi Teller](#4-aplikasi-teller-desktop)
   - [4.1 Membuka Sesi Kas](#41-membuka-sesi-kas)
   - [4.2 Setor Tunai Nasabah](#42-setor-tunai-nasabah)
   - [4.3 Tarik Tunai Nasabah](#43-tarik-tunai-nasabah)
   - [4.4 Transfer Antar Rekening](#44-transfer-antar-rekening)
   - [4.5 Cetak Slip Transaksi](#45-cetak-slip-transaksi)
   - [4.6 Menutup Sesi Kas](#46-menutup-sesi-kas)
   - [4.7 Pencarian Nasabah](#47-pencarian-nasabah)
5. [Aplikasi Management](#5-aplikasi-management-web)
   - [5.1 Dashboard & Laporan Singkat](#51-dashboard--laporan-singkat)
   - [5.2 Manajemen Nasabah](#52-manajemen-nasabah)
   - [5.3 Pembiayaan (Kredit Syariah)](#53-pembiayaan-kredit-syariah)
   - [5.4 Laporan Keuangan](#54-laporan-keuangan)
   - [5.5 Jurnal Manual](#55-jurnal-manual-staf-finance)
   - [5.6 Settings BMT](#56-settings-bmt)
6. [Aplikasi Pondok](#6-aplikasi-pondok-webmobile)
   - [6.1 Manajemen Santri](#61-manajemen-santri)
   - [6.2 Absensi Kelas](#62-absensi-kelas)
   - [6.3 Input Nilai & Raport](#63-input-nilai--raport)
   - [6.4 Tagihan SPP](#64-tagihan-spp)
   - [6.5 Perpustakaan](#65-perpustakaan)
   - [6.6 PPDB](#66-ppdb-penerimaan-peserta-didik-baru)
7. [Aplikasi Merchant](#7-aplikasi-merchant-androidios)
   - [7.1 Transaksi NFC (Kasir)](#71-transaksi-nfc-kasir)
   - [7.2 Pembayaran Manual (tanpa NFC)](#72-pembayaran-manual-tanpa-nfc)
   - [7.3 Laporan Penjualan Harian](#73-laporan-penjualan-harian)
8. [Kiosk Cek Saldo](#8-kiosk-cek-saldo)
9. [Notifikasi](#9-notifikasi)
10. [Fitur Autodebet](#10-fitur-autodebet)
11. [Keamanan & Privasi](#11-keamanan--privasi)
12. [Pertanyaan Umum (FAQ)](#12-pertanyaan-umum-faq)
13. [Kontak & Dukungan](#13-kontak--dukungan)

---

## 1. Pendahuluan

Platform Pesantren Terpadu adalah sistem digital yang dirancang khusus untuk kebutuhan pondok pesantren modern di Indonesia. Platform ini menyatukan tiga layanan utama dalam satu ekosistem yang terintegrasi:

- **🏦 CBS Syariah (Core Banking System)** — Layanan simpan-pinjam syariah: tabungan, simpanan berjangka (deposito), pembiayaan (kredit syariah), transaksi teller, dan autodebet.
- **🏫 ERP Pondok** — Pengelolaan operasional pesantren: data santri, absensi, nilai akademik, raport, perpustakaan, PPDB, dan tagihan SPP.
- **🛒 E-commerce OPOP** — Toko online produk-produk pondok pesantren, mendukung program One Pesantren One Product.

Platform ini hadir dalam **7 aplikasi** yang masing-masing dirancang sesuai peran penggunanya, mulai dari nasabah dan wali santri, staf teller, pengelola pondok, hingga merchant kantin dan koperasi.

### Manfaat Utama

- **Satu login, semua layanan** — Nasabah, wali santri, dan santri cukup satu akun untuk mengakses tabungan, data akademik, dan toko.
- **Transaksi real-time** — Setoran, penarikan, dan transfer langsung tercatat dan dapat dipantau.
- **Paperless & audit trail** — Semua transaksi tersimpan digital dengan jejak audit lengkap.
- **Sesuai prinsip syariah** — Seluruh perhitungan, bagi hasil, dan alur transaksi mengikuti ketentuan syariah Islam.
- **Aman & terpisah antar lembaga** — Data setiap BMT dan pondok sepenuhnya terisolasi dari lembaga lain.

---

## 2. Cara Login & Keamanan Akun

### Login Pertama Kali

Akun Anda didaftarkan oleh admin BMT atau pondok. Anda akan menerima notifikasi (WhatsApp atau SMS) berisi informasi akun awal.

**Langkah login:**

1. Buka aplikasi yang sesuai dengan peran Anda.
2. Masukkan **nomor telepon** atau **email** yang terdaftar.
3. Ketuk **Kirim OTP** — kode 6 digit akan dikirim via WhatsApp atau SMS.
4. Masukkan kode OTP dalam waktu 5 menit.
5. Pada login pertama, Anda akan diminta membuat **PIN** (6 digit) untuk keamanan tambahan.

> **⚠️ Penting:** OTP hanya berlaku **5 menit**. Jika kode sudah kadaluarsa, ketuk "Kirim Ulang OTP" untuk mendapatkan kode baru.

### Verifikasi 2 Langkah (2FA)

Untuk akun dengan akses sensitif (teller, manager, finance), platform menerapkan verifikasi dua langkah secara otomatis. Setelah memasukkan OTP, Anda juga perlu mengonfirmasi identitas melalui PIN yang sudah Anda buat.

### Reset Password / Lupa PIN

Jika Anda lupa PIN atau password:

1. Di halaman login, ketuk **"Lupa PIN?"** atau **"Lupa Password?"**.
2. Masukkan nomor telepon atau email Anda.
3. Kode verifikasi akan dikirim ke nomor/email terdaftar.
4. Masukkan kode, lalu buat PIN/password baru.

Jika nomor telepon Anda sudah berganti, hubungi admin BMT atau admin pondok untuk mereset akun secara manual.

### Logout Aman

Selalu logout setelah selesai menggunakan aplikasi, terutama di perangkat yang digunakan bersama.

- Di pojok kanan atas (atau menu profil), ketuk ikon **keluar** atau tombol **Logout**.
- Aplikasi akan mengakhiri sesi dan menghapus token login dari perangkat.

> **Tips keamanan:** Jangan simpan password di browser publik atau perangkat orang lain.

### Keamanan OTP — Peringatan Penting

> **🚨 JANGAN PERNAH bagikan OTP kepada siapapun** — termasuk kepada orang yang mengaku staf BMT, admin pondok, atau customer service platform. Pihak resmi platform **tidak pernah** meminta OTP melalui telepon, WhatsApp, atau pesan apapun.

---

## 3. Aplikasi Nasabah (Android/iOS)

📱 Aplikasi Nasabah diperuntukkan bagi **nasabah BMT**, **wali santri**, dan **alumni**. Tersedia di Google Play Store dan Apple App Store.

### 3.1 Beranda & Saldo

Setelah login, halaman **Beranda** menampilkan ringkasan akun Anda:

- **Saldo aktif** semua rekening yang Anda miliki (tabungan sukarela, simpanan wajib, deposito, dll.).
- **Notifikasi terbaru** — transaksi masuk/keluar, tagihan jatuh tempo, pengumuman pondok.
- **Pintasan cepat** — Transfer, Top-up Kartu, Bayar SPP, Lihat Mutasi.

Ketuk nama rekening untuk melihat detail, termasuk nomor rekening lengkap dan produk simpanan yang berlaku.

### 3.2 Riwayat Transaksi (Mutasi)

1. Dari Beranda, pilih rekening yang ingin dilihat mutasinya.
2. Ketuk **"Mutasi"** atau **"Riwayat Transaksi"**.
3. Mutasi ditampilkan urut terbaru. Setiap baris menampilkan:
   - Tanggal & waktu transaksi
   - Keterangan (setor, tarik, transfer, autodebet, dll.)
   - Jumlah (hijau = masuk, merah = keluar)
   - Saldo akhir setelah transaksi

**Filter mutasi:** Gunakan tombol filter di pojok kanan atas untuk menyaring berdasarkan tanggal atau jenis transaksi.

**Unduh mutasi:** Ketuk ikon unduh untuk mengekspor mutasi dalam format PDF (tersedia untuk periode maksimal 3 bulan terakhir).

### 3.3 Transfer Antar Rekening

Transfer dapat dilakukan antar rekening dalam satu BMT (cabang yang sama maupun berbeda).

**Langkah transfer:**

1. Dari Beranda, ketuk **"Transfer"**.
2. Pilih rekening sumber (rekening Anda yang akan didebet).
3. Masukkan **nomor rekening tujuan** atau cari berdasarkan nama nasabah.
4. Konfirmasi nama pemilik rekening tujuan yang muncul di layar.
5. Masukkan **jumlah** yang akan ditransfer.
6. Tambahkan **catatan** (opsional, contoh: "Titipan uang saku").
7. Ketuk **Lanjut** → periksa ringkasan transaksi → ketuk **Konfirmasi**.
8. Masukkan **PIN** Anda untuk otorisasi.
9. Slip transaksi digital akan ditampilkan. Ketuk **Bagikan** untuk mengirim bukti via WhatsApp.

> **Catatan:** Transfer antar rekening dalam platform ini tidak dikenakan biaya (sesuai kebijakan masing-masing BMT). Untuk transfer ke rekening bank lain di luar platform, hubungi teller di kantor cabang.

### 3.4 Kartu NFC Virtual — Top-up & Cek Saldo

Kartu NFC digunakan santri untuk bertransaksi di kantin, koperasi, dan merchant pondok tanpa uang tunai.

**Melihat saldo kartu NFC:**

- Di halaman Beranda, ketuk ikon **Kartu NFC** atau menu **"Kartu Saya"**.
- Saldo kartu ditampilkan beserta 5 transaksi terakhir.

**Top-up saldo kartu NFC:**

1. Di halaman Kartu NFC, ketuk **"Top-up"**.
2. Pilih sumber dana (rekening tabungan Anda).
3. Masukkan jumlah top-up (minimal sesuai ketentuan BMT).
4. Konfirmasi dengan PIN.
5. Saldo kartu akan bertambah secara instan.

> **Tips:** Atur top-up otomatis agar saldo kartu santri selalu tercukupi. Fitur ini dapat diaktifkan di menu **Pengaturan → Auto Top-up Kartu**.

**Pemblokiran kartu NFC:**

Jika kartu hilang atau dicuri, segera blokir melalui:
- Aplikasi Nasabah → Kartu Saya → **"Blokir Kartu"**, atau
- Hubungi admin BMT/pondok untuk pemblokiran manual.

### 3.5 Belanja di OPOP (E-commerce Pondok)

OPOP (One Pesantren One Product) adalah toko online produk-produk hasil pondok pesantren.

**Cara belanja:**

1. Dari menu bawah, ketuk **"OPOP"** atau **"Toko"**.
2. Jelajahi produk berdasarkan kategori atau gunakan pencarian.
3. Ketuk produk untuk melihat detail, foto, deskripsi, dan harga.
4. Ketuk **"Tambah ke Keranjang"** lalu lanjut ke **"Checkout"**.
5. Pilih alamat pengiriman dan metode pengiriman.
6. Pilih metode pembayaran:
   - **Saldo Rekening** — langsung didebet dari rekening tabungan.
   - **Kartu NFC** — menggunakan saldo kartu.
7. Konfirmasi pesanan → masukkan PIN → pesanan berhasil dibuat.
8. Pantau status pesanan di menu **"Pesanan Saya"**.

### 3.6 Melihat Raport & Akademik Anak (Wali Santri)

Wali santri dapat memantau perkembangan akademik anak langsung dari aplikasi.

**Cara mengakses:**

1. Di Beranda, ketuk tab **"Anak Saya"** atau menu **"Akademik"**.
2. Jika memiliki lebih dari satu anak yang terdaftar, pilih nama santri.
3. Tersedia informasi:
   - **Absensi** — rekap kehadiran bulan ini dan riwayat absensi.
   - **Nilai** — nilai per mata pelajaran dan ujian.
   - **Raport** — raport semester yang sudah diterbitkan (format PDF).
   - **Pengumuman Kelas** — dari wali kelas atau ustadz.
   - **Hafalan** — progress hafalan Al-Qur'an (jika pondok menggunakan fitur ini).

> **Catatan:** Data akademik hanya tampil setelah admin pondok menerbitkannya. Jika raport belum muncul, berarti belum diterbitkan oleh pihak pondok.

### 3.7 Pembayaran SPP

**Cara bayar SPP:**

1. Dari Beranda, ketuk **"Bayar SPP"** atau masuk ke menu **"Tagihan"**.
2. Sistem akan menampilkan tagihan SPP yang belum terbayar beserta rincian (bulan, jumlah, batas waktu).
3. Pilih tagihan yang ingin dibayar.
4. Pilih sumber dana (rekening tabungan).
5. Konfirmasi → masukkan PIN → pembayaran selesai.
6. Bukti pembayaran dapat diunduh atau dibagikan via WhatsApp.

> **Info:** SPP yang sudah lewat jatuh tempo akan ditandai merah. Hubungi admin pondok jika ada perbedaan tagihan.

---

## 4. Aplikasi Teller (Desktop)

🏦 Aplikasi Teller digunakan oleh **petugas teller** di kantor cabang BMT untuk melayani transaksi tunai nasabah. Aplikasi ini berjalan di komputer desktop.

### 4.1 Membuka Sesi Kas

Sebelum melayani nasabah, teller wajib membuka sesi kas terlebih dahulu. Sesi kas adalah periode kerja teller yang terikat dengan jumlah uang fisik yang dipegang.

**Langkah membuka sesi kas:**

1. Login ke Aplikasi Teller dengan akun Anda.
2. Dari halaman utama, ketuk **"Buka Sesi Kas"**.
3. Sistem akan menampilkan formulir **hitung uang fisik per pecahan**:

| Pecahan | Jumlah Lembar/Keping | Subtotal |
|---------|---------------------|---------|
| Rp 100.000 | ___ | Rp ___ |
| Rp 50.000 | ___ | Rp ___ |
| Rp 20.000 | ___ | Rp ___ |
| Rp 10.000 | ___ | Rp ___ |
| Rp 5.000 | ___ | Rp ___ |
| Rp 2.000 | ___ | Rp ___ |
| Rp 1.000 | ___ | Rp ___ |
| Koin | ___ | Rp ___ |

4. Isi **jumlah lembar/keping** untuk setiap pecahan (bukan total langsung). Sistem menghitung total otomatis.
5. Periksa total yang muncul — pastikan sesuai dengan uang fisik yang Anda hitung.
6. Ketuk **"Konfirmasi Buka Sesi"**.

> **⚠️ Wajib:** Pengisian harus per pecahan, bukan langsung memasukkan angka total. Ini memastikan akurasi hitung dan mempermudah audit.

Sesi kas yang berhasil dibuka akan ditampilkan di sudut kiri atas layar beserta waktu mulai dan saldo awal kas.

### 4.2 Setor Tunai Nasabah

**Langkah melayani setoran:**

1. Dari menu utama, pilih **"Setor Tunai"**.
2. Cari nasabah berdasarkan nomor rekening, NIK, atau nama (lihat [Pencarian Nasabah](#47-pencarian-nasabah)).
3. Konfirmasi identitas nasabah — cocokkan dengan KTP atau buku tabungan.
4. Masukkan **jumlah setoran** yang diterima dari nasabah.
5. Pilih rekening tujuan (jika nasabah punya lebih dari satu rekening).
6. Periksa ringkasan transaksi:
   - Nama & nomor rekening nasabah
   - Jumlah setoran
   - Saldo sebelum dan sesudah setoran
7. Ketuk **"Proses Setoran"**.
8. Sistem otomatis mencetak slip transaksi (atau tampilkan di layar untuk ditandatangani nasabah).

### 4.3 Tarik Tunai Nasabah

**Langkah melayani penarikan:**

1. Pilih **"Tarik Tunai"** dari menu utama.
2. Cari dan verifikasi identitas nasabah.
3. Masukkan **jumlah penarikan**.
4. Sistem akan memvalidasi:
   - Saldo mencukupi (saldo ≥ jumlah tarik + saldo minimum rekening).
   - Tidak melebihi limit penarikan harian.
5. Periksa ringkasan → ketuk **"Proses Penarikan"**.
6. Serahkan uang tunai kepada nasabah setelah transaksi berhasil.
7. Cetak slip transaksi.

> **Catatan:** Jika nasabah ingin menarik seluruh saldo (tutup rekening), proses ini harus melalui menu **Manajemen Nasabah** di Aplikasi Management, bukan dari menu Tarik Tunai biasa.

### 4.4 Transfer Antar Rekening

Teller dapat membantu nasabah melakukan transfer antar rekening.

1. Pilih **"Transfer"** dari menu utama.
2. Cari rekening nasabah sebagai sumber dana.
3. Masukkan nomor rekening tujuan.
4. Konfirmasi nama pemilik rekening tujuan.
5. Masukkan jumlah transfer dan keterangan.
6. Proses dan cetak slip.

### 4.5 Cetak Slip Transaksi

Slip transaksi dicetak otomatis setelah setiap transaksi berhasil. Jika printer tidak siap atau slip gagal tercetak:

1. Buka menu **"Riwayat Transaksi"** di sesi aktif.
2. Cari transaksi yang dimaksud.
3. Ketuk **"Cetak Ulang Slip"**.

Slip berisi: nomor transaksi, tanggal/waktu, nama nasabah, nomor rekening, jenis transaksi, jumlah, saldo akhir, dan nama teller.

### 4.6 Menutup Sesi Kas

Di akhir hari kerja atau akhir shift, teller wajib menutup sesi kas.

**Langkah menutup sesi kas:**

1. Pastikan semua transaksi nasabah sudah selesai.
2. Hitung fisik uang kas yang tersisa.
3. Dari menu utama, ketuk **"Tutup Sesi Kas"**.
4. Isi formulir hitung uang fisik per pecahan (sama seperti saat buka sesi).
5. Sistem akan menampilkan:
   - **Saldo awal** (saat buka sesi)
   - **Total setoran masuk** (selama sesi)
   - **Total penarikan keluar** (selama sesi)
   - **Saldo kas seharusnya** (perhitungan sistem)
   - **Saldo kas aktual** (hasil hitung fisik Anda)
   - **Selisih** (jika ada)

6. Jika **selisih = Rp 0**: Ketuk **"Konfirmasi Tutup Sesi"**. Sesi ditutup.
7. Jika **ada selisih**: Sistem akan memblokir penutupan dan meminta Anda menghitung ulang uang fisik.

> **⚠️ Jika selisih tetap ada setelah hitung ulang:** Hubungi supervisor/kepala teller. Selisih akan dicatat sebagai **Ta'zir** (denda) sesuai ketentuan BMT dan dimasukkan ke jurnal selisih kas.

### 4.7 Pencarian Nasabah

Gunakan fitur pencarian untuk menemukan nasabah dengan cepat:

- **Nomor rekening** — cara tercepat dan paling akurat.
- **NIK (Nomor KTP)** — masukkan 16 digit NIK.
- **Nama** — masukkan minimal 3 karakter nama. Hasil pencarian akan menampilkan daftar nasabah yang cocok.
- **Nomor HP** — masukkan nomor HP yang terdaftar.

Setelah menemukan nasabah, selalu **konfirmasi identitas fisik** (KTP atau buku tabungan) sebelum memproses transaksi.

---

## 5. Aplikasi Management (Web)

🖥️ Aplikasi Management diakses melalui browser web, diperuntukkan bagi **staf BMT, manajer, staf keuangan, dan account officer**.

### 5.1 Dashboard & Laporan Singkat

Halaman Dashboard menampilkan ringkasan kondisi BMT secara real-time:

- **Total Aset** — jumlah aset per hari ini.
- **DPK (Dana Pihak Ketiga)** — total simpanan nasabah.
- **Outstanding Pembiayaan** — total pembiayaan yang masih berjalan.
- **Likuiditas** — posisi kas dan setara kas.
- **Grafik tren** — perkembangan DPK dan pembiayaan 30 hari terakhir.
- **Teller aktif** — daftar teller yang sedang membuka sesi hari ini.
- **Transaksi hari ini** — jumlah dan volume transaksi yang sudah diproses.

Widget dashboard dapat dikustomisasi sesuai peran pengguna.

### 5.2 Manajemen Nasabah

#### Membuka Rekening Baru

> **Penting:** Semua perubahan data nasabah dan pembukaan rekening wajib melalui **formulir pengajuan** yang disetujui oleh pejabat berwenang. Tidak ada edit langsung.

**Langkah membuka rekening baru:**

1. Masuk ke menu **"Nasabah"** → **"Pengajuan Rekening Baru"**.
2. Isi formulir data calon nasabah:
   - Data pribadi (NIK, nama lengkap, tempat/tanggal lahir, alamat)
   - Data kontak (nomor HP aktif, email)
   - Jenis rekening yang diajukan
   - Setoran awal
3. Unggah dokumen: scan KTP, foto nasabah, tanda tangan digital.
4. Ketuk **"Ajukan"** → formulir masuk ke antrean persetujuan.
5. Pejabat berwenang (kepala cabang atau supervisor) akan menyetujui atau menolak.
6. Setelah disetujui, rekening otomatis aktif dan nasabah menerima notifikasi.

#### Update Data Nasabah

Perubahan data nasabah (alamat, nomor HP, status pernikahan, dll.) tidak dapat dilakukan langsung — harus melalui form pengajuan perubahan data:

1. Cari nasabah → ketuk **"Ajukan Perubahan Data"**.
2. Isi data yang berubah dan alasan perubahan.
3. Unggah dokumen pendukung (KTP baru, surat keterangan, dll.).
4. Ajukan → tunggu persetujuan.

> **Mengapa harus melalui form?** Ini adalah mekanisme keamanan untuk mencegah perubahan data tidak sah yang bisa digunakan untuk fraud atau pembobolan rekening.

#### Pencarian & Detail Nasabah

- Cari nasabah via nomor rekening, NIK, nama, atau nomor HP.
- Halaman detail nasabah menampilkan: data pribadi, daftar rekening, riwayat transaksi, status pembiayaan, dan histori perubahan data.

### 5.3 Pembiayaan (Kredit Syariah)

#### Pengajuan Pembiayaan Baru

1. Masuk ke menu **"Pembiayaan"** → **"Pengajuan Baru"**.
2. Cari nasabah yang mengajukan.
3. Isi formulir pengajuan:
   - Jenis pembiayaan (Murabahah, Mudharabah, Musyarakah, dll.)
   - Jumlah yang diajukan
   - Tujuan penggunaan dana
   - Jangka waktu dan pola angsuran
   - Agunan (jika ada)
4. Ajukan → status berubah menjadi **"Menunggu Analisis"**.

#### Alur Persetujuan Pembiayaan

Setiap pengajuan pembiayaan melewati empat tahap:

```
Pengajuan → Analisis (AO) → Komite → Akad → Pencairan
```

| Tahap | Pelaku | Keterangan |
|-------|--------|-----------|
| **Analisis** | Account Officer (AO) | AO melakukan survei, analisis karakter & kapasitas nasabah, menyiapkan memorandum analisis |
| **Komite** | Komite Pembiayaan | Rapat komite memutuskan: setuju, ditolak, atau setuju dengan syarat |
| **Akad** | Staf Legal & Nasabah | Penandatanganan akad pembiayaan sesuai jenis akad syariah |
| **Pencairan** | Teller / Finance | Dana dicairkan ke rekening nasabah atau pihak ketiga (untuk Murabahah) |

Setiap tahap dapat dilihat statusnya di menu **"Pembiayaan"** → **"Daftar Pengajuan"**.

#### Monitoring Angsuran & Kolektibilitas

- **Jadwal angsuran** — lihat jadwal lengkap angsuran per pembiayaan beserta status pembayaran (lunas, belum, terlambat).
- **Kolektibilitas** — tingkat kelancaran pembayaran sesuai standar OJK (5 level):
  - **1 - Lancar** — tidak ada tunggakan.
  - **2 - Dalam Perhatian Khusus** — tunggakan 1–90 hari.
  - **3 - Kurang Lancar** — tunggakan 91–120 hari.
  - **4 - Diragukan** — tunggakan 121–180 hari.
  - **5 - Macet** — tunggakan lebih dari 180 hari.
- **PPAP (Penyisihan Penghapusan Aktiva Produktif)** — sistem menghitung kebutuhan PPAP otomatis berdasarkan kolektibilitas.

#### Pembayaran Angsuran

Angsuran dapat dibayar melalui:
- **Autodebet** — ditarik otomatis dari rekening tabungan nasabah pada tanggal jatuh tempo.
- **Bayar via Teller** — nasabah datang ke kantor dan bayar tunai.
- **Bayar via Aplikasi Nasabah** — transfer dari rekening ke nomor rekening angsuran.

### 5.4 Laporan Keuangan

Semua laporan dapat diekspor ke PDF atau Excel.

#### Neraca (Balance Sheet)

Menampilkan posisi aset, kewajiban, dan ekuitas BMT pada tanggal tertentu. Akses: **Laporan → Neraca** → pilih tanggal laporan.

#### SHU (Sisa Hasil Usaha)

Laporan pendapatan dan beban operasional BMT dalam periode tertentu. Setara laporan laba-rugi di koperasi syariah.

#### Laporan Arus Kas

Menampilkan aliran kas dari tiga aktivitas:
- Aktivitas Operasi
- Aktivitas Investasi
- Aktivitas Pendanaan

#### Laporan Kolektibilitas

Laporan portofolio pembiayaan yang dikelompokkan berdasarkan 5 level kolektibilitas OJK. Wajib dilaporkan ke OJK secara berkala.

#### Laporan Bagi Hasil Deposito

Menampilkan perhitungan bagi hasil (nisbah) untuk nasabah deposito berdasarkan realisasi pendapatan BMT, bukan proyeksi/estimasi — sesuai prinsip syariah.

### 5.5 Jurnal Manual (Staf Finance)

Untuk koreksi akuntansi atau pencatatan transaksi non-operasional, staf finance dapat membuat jurnal manual.

**Membuat jurnal manual:**

1. Masuk ke menu **"Akuntansi"** → **"Jurnal Manual"** → **"Buat Jurnal"**.
2. Isi tanggal efektif jurnal.
3. Tambahkan baris debet dan kredit:
   - Pilih akun dari Bagan Akun (Chart of Accounts).
   - Masukkan jumlah.
   - Tambahkan keterangan per baris.
4. Sistem memvalidasi: **Σ Debet harus = Σ Kredit**. Jika tidak seimbang, jurnal tidak dapat disimpan.
5. Ketuk **"Simpan Draft"**.

**Menyetujui (posting) jurnal:**

Jurnal draft tidak langsung aktif — harus disetujui oleh pejabat berwenang (kepala accounting atau manager):

1. Masuk ke **"Jurnal Manual"** → **"Menunggu Persetujuan"**.
2. Review jurnal yang diajukan.
3. Ketuk **"Setujui & Posting"** — jurnal resmi tercatat di buku besar.

> **Catatan:** Jurnal yang sudah diposting tidak dapat dihapus. Koreksi hanya bisa dilakukan dengan jurnal koreksi baru (jurnal balik/reverse).

### 5.6 Settings BMT

Menu Settings hanya dapat diakses oleh admin/manager BMT. Semua konfigurasi sistem disimpan di database — tidak ada hardcode di sistem.

**Yang dapat dikonfigurasi:**

- **Produk Simpanan** — nama produk, nisbah bagi hasil, setoran minimum, biaya admin.
- **Produk Pembiayaan** — jenis akad, margin, biaya-biaya.
- **Autodebet** — jadwal, retry, kebijakan partial debit.
- **Notifikasi** — template pesan WhatsApp, SMS, dan email (dapat dikustomisasi per BMT).
- **Limit Transaksi** — limit harian per nasabah, per teller.
- **Jadwal Operasional** — jam buka dan tutup cabang.
- **Bagan Akun** — nomor dan nama akun sesuai standar akuntansi BMT.

---

## 6. Aplikasi Pondok (Web/Mobile)

🏫 Aplikasi Pondok digunakan oleh **admin pondok, operator, ustadz/ustadzah, dan staf perpustakaan** untuk mengelola kegiatan pesantren.

### 6.1 Manajemen Santri

**Melihat daftar santri:**

1. Masuk ke menu **"Santri"** → **"Daftar Santri"**.
2. Filter berdasarkan: kelas, angkatan, status aktif/alumni, asrama.
3. Ketuk nama santri untuk melihat profil lengkap.

**Menambah data santri:**

1. Ketuk **"+ Tambah Santri"**.
2. Isi formulir: data pribadi, asal sekolah, data orang tua/wali, kelas, asrama.
3. Unggah foto santri dan dokumen (akta lahir, kartu keluarga).
4. Simpan → santri masuk ke sistem dan nomor induk santri diterbitkan otomatis.

**Mengubah data santri:**

Sama seperti nasabah BMT, perubahan data santri harus melalui formulir perubahan yang disetujui kepala pondok atau admin.

**Mutasi santri** (pindah kelas, pindah asrama, keluar):

Gunakan menu **"Mutasi Santri"** — setiap perubahan status tercatat lengkap dengan tanggal dan alasan.

### 6.2 Absensi Kelas

**Input absensi harian:**

1. Masuk ke menu **"Absensi"** → **"Input Absensi"**.
2. Pilih kelas dan tanggal.
3. Untuk setiap santri, pilih status: **Hadir / Sakit / Izin / Tanpa Keterangan (Alpha)**.
4. Untuk izin/sakit, lampirkan surat keterangan (opsional).
5. Ketuk **"Simpan Absensi"**.

Absensi yang sudah disimpan dapat dikoreksi oleh admin dalam batas waktu yang ditentukan (biasanya H+1).

**Rekap absensi:**

- Per santri: lihat dari profil santri → tab "Absensi".
- Per kelas: menu **"Absensi"** → **"Rekap"** → pilih kelas dan periode.
- Ekspor ke Excel untuk pelaporan.

> **Info untuk Wali Santri:** Jika anak Anda tercatat alpha (tanpa keterangan), Anda dapat melihatnya langsung di Aplikasi Nasabah dan menghubungi wali kelas.

### 6.3 Input Nilai & Raport

**Input nilai:**

1. Masuk ke menu **"Akademik"** → **"Input Nilai"**.
2. Pilih kelas, mata pelajaran, dan jenis penilaian (Ulangan Harian, UTS, UAS, Tugas, dll.).
3. Masukkan nilai untuk setiap santri.
4. Simpan.

**Menerbitkan raport:**

1. Pastikan semua nilai sudah lengkap untuk semua mata pelajaran.
2. Masuk ke menu **"Raport"** → **"Generate Raport"**.
3. Pilih kelas dan semester.
4. Sistem akan menghasilkan raport PDF otomatis berdasarkan nilai yang sudah diinput.
5. Review raport → ketuk **"Terbitkan"**.
6. Raport langsung dapat diakses oleh wali santri melalui Aplikasi Nasabah.

### 6.4 Tagihan SPP

**Membuat tagihan SPP:**

1. Masuk ke menu **"Keuangan"** → **"Tagihan SPP"** → **"Buat Tagihan"**.
2. Pilih periode (bulan dan tahun).
3. Pilih kelas atau pilih "Semua Santri Aktif".
4. Sistem akan membuat tagihan untuk setiap santri secara otomatis berdasarkan besaran SPP di kelasnya.
5. Konfirmasi → tagihan aktif dan muncul di Aplikasi Nasabah wali santri.

**Pemantauan pembayaran SPP:**

- Menu **"Keuangan"** → **"Status SPP"** menampilkan daftar santri beserta status pembayaran (Lunas / Belum Bayar / Terlambat).
- Filter berdasarkan kelas, bulan, atau status pembayaran.

### 6.5 Perpustakaan

**Manajemen koleksi:**

1. Masuk ke menu **"Perpustakaan"** → **"Koleksi"** → **"+ Tambah Buku"**.
2. Isi: judul, pengarang, penerbit, tahun terbit, nomor ISBN, jumlah eksemplar, lokasi rak.
3. Sistem menerbitkan nomor inventaris otomatis.

**Peminjaman buku:**

1. Menu **"Transaksi"** → **"Pinjam Buku"**.
2. Scan atau ketik nomor anggota peminjam.
3. Scan atau ketik nomor inventaris buku.
4. Sistem otomatis menentukan tanggal kembali (sesuai kebijakan pondok, contoh: 7 hari).
5. Konfirmasi → buku tercatat dipinjam.

**Pengembalian buku:**

1. Menu **"Transaksi"** → **"Kembali Buku"**.
2. Scan nomor inventaris buku.
3. Sistem menampilkan data peminjam dan tanggal kembali.
4. Jika terlambat, sistem menghitung denda otomatis.
5. Konfirmasi pengembalian.

### 6.6 PPDB (Penerimaan Peserta Didik Baru)

**Membuka gelombang PPDB:**

1. Masuk ke menu **"PPDB"** → **"Gelombang"** → **"Buka Gelombang Baru"**.
2. Isi: nama gelombang, tanggal buka & tutup pendaftaran, kuota, biaya pendaftaran.
3. Aktifkan → formulir PPDB online tersedia untuk umum.

**Memproses pendaftar:**

1. Menu **"PPDB"** → **"Pendaftar"** — daftar semua pendaftar yang masuk.
2. Untuk setiap pendaftar: lihat berkas, verifikasi dokumen, beri status (Lengkap / Perlu Perbaikan).
3. Setelah verifikasi: proses seleksi, tes masuk (jika ada), pengumuman kelulusan.
4. Pendaftar yang diterima dapat langsung dikonversi menjadi data santri.

---

## 7. Aplikasi Merchant (Android/iOS)

🛒 Aplikasi Merchant digunakan oleh **kasir dan pemilik merchant** (kantin, koperasi, toko) di lingkungan pondok pesantren.

### 7.1 Transaksi NFC (Kasir)

Santri membayar pembelian cukup dengan menempelkan kartu NFC ke perangkat kasir — tanpa uang tunai.

**Cara memproses transaksi NFC:**

1. Kasir membuat daftar belanja di aplikasi dan konfirmasi total tagihan.
2. Minta santri **menempelkan kartu NFC** ke bagian belakang smartphone kasir (atau ke reader NFC eksternal).
3. Kartu terbaca → aplikasi menampilkan:
   - Nama santri
   - Saldo kartu saat ini
   - Total tagihan yang harus dibayar
4. Kasir meminta santri **konfirmasi** ("Apakah Anda setuju membayar Rp X?").
5. Santri **menekan tombol konfirmasi** di layar (atau menginput PIN NFC jika nominal melebihi batas tanpa PIN).
6. Transaksi berhasil → saldo kartu berkurang → struk digital muncul di layar.
7. Ketuk **"Cetak Struk"** (jika tersedia printer) atau **"Selesai"**.

> **Catatan keamanan:** Jika saldo kartu santri tidak mencukupi, transaksi otomatis ditolak dan kasir akan mendapat notifikasi. Tidak ada utang atau overdraft dari kartu NFC.

### 7.2 Pembayaran Manual (tanpa NFC)

Jika kartu santri tidak tersedia atau rusak:

1. Di halaman kasir, ketuk **"Bayar Manual"**.
2. Masukkan **nomor rekening** atau **nomor kartu** santri secara manual.
3. Konfirmasi nama santri yang muncul.
4. Lanjutkan proses seperti transaksi NFC biasa.

Cara ini memerlukan santri menginput PIN di layar aplikasi kasir.

### 7.3 Laporan Penjualan Harian

**Melihat laporan harian:**

1. Masuk ke menu **"Laporan"** → **"Penjualan Hari Ini"**.
2. Laporan menampilkan:
   - Jumlah transaksi
   - Total omzet
   - Daftar transaksi per jam
   - Produk terlaris
3. Untuk laporan periode lain: pilih **"Laporan Periode"** dan tentukan tanggal awal–akhir.
4. Ekspor ke PDF untuk diserahkan ke pengelola toko.

**Rekonsiliasi:**

Pengelola merchant dapat membandingkan total penjualan di aplikasi dengan catatan manual (jika ada) melalui menu **"Rekonsiliasi"**.

---

## 8. Kiosk Cek Saldo

💳 Kiosk Cek Saldo adalah terminal mandiri (layar sentuh atau komputer) yang ditempatkan di lokasi strategis di pondok (asrama, perpustakaan, kantor) agar santri dapat mengecek saldo kartu NFC kapan saja.

### Cara Menggunakan Kiosk

1. Santri **menempelkan kartu NFC** ke reader yang tersedia di kiosk.
2. Layar menampilkan secara otomatis:
   - **Nama santri**
   - **Saldo kartu saat ini**
   - **5 transaksi terakhir** (nama merchant, tanggal, jumlah)
3. Informasi tampil selama 30 detik, lalu layar kembali ke halaman awal.

### Hal yang Perlu Diketahui

- **Tidak perlu login** — kiosk langsung membaca kartu.
- **Tidak ada PIN di kiosk** — keamanan dilakukan melalui pembatasan akses jaringan (hanya terminal yang terdaftar yang dapat mengakses sistem).
- **Kiosk hanya untuk lihat saldo** — tidak ada transaksi transfer atau pembayaran dari kiosk.
- Jika kiosk tidak merespons kartu, coba tempelkan kartu lebih lama (tahan 2–3 detik). Jika masih gagal, kartu mungkin perlu diganti — hubungi admin pondok.

---

## 9. Notifikasi

🔔 Platform mengirimkan notifikasi melalui berbagai saluran agar Anda selalu mendapat informasi terkini.

### Jenis Notifikasi

| Saluran | Contoh Penggunaan |
|---------|------------------|
| **Push Notification (FCM)** | Transaksi masuk/keluar, tagihan baru, pengumuman |
| **WhatsApp** | OTP login, konfirmasi transaksi besar, tagihan jatuh tempo |
| **SMS** | OTP (cadangan jika WhatsApp tidak tersedia) |
| **Email** | Laporan bulanan, dokumen penting, reset password |

### Cara Mengaktifkan Notifikasi di Aplikasi

Saat pertama kali membuka aplikasi, sistem akan meminta izin notifikasi. Ketuk **"Izinkan"** agar notifikasi push dapat diterima.

Jika sebelumnya Anda menolak:
1. Buka **Pengaturan** di smartphone Anda.
2. Cari nama aplikasi (misal: "BMT Santri" atau sesuai nama BMT Anda).
3. Masuk ke **Notifikasi** → aktifkan semua izin notifikasi.

**Mengatur preferensi notifikasi di dalam aplikasi:**

Masuk ke **Profil → Pengaturan Notifikasi** untuk memilih jenis notifikasi yang ingin Anda terima (contoh: matikan notifikasi promosi, tapi tetap terima notifikasi transaksi).

### Kustomisasi Template (Admin BMT)

Admin BMT dapat mengubah isi pesan notifikasi melalui **Settings → Notifikasi → Template**. Tersedia placeholder seperti `{nama_nasabah}`, `{jumlah}`, `{nomor_rekening}` yang akan digantikan otomatis oleh sistem.

---

## 10. Fitur Autodebet

⚙️ Autodebet adalah fitur penarikan otomatis dari rekening nasabah pada tanggal yang sudah ditentukan, tanpa perlu nasabah datang ke kantor.

### Apa yang Di-autodebet?

Platform mendukung tiga jenis autodebet:

1. **Simpanan Wajib Bulanan** — Setiap bulan, sejumlah nominal tertentu ditarik otomatis dari rekening tabungan sukarela ke rekening simpanan wajib (sesuai produk yang disepakati saat buka rekening).

2. **Biaya Admin Rekening** — Biaya administrasi bulanan rekening ditarik otomatis sesuai tarif produk.

3. **Angsuran Pembiayaan** — Cicilan pembiayaan (kredit syariah) ditarik otomatis dari rekening nasabah pada tanggal jatuh tempo yang tercantum di jadwal angsuran.

### Kapan Autodebet Dijalankan?

Autodebet diproses oleh sistem pada malam hari (biasanya pukul 00.00–02.00 WIB) sesuai tanggal yang ditentukan. Pastikan saldo mencukupi sebelum tanggal jatuh tempo.

### Apa yang Terjadi Jika Saldo Tidak Cukup?

> **Penting:** Jika saldo tidak mencukupi, autodebet **tidak dibatalkan seluruhnya**. Sistem menerapkan **partial debit** (debit sebagian):

- Sistem akan mendebit **sebanyak saldo yang tersedia** dari rekening.
- Sisa yang belum terbayar dicatat sebagai **tunggakan** (piutang).
- Nasabah menerima notifikasi bahwa autodebet hanya berhasil sebagian beserta jumlah tunggakan yang tersisa.

Contoh: Angsuran Rp 500.000, saldo Rp 200.000 → Rp 200.000 didebet, tunggakan Rp 300.000 dicatat.

### Cara Membayar Tunggakan

1. **Via Aplikasi Nasabah** → menu **"Tagihan"** → **"Tunggakan"** → pilih tunggakan → bayar dari saldo rekening.
2. **Via Teller** — datang ke kantor, sampaikan ingin membayar tunggakan autodebet, teller akan memproses.
3. **Via Transfer** — transfer ke nomor rekening angsuran yang tertera di detail pembiayaan.

> **Tips:** Aktifkan notifikasi WhatsApp agar Anda mendapat pengingat H-3 sebelum tanggal autodebet. Ini memberi waktu untuk memastikan saldo cukup.

---

## 11. Keamanan & Privasi

🔒 Keamanan data nasabah dan santri adalah prioritas utama platform.

### Isolasi Data Antar Lembaga

Data setiap BMT dan pondok **sepenuhnya terpisah** dari lembaga lain. Staf BMT A tidak dapat melihat data nasabah BMT B, dan sebaliknya. Isolasi ini diterapkan di tingkat database — bukan hanya di tampilan aplikasi.

### OTP Berlaku Hanya 5 Menit

Kode OTP untuk login atau konfirmasi transaksi hanya valid selama 5 menit. Setelah itu, kode kedaluwarsa dan Anda perlu meminta kode baru. Ini mencegah penyalahgunaan kode OTP yang terlambat digunakan.

### Audit Log — Semua Tercatat

Setiap transaksi, perubahan data, dan aktivitas login tercatat di **audit log** dengan informasi:
- Waktu kejadian
- Pengguna yang melakukan
- Jenis aktivitas
- Data sebelum dan sesudah perubahan
- Perangkat dan lokasi (IP address)

Audit log tidak dapat dihapus dan hanya dapat dilihat oleh admin/manager berwenang.

### Lindungi Akun Anda

- **Jangan bagikan OTP** kepada siapapun, termasuk orang yang mengaku dari pihak BMT atau pondok.
- **Jangan bagikan password dan PIN** kepada siapapun.
- **Jangan gunakan jaringan WiFi publik** saat melakukan transaksi keuangan.
- **Aktifkan kunci layar** di smartphone Anda.
- **Logout** setelah selesai menggunakan aplikasi di perangkat bersama.
- **Segera hubungi BMT** jika Anda mencurigai akun Anda diakses orang lain.

### Hak Privasi Anda

Sesuai ketentuan privasi platform, data Anda:
- Hanya digunakan untuk keperluan operasional BMT dan pondok Anda.
- Tidak dijual atau dibagikan kepada pihak ketiga tanpa izin Anda.
- Disimpan dengan enkripsi di server yang aman.

---

## 12. Pertanyaan Umum (FAQ)

**Q: Apakah data saya aman?**

A: Ya. Data Anda disimpan di server dengan enkripsi dan hanya dapat diakses oleh staf BMT atau pondok Anda yang berwenang. Data antar lembaga sepenuhnya terpisah. Semua akses tercatat di audit log yang tidak dapat dimanipulasi.

---

**Q: Bagaimana jika saldo autodebet tidak cukup?**

A: Sistem akan mendebit sebanyak saldo yang tersedia (partial debit) dan mencatat sisanya sebagai tunggakan. Anda akan mendapat notifikasi dengan rincian tunggakan. Bayar tunggakan melalui Aplikasi Nasabah (menu Tagihan → Tunggakan) atau datang ke kantor teller.

---

**Q: Bagaimana cara bayar angsuran pembiayaan?**

A: Ada tiga cara: (1) Autodebet otomatis — angsuran ditarik otomatis dari rekening Anda pada tanggal jatuh tempo; (2) Transfer dari Aplikasi Nasabah ke nomor rekening angsuran; (3) Bayar tunai di teller kantor cabang.

---

**Q: Apakah bisa top-up kartu NFC dari aplikasi?**

A: Ya. Buka Aplikasi Nasabah → Kartu Saya → Top-up → pilih rekening sumber → masukkan jumlah → konfirmasi dengan PIN. Saldo kartu bertambah secara instan.

---

**Q: Bagaimana cara melihat mutasi rekening?**

A: Buka Aplikasi Nasabah → pilih rekening → ketuk "Mutasi". Anda dapat memfilter berdasarkan tanggal dan mengunduh mutasi dalam format PDF (maksimal 3 bulan terakhir).

---

**Q: Apa itu kolektibilitas?**

A: Kolektibilitas adalah tingkat kelancaran pembayaran angsuran pembiayaan, dikelompokkan menjadi 5 level sesuai standar OJK: (1) Lancar, (2) Dalam Perhatian Khusus, (3) Kurang Lancar, (4) Diragukan, (5) Macet. Semakin tinggi level, semakin besar risiko dan potensi dampak terhadap bunga/margin pembiayaan Anda.

---

**Q: Bagaimana cara reset PIN NFC?**

A: Hubungi admin BMT atau admin pondok. PIN NFC tidak dapat direset sendiri melalui aplikasi — memerlukan verifikasi identitas fisik di kantor. Setelah reset, Anda akan diminta membuat PIN baru saat transaksi NFC berikutnya.

---

**Q: Bisa transfer ke rekening bank lain (di luar platform)?**

A: Transfer ke rekening bank lain (BCA, BRI, Mandiri, dll.) saat ini belum tersedia langsung di Aplikasi Nasabah. Untuk transfer keluar, Anda dapat melakukannya melalui teller di kantor cabang BMT. Hubungi BMT Anda untuk informasi biaya dan ketentuan transfer antarbank.

---

**Q: Bagaimana kalau aplikasi error / tidak bisa login?**

A: Coba langkah berikut:
1. Pastikan koneksi internet Anda stabil.
2. Tutup aplikasi sepenuhnya lalu buka kembali.
3. Periksa apakah ada pembaruan aplikasi di Play Store / App Store.
4. Jika masih tidak bisa, coba logout dan login ulang.
5. Jika masih bermasalah, hubungi admin BMT atau pondok Anda — sertakan screenshot pesan error yang muncul.

---

**Q: Apakah ada biaya transaksi di aplikasi?**

A: Kebijakan biaya transaksi ditentukan oleh masing-masing BMT. Beberapa BMT menerapkan biaya admin bulanan yang sudah termasuk semua transaksi internal. Tanyakan ke BMT Anda untuk informasi biaya yang berlaku.

---

**Q: Bagaimana wali santri bisa memantau pengeluaran anak di kantin?**

A: Setiap transaksi NFC kartu santri dapat dilihat di Aplikasi Nasabah (Wali Santri) → Kartu Saya → Riwayat Transaksi Kartu. Anda dapat melihat di mana dan berapa santri berbelanja setiap hari.

---

## 13. Kontak & Dukungan

Untuk bantuan terkait penggunaan aplikasi, data rekening, atau masalah teknis — **hubungi admin BMT atau admin pondok Anda** secara langsung.

Setiap BMT dan pondok yang menggunakan platform ini memiliki tim dukungan sendiri yang siap membantu anggota dan wali santri mereka.

> **Catatan:** Pihak platform pusat tidak melayani pertanyaan langsung dari nasabah atau wali santri — semua pertanyaan disalurkan melalui admin BMT atau pondok Anda masing-masing.

**Informasi yang perlu Anda siapkan saat menghubungi admin:**

- Nama lengkap dan nomor rekening (untuk pertanyaan terkait keuangan).
- Nama santri dan nomor induk santri (untuk pertanyaan akademik).
- Screenshot pesan error (untuk masalah teknis).
- Waktu kejadian masalah.

---

*Panduan ini akan diperbarui secara berkala seiring pengembangan platform. Pastikan Anda selalu menggunakan versi aplikasi terbaru untuk mendapatkan fitur dan perbaikan keamanan terkini.*

---

**Platform Pesantren Terpadu** · CBS Syariah · ERP Pondok · E-commerce OPOP
