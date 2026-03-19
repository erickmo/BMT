# Domain: Konvensi Koding & Prinsip Utama

## 6 Prinsip Utama (Wajib Dipatuhi)

```
1. TIDAK ADA HARDCODE — semua nilai konfigurasi dari settings DB
2. TIDAK ADA UPDATE LANGSUNG — data nasabah/rekening via form + approval
3. SETIAP TRANSAKSI KEUANGAN = DB transaction + jurnal (module-vernon-accounting) + usage_log
4. SETIAP QUERY = di-scope bmt_id + cabang_id (tenant isolation)
5. UANG = int64 (Money), TIDAK PERNAH float
6. AUTODEBET GAGAL = partial debit + INSERT tunggakan (bukan skip atau error)
```

---

## Pola Go yang Benar

```go
// ✅ Settings dari DB — bukan konstanta
jam := settings.Resolve(ctx, bmtID, cabangID, "operasional.jam_buka")
metode := settings.ResolveJSON(ctx, bmtID, cabangID, "pondok.absensi_metode")

// ✅ Uang sebagai int64
type Money int64
var biaya Money = 5_000  // ✅
var biaya float64 = 5000 // ❌

// ✅ Transaksi keuangan: DB tx + jurnal + usage_log
s.db.WithTx(ctx, func(tx pgx.Tx) error {
    // 1. SELECT FOR UPDATE
    // 2. Validasi domain
    // 3. Update saldo
    // 4. INSERT transaksi_rekening
    // 5. journal.Post(...)
    // 6. usageLog.Catat(...)
    return nil
})

// ✅ Autodebet partial
berhasil := min(rekening.Saldo, jadwal.NominalTarget)
sisa     := jadwal.NominalTarget - berhasil
if sisa > 0 { tunggakanRepo.Insert(...) }

// ✅ Pecahan uang dari DB
pecahans, _ := pecahanRepo.GetAktif(ctx)  // BUKAN: []Pecahan{{100,"LOGAM"},...}

// ✅ Tenant isolation — semua query
WHERE bmt_id = @bmt_id AND cabang_id = @cabang_id
```

---

## Error Sentinel

```go
var (
    ErrSaldoTidakCukup           = errors.New("saldo tidak mencukupi")
    ErrRekeningBeku              = errors.New("rekening dalam status blokir")
    ErrTidakAdaSesiTeller        = errors.New("tidak ada sesi teller aktif")
    ErrSesiTellerSelisih         = errors.New("saldo fisik tidak sesuai, sesi ditolak")
    ErrFormTidakBisaDiubah       = errors.New("status form tidak mengizinkan perubahan")
    ErrApproverTidakBerwenang    = errors.New("role tidak berwenang menyetujui form ini")
    ErrKartuNFCTidakAktif        = errors.New("kartu NFC tidak aktif atau expired")
    ErrKartuNFCPINSalah          = errors.New("PIN kartu tidak sesuai")
    ErrFiturTidakAktif           = errors.New("fitur tidak diaktifkan di kontrak BMT")
    ErrJurnalTidakBalance        = errors.New("jurnal tidak balance: debit ≠ kredit")
    ErrSettingsNotFound          = errors.New("konfigurasi tidak ditemukan di settings")
    ErrMidtransSignatureInvalid  = errors.New("signature Midtrans tidak valid")
    ErrStokTidakCukup            = errors.New("stok produk tidak mencukupi")
    ErrIPKioskTidakDiizinkan     = errors.New("IP tidak terdaftar sebagai terminal kiosk")
)
```

---

## Checklist Syariah (Semua Domain)

### CBS
- [ ] Tidak ada riba — margin/nisbah/ujrah transparan & disepakati sebelum akad
- [ ] Ta'zir 100% masuk akun 211 — bukan pendapatan
- [ ] Bagi hasil dari realisasi pendapatan, bukan % nominal pokok
- [ ] Autodebet angsuran partial tetap menghasilkan jurnal syariah yang benar
- [ ] Biaya admin rekening: akad jelas saat buka rekening

### E-commerce
- [ ] Harga produk OPOP: transparan, tidak ada penipuan (gharar)
- [ ] Komisi platform dari e-commerce: akad jelas (ujrah/wakalah)

### Donasi & Wakaf
- [ ] Dana sosial tidak bercampur dengan operasional BMT
- [ ] Wakaf produktif: hasil usaha dibagikan ke mauquf alaih, bukan ke BMT
- [ ] Infaq/shadaqah: penyaluran tercatat dengan mustahiq

### SDM & Sosial
- [ ] Gaji pegawai: tidak ada unsur riba — gaji flat, bukan % keuntungan
- [ ] Denda keterlambatan buku perpustakaan: masuk dana sosial (akun 611)

---

## Glosarium

| Istilah | Definisi |
|---------|----------|
| **2FA** | Two-Factor Authentication — verifikasi dua langkah via OTP |
| **Alumni** | Santri yang sudah lulus dari pondok pesantren |
| **Anti-Fraud** | Sistem deteksi transaksi mencurigakan berbasis rules engine |
| **Autodebet** | Debit otomatis terjadwal; tanggal diset per rekening via management BMT |
| **Beasiswa** | Potongan biaya SPP/pembiayaan santri, ditetapkan admin pondok |
| **BMT** | Baitul Maal wa Tamwil — koperasi simpan pinjam berbasis syariah |
| **CBS** | Core Banking System — sistem inti perbankan BMT |
| **DAPODIK** | Data Pokok Pendidikan — database siswa nasional Kemendikbud |
| **EMIS** | Education Management Information System — database Kemenag |
| **ERP Pondok** | Enterprise Resource Planning untuk administrasi pondok pesantren |
| **FCM** | Firebase Cloud Messaging — layanan push notification Google |
| **Finance** | Role staf yang mengelola jurnal manual & biaya operasional |
| **Form Pengajuan** | Mekanisme wajib untuk semua perubahan data nasabah & rekening |
| **GMV** | Gross Merchandise Value — total nilai transaksi OPOP bruto |
| **Hardcode** | Nilai yang tertanam di kode — **dilarang**, semua harus dari settings DB |
| **Infaq** | Pengeluaran sukarela di jalan Allah (termasuk dana ta'zir) |
| **Jenis Rekening** | Tipe rekening dengan aturan & tarif sendiri (CRUD management BMT) |
| **Kartu NFC** | Kartu fisik santri untuk transaksi di merchant pondok |
| **Kiosk** | Terminal cek saldo NFC tanpa login |
| **Komponen Nilai** | Komponen penilaian (UH, UTS, UAS) dengan bobot % |
| **Modul Vernon** | `module-vernon-accounting` — mesin double-entry akuntansi internal |
| **Nasabah** | Anggota BMT yang memiliki rekening |
| **Nazhir** | Pengelola aset wakaf |
| **NPSN** | Nomor Pokok Sekolah Nasional — ID unik sekolah di DAPODIK |
| **NSM** | Nomor Statistik Madrasah — ID unik madrasah di EMIS |
| **OPOP** | One Pondok One Product — marketplace produk UMKM antar pondok |
| **Partial Debit** | Debit sebesar saldo tersedia saat autodebet gagal, sisa jadi tunggakan |
| **Payroll** | Sistem penggajian — transfer gaji otomatis ke rekening BMT pegawai |
| **Pecahan Uang** | Data redenominasi Rupiah di DB — bisa diupdate tanpa deploy ulang |
| **PPDB** | Penerimaan Peserta Didik Baru — pendaftaran santri baru online |
| **RPP** | Rencana Pelaksanaan Pembelajaran |
| **Santri** | Pelajar aktif pondok pesantren |
| **Sesi Teller** | Periode kerja teller satu hari dengan pembukuan kas berbasis redenominasi |
| **Settings Engine** | Sistem 3-level (platform→BMT→cabang) resolusi konfigurasi dari DB |
| **SHU** | Sisa Hasil Usaha — "laba" koperasi yang dibagikan ke nasabah |
| **SPP** | Sumbangan Pembinaan Pendidikan — iuran bulanan santri |
| **Tahfidz** | Hafalan Al-Quran — dinilai tersendiri dalam akademik pondok |
| **Ta'zir** | Denda keterlambatan — 100% masuk dana sosial, bukan pendapatan |
| **Tenant** | Satu BMT beserta seluruh cabangnya |
| **Toko OPOP** | Toko pondok di marketplace OPOP lintas pondok |
| **Tunggakan** | Sisa kewajiban autodebet yang belum terbayar akibat saldo kurang |
| **Usage Log** | Catatan transaksi yang dikenai biaya admin platform |
| **Wakaf** | Aset yang dibekukan manfaatnya untuk kepentingan umum |
| **Wakaf Produktif** | Aset wakaf yang dikelola BMT untuk usaha produktif |
| **Wakif** | Pemberi wakaf |
| **White-label** | Custom branding app per pondok/BMT |
