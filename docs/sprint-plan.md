# Sprint Plan — Platform Pesantren Terpadu

> **Dibuat:** 20 Maret 2026
> **Basis:** Gap analysis dari state implementasi aktual vs dokumentasi
> **Durasi per sprint:** 2 minggu
> **Total:** 8 sprint (20 Mar — 10 Jul 2026)

---

## Prinsip Prioritas

1. **CBS dulu** — revenue utama BMT, harus jalan sebelum fitur lain
2. **Prinsip wajib ditegakkan** — feature gate, form workflow, audit log harus ada sebelum fitur baru
3. **Notifikasi menyertai fitur** — setiap fitur yang ship harus punya notifikasi
4. **Worker follow service** — worker diimplementasi bersamaan dengan service-nya

---

## Sprint 1 — CBS Core Service Layer
**20 Mar – 3 Apr 2026**

Fokus: Wire semua CBS handler ke service. Ini blockernya semua fitur lain.

### Backend (Go)
- [ ] `RekeningService.Tarik()` — validasi saldo, lock, transaksi, jurnal, usageLog
- [ ] `RekeningService.Transfer()` — antar rekening (in-branch + cross-branch)
- [ ] `NasabahService` — CRUD lengkap (CreateNasabah, GetNasabah, ListNasabah, SearchNasabah)
- [ ] `SesiTellerService` — BukaEsi, TutupSesi, validasi selisih kas
- [ ] Wire handler `/teller/*` → services (Setor, Tarik, Transfer, CariNasabah, CariRekening)
- [ ] Wire handler `/nasabah/*` → services (Profil, Mutasi, Detail rekening)
- [ ] `middleware/feature_gate.go` — `RequireFeature(kode)` cek paket BMT + add-on
- [ ] `middleware/audit_log.go` — catat setiap aksi ke `audit_log` table
- [ ] Error sentinel package — 17 sentinel vars standar

### Target akhir sprint
Teller bisa buka sesi, setor, tarik, transfer. Nasabah bisa lihat saldo + mutasi.

---

## Sprint 2 — Form Workflow + Auth Keamanan
**3 Apr – 17 Apr 2026**

Fokus: Tegakkan prinsip "tidak ada update langsung" dan keamanan dasar.

### Backend (Go)
- [ ] `FormService` — CreateForm, AjukanForm, SetujuiForm, TolakForm
- [ ] Auto-execute saat form DISETUJUI (CreateNasabah, BukaRekening, BlokirRekening, TutupRekening)
- [ ] Wire handler `/api/form/*` → FormService
- [ ] `OTPService` — GenerateOTP (6 digit, TTL 5 menit, Redis), ValidasiOTP, BlokirBrute-force
- [ ] Kirim OTP via SMS (Zenziva) dan Email (SMTP) — provider dari settings
- [ ] `SessionService` — CreateSession (device_id, IP), ListSesi, CabutSesi
- [ ] `SettingsService.Resolve()` — 3-level inheritance (platform → BMT → cabang) + cache Redis
- [ ] Handler `/auth/*` — login staf + nasabah + pondok dengan OTP 2FA

### Target akhir sprint
Nasabah/rekening baru hanya bisa dibuat via form + approval. Login wajib OTP jika 2FA aktif.

---

## Sprint 3 — Autodebet Lengkap + Workers CBS
**17 Apr – 1 Mei 2026**

Fokus: Autodebet harus production-ready (partial debit + tunggakan).

### Backend (Go)
- [ ] `AutodebetService.ExecuteBulanan()` — debit sejumlah saldo tersedia, INSERT tunggakan jika kurang
- [ ] `AutodebetService.GenerateJadwal()` — generate jadwal autodebet bulan depan (tgl 25)
- [ ] `TunggakanService` — list tunggakan, bayar tunggakan, riwayat
- [ ] Worker `HandleAutodebetBulanan` — wired ke service bulanan
- [ ] Worker `HandleGenerateJadwalAutodebet` — tgl 25 tiap bulan
- [ ] Worker `HandleUpdateKolektibilitas` — 00:05 harian, klasifikasi OJK 5 level
- [ ] Worker `HandleReminderAngsuran` — H-N hari sebelum jatuh tempo (N dari settings)
- [ ] Worker `HandleDistribusiBagiHasil` — akhir bulan, hitung dari realisasi pendapatan

### Target akhir sprint
Autodebet berjalan otomatis, partial debit benar, tunggakan tercatat. Kolektibilitas update harian.

---

## Sprint 4 — Notifikasi + Midtrans
**1 Mei – 15 Mei 2026**

Fokus: Setiap event penting harus ada notifikasinya.

### Backend (Go)
- [ ] `NotifikasiService.Kirim()` — render template + queue ke `notifikasi_antrian`
- [ ] Worker `HandleKirimNotifikasi` — proses antrian, delivery per channel
- [ ] FCM delivery — push notification via Firebase (konfigurasi dari settings)
- [ ] WhatsApp delivery — Fonnte/Wablas (provider dari settings per BMT)
- [ ] SMS delivery — Zenziva/Twilio (provider dari settings)
- [ ] Email delivery — SMTP (konfigurasi dari settings)
- [ ] Retry with backoff (3x, exponential) + update status delivery
- [ ] Event trigger notifikasi: Setor, Tarik, Transfer, Autodebet gagal, Angsuran H-3/H-1
- [ ] Midtrans webhook `/webhook/midtrans` — verifikasi SHA512, post transaksi + jurnal + usageLog
- [ ] Worker `HandleCekMidtransPending` — poll PENDING > 30 menit setiap 15 menit

### Target akhir sprint
Nasabah dapat notifikasi FCM + WA setiap transaksi. Midtrans webhook production-ready.

---

## Sprint 5 — Pembiayaan + Akuntansi Bridge
**15 Mei – 29 Mei 2026**

Fokus: Pembiayaan end-to-end + jurnal otomatis semua transaksi.

### Backend (Go)
- [ ] `PembiayaanService` — state machine lengkap (PENGAJUAN→ANALISIS→KOMITE→AKAD→PENCAIRAN→AKTIF→LUNAS)
- [ ] Approval komite via form workflow
- [ ] `PembiayaanService.BayarAngsuran()` — update saldo, jurnal, reminder
- [ ] `AkuntansiService` — bridge handler `/finance/*` ke `module-vernon-accounting`
- [ ] Default COA seed per BMT baru (akun 1xx-5xx)
- [ ] Laporan: Neraca, SHU, Arus Kas, Kolektibilitas, Bagi Hasil — generate dari module-vernon-accounting
- [ ] Worker `HandleHitungZakat` — akhir tahun, basis nisab dari settings
- [ ] Notifikasi: Pembiayaan disetujui, angsuran jatuh tempo, lunas

### Target akhir sprint
Pembiayaan bisa diajukan, disetujui komite, dicairkan, dan dibayar. Laporan keuangan bisa di-generate.

---

## Sprint 6 — Pondok Akademik + Operasional
**29 Mei – 12 Jun 2026**

Fokus: Pondok MVP — administrasi santri, akademik, SPP.

### Backend (Go)
- [ ] `SantriService` — CRUD santri, link nasabah, kamar
- [ ] `AbsensiService` — catat absensi (MANUAL/NFC/BIOMETRIK), rekap harian
- [ ] `NilaiService` — input nilai per komponen, hitung nilai akhir (bobot dari settings)
- [ ] `RaportService` — generate raport DRAFT → FINAL → DITERBITKAN
- [ ] `TagihanSPPService` — generate tagihan, hitung beasiswa, bayar via rekening BMT
- [ ] `PPDBService` — state machine (MENDAFTAR→SELEKSI→DITERIMA/DITOLAK), payment Midtrans
- [ ] Wire handler `/pondok/*` → services (administrasi, akademik, keuangan pondok)
- [ ] Worker `HandleGenerateTagihanSPP` — tgl 25 tiap bulan
- [ ] Worker `HandleReminderSPP` — H-N sebelum jatuh tempo
- [ ] Worker `HandleGenerateRaport` — akhir semester (jadwal dari settings)
- [ ] Worker `HandlePayrollBulanan` — via module-vernon-hrm

### Target akhir sprint
Admin pondok bisa kelola santri, input nilai, generate raport, tagih SPP.

---

## Sprint 7 — OPOP Ecommerce + Storage + PDF
**12 Jun – 26 Jun 2026**

Fokus: OPOP marketplace siap transaksi + dokumen bisa di-generate.

### Backend (Go)
- [ ] `PesananService` — state machine (MENUNGGU→DIBAYAR→DIPROSES→DIKIRIM→SELESAI/DIBATALKAN)
- [ ] Pembayaran pesanan via Midtrans + Rekening BMT + NFC
- [ ] `MerchantService` — onboarding, terminal kiosk, settlement komisi
- [ ] `NfcService` — issue kartu, PIN management, tap transaction, validasi IP kiosk
- [ ] MinIO integration `pkg/storage/` — upload file (foto produk, portfolio, dokumen)
- [ ] PDF generation `pkg/pdfgen/` — slip setoran/penarikan, kartu anggota (A6), buku tabungan (A5), akad pembiayaan, kartu angsuran
- [ ] Worker `HandleAnalyticsHarian` — snapshot GMV, pesanan, pembeli (23:00)
- [ ] Worker `HandleHitungKomisiOPOP` — akhir bulan, % dari kontrak per toko
- [ ] Notifikasi: status pesanan update (FCM nasabah + FCM merchant)

### Target akhir sprint
Wali santri bisa beli produk pondok, toko bisa kelola pesanan, slip transaksi bisa dicetak.

---

## Sprint 8 — SaaS Platform + Fraud + DAPODIK
**26 Jun – 10 Jul 2026**

Fokus: Platform production-ready — billing, keamanan, integrasi eksternal.

### Backend (Go)
- [ ] Feature gate enforcement — `RequireFeature()` dipasang di semua route sesuai docs/14-api.md
- [ ] `MonetisasiService` — usage metering per BMT (transaksi, nasabah, storage)
- [ ] `BillingService` — generate invoice bulanan, kirim email ke ADMIN_BMT
- [ ] Worker `HandleHitungUsageBulanan` — tgl 1 tiap bulan
- [ ] Worker `HandleGenerateSaaSInvoice` — tgl dari settings platform
- [ ] Worker `HandleCekSaaSExpiry` — 08:00, suspend BMT expired
- [ ] Worker `HandleCekListingExpiry` — 08:00 + reminder H-7
- [ ] `FraudService` — eval rule JSONB per transaksi, aksi ALERT/BLOCK/REQUIRE_APPROVAL
- [ ] Worker `HandleFraudDetection` — real-time per event transaksi
- [ ] DAPODIK HTTP client — pull santri + guru dari API DAPODIK (config NPSN dari settings)
- [ ] EMIS HTTP client — pull data dari API EMIS (config NSM dari settings)
- [ ] Worker `HandleSinkronDAPODIK` + `HandleSinkronEMIS` — jadwal dari settings
- [ ] Worker `HandleBackupDatabase` — pg_dump → gzip → MinIO 02:00 WIB
- [ ] Dashboard WebSocket — transaksi real-time, fraud alert, sesi teller aktif

### Target akhir sprint
Platform bisa charge BMT, fraud ter-detect, data pondok tersinkron ke DAPODIK/EMIS.

---

## Ringkasan Timeline

```
Mar  ████████ Sprint 1: CBS Core Service Layer
Apr  ████████ Sprint 2: Form Workflow + Auth
     ████████ Sprint 3: Autodebet + Workers CBS
Mei  ████████ Sprint 4: Notifikasi + Midtrans
     ████████ Sprint 5: Pembiayaan + Akuntansi
Jun  ████████ Sprint 6: Pondok Akademik + Ops
     ████████ Sprint 7: OPOP + Storage + PDF
Jul  ████████ Sprint 8: SaaS + Fraud + DAPODIK
```

## Total Backlog per Kategori

| Kategori | Item |
|---|---|
| Service layer baru | 28 service methods |
| Worker handlers | 31 workers |
| Handler wiring | 19 handler groups |
| Infrastructure (storage, pdf, notif delivery) | 8 packages |
| Middleware baru | 3 (feature_gate, audit_log, session) |
| Integrasi eksternal | 4 (FCM, WA, DAPODIK, EMIS) |
