# 02 — Stack & Aplikasi

> **Terakhir diperbarui:** 20 Maret 2026

## Backend
```
API              : Go 1.23 (net/http + Chi router)
Database         : PostgreSQL 17
Cache            : Redis 7
Queue/Worker     : asynq (Redis-backed)
Auth             : JWT — access token 15m + refresh token 7d
Storage          : MinIO (self-hosted S3-compatible)
PDF              : chromedp (Go → headless Chrome)
Migration        : golang-migrate
Query            : sqlc (type-safe SQL → Go, BUKAN ORM)
Payment          : Midtrans (Snap + Core API)
Email            : SMTP — dikonfigurasi per BMT di settings
Biometrik        : Fingerprint SDK via REST
Search Listing   : PostgreSQL full-text search (fase 1), Elasticsearch (fase 2)
```

### Kenapa sqlc bukan ORM?
Sistem perbankan butuh kontrol penuh atas query SQL. sqlc menghasilkan Go type-safe dari `.sql` — tidak ada N+1, mudah diaudit.

## module-vernon-accounting (Internal Go)

Mesin akuntansi double-entry. Semua transaksi keuangan **wajib** memanggil modul ini.

```
module-vernon-accounting/
├── account/   # Chart of accounts, tipe akun (ASET/KEWAJIBAN/EKUITAS/PENDAPATAN/BEBAN)
├── journal/   # Double-entry engine — validasi Σdebit=Σkredit sebelum persist
├── ledger/    # Buku besar per akun
├── period/    # Periode akuntansi, closing bulanan/tahunan
└── report/    # neraca/, laba_rugi/, arus_kas/
```

```go
import "bmt-saas/module-vernon-accounting/journal"

journal.Post(ctx, Journal{
    BMTID: bmtID, CabangID: cabangID,
    Tanggal: time.Now(), Keterangan: "Setoran tunai",
    Referensi: transaksiID.String(),
    Entries: []Entry{
        {KodeAkun: "101", Posisi: DEBIT,  Nominal: nominal},
        {KodeAkun: "202", Posisi: KREDIT, Nominal: nominal},
    },
}) // error jika Σdebit ≠ Σkredit
```

**Aturan:**
- Jurnal tidak pernah dihapus — koreksi via reversal entry
- Setiap `journal.Post()` wajib berada di dalam DB transaction yang sama dengan transaksi keuangannya
- Kode akun dapat dikustomisasi per BMT via settings

---

## module-vernon-hrm (Internal Go)

Mesin SDM (Human Resource Management). Semua proses kepegawaian pondok dan BMT **wajib** memanggil modul ini.

```
module-vernon-hrm/
├── employee/      # Data karyawan/pengajar — canonical record SDM
├── contract/      # Kontrak kerja, tipe kontrak, periode
├── attendance/    # Rekap absensi, jam kerja, keterlambatan
├── leave/         # Pengajuan & approval cuti/izin
├── payroll/       # Kalkulasi gaji, tunjangan, potongan
├── payslip/       # Generate slip gaji (data + PDF trigger)
├── performance/   # KPI, penilaian kinerja periodik
└── README.md
```

### Prinsip Modul HRM

- **Canonical employee record** — satu sumber kebenaran data pegawai; `pondok_pengajar` dan `pondok_karyawan` mereferensikan `hrm.Employee`, bukan sebaliknya
- **Payroll pipeline** — kalkulasi deterministik: `GajiPokok + Tunjangan - PotonganAbsen - PotonganLain = GajiBersih`
- **Terintegrasi dengan module-vernon-accounting** — setiap run payroll otomatis posting jurnal via `journal.Post()`
- **Terintegrasi dengan absensi** — `attendance.Rekap(employeeID, periode)` menghasilkan `HariAbsen` yang digunakan payroll

```go
import (
    "bmt-saas/module-vernon-hrm/payroll"
    "bmt-saas/module-vernon-hrm/payslip"
    "bmt-saas/module-vernon-hrm/attendance"
)

// Kalkulasi gaji satu karyawan
hasil, err := payroll.Hitung(ctx, payroll.Input{
    KontrakID: kontrakID,
    Periode:   "2025-01",
    RekapAbsensi: attendance.Rekap(ctx, employeeID, "2025-01"),
})
// hasil.GajiBersih, hasil.Detail (tunjangan + potongan itemized)

// Generate slip gaji
slip, err := payslip.Generate(ctx, payslip.Input{
    HasilPayroll: hasil,
    TemplateBMT:  templateURL, // dari MinIO
})
// slip.FileURL → path MinIO PDF siap download

// Posting jurnal otomatis
journal.Post(ctx, slip.JurnalEntries...)
// Debit: 504 (Beban Gaji), Kredit: 202 (Hutang Gaji) → lalu Kredit: 101 (Kas) saat transfer
```

### Struktur Data Utama HRM

```go
// module-vernon-hrm/employee/employee.go
type Employee struct {
    ID          uuid.UUID
    BMTID       uuid.UUID
    CabangID    uuid.UUID
    NIP         string
    Nama        string
    Jabatan     string
    Tipe        EmployeeTipe  // PENGAJAR | KARYAWAN | STAF_BMT
    NasabahID   *uuid.UUID    // nullable — terhubung ke nasabah BMT
    StatusAktif bool
}

// module-vernon-hrm/contract/contract.go
type Contract struct {
    ID                    uuid.UUID
    EmployeeID            uuid.UUID
    TipeKontrak           string    // TETAP|KONTRAK|HONORER|PARUH_WAKTU
    GajiPokok             money.Money
    Tunjangan             []TunjanganItem
    PotonganPerHariAbsen  money.Money
    RekeningGajiID        uuid.UUID  // rekening BMT tujuan transfer gaji
    TanggalMulai          time.Time
    TanggalSelesai        *time.Time
}

// module-vernon-hrm/attendance/attendance.go
type RingkasanAbsensi struct {
    EmployeeID  uuid.UUID
    Periode     string
    TotalHari   int
    HadiR       int
    Sakit        int
    Izin         int
    Alfa         int
    Terlambat    int
    JamKerjaTotal time.Duration
}

// module-vernon-hrm/leave/leave.go
type PengajuanCuti struct {
    ID          uuid.UUID
    EmployeeID  uuid.UUID
    JenisCuti   string  // TAHUNAN|SAKIT|MELAHIRKAN|PENTING|TANPA_GAJI
    TanggalMulai time.Time
    TanggalSelesai time.Time
    Status      string  // MENUNGGU|DISETUJUI|DITOLAK
    Pengganti   *uuid.UUID
}
```

### Integrasi dengan Domain Pondok

```
pondok_pengajar.hrm_employee_id → module-vernon-hrm/employee (canonical)
pondok_karyawan.hrm_employee_id → module-vernon-hrm/employee (canonical)
pengguna (staf BMT).hrm_employee_id → module-vernon-hrm/employee (jika applicable)

pondok_absensi → feed ke module-vernon-hrm/attendance (sinkron harian)
surat_izin (disetujui) → create module-vernon-hrm/leave.PengajuanCuti otomatis
```

### Feature Gate
`SDM_PAYROLL` — add-on yang diperlukan untuk mengakses payroll pipeline penuh.
Fitur attendance & leave tersedia di tier BASIC ke atas.

## 7 Aplikasi Flutter

| App | Dir | Platform | Pengguna |
|-----|-----|----------|---------|
| Nasabah | `app/nasabah/` | Android, iOS | Nasabah BMT (e-banking + santri + listing) |
| Management | `app/management/` | Web, Desktop, Mobile | Staf & management BMT |
| Developer | `app/developer/` | Web, Desktop, Mobile | Developer platform |
| Teller | `app/teller/` | Desktop | Teller BMT |
| Merchant | `app/merchant/` | Android, iOS | Kasir + owner merchant NFC pondok |
| Cek Saldo | `app/ceksaldo/` | Android / Kiosk | Santri cek saldo NFC |
| Pondok | `app/pondok/` | Web, Mobile | Admin pondok |

### `app/developer` — Portal Developer
Akses via `Developer-Token`. Platform-wide management:
- Dashboard metrics (BMT aktif, MRR, transaksi/hari)
- **CRUD BMT** — tambah BMT baru, set paket tier, fitur aktif
- **Kelola paket tier** — harga, fitur yang termasuk
- **Kelola add-on** — fitur à la carte, harga
- **Kelola listing** — approve/reject pendaftar stakeholder
- **Kelola kategori listing** — tambah/edit kategori
- Pecahan uang Rupiah, tarif transaksi, kontrak BMT
- Platform settings, health check, maintenance mode

### `app/nasabah` — Dua Mode
**Mode Nasabah Biasa:** e-banking (saldo, setor, riwayat, pembiayaan)

**Mode Wali Santri:** semua fitur nasabah + profil santri, saldo NFC, raport, SPP, belanja OPOP, **direktori listing stakeholder** (cari guru les, antar jemput, dll. di sekitar pondok)

### `app/management`
Role: `ADMIN_BMT`, `MANAJER_BMT`, `MANAJER_CABANG`, `AO`, `KOMITE`, `FINANCE`, `AUDITOR`
Fitur: dashboard, laporan, form approval, nasabah (view only), rekening, pembiayaan, settings BMT

### `app/teller`
Semua tombol transaksi disabled tanpa sesi aktif.
Buka/tutup sesi: redenominasi dari DB (bukan hardcode).

### `app/pondok`
Role: `ADMIN_PONDOK`, `OPERATOR_PONDOK`, `BENDAHARA_PONDOK` + role spesialis
Fitur: administrasi santri/guru, akademik, jadwal, absensi, penilaian, raport, SPP, toko OPOP

### `app/merchant`
- Mode Kasir: input nominal → tap NFC → PIN → konfirmasi → struk
- Mode Owner: dashboard penjualan, riwayat, laporan, export CSV
