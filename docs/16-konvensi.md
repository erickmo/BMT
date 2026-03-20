# 16 — Konvensi Koding & Testing

> **Terakhir diperbarui:** 20 Maret 2026

## Go — 7 Aturan Wajib

```go
// 1. UANG = int64, TIDAK PERNAH float
type Money int64
var biaya Money = 5_000   // ✅
var biaya float64 = 5000  // ❌

// 2. SETTINGS dari DB, bukan konstanta
jam := settings.Resolve(ctx, bmtID, cabangID, "operasional.jam_buka")  // ✅
const jamBuka = "08:00"  // ❌

// 3. FEATURE GATE wajib sebelum proses
if !featureGate.IsEnabled(ctx, bmtID, "CBS_DEPOSITO") {
    return ErrFiturTidakAktif
}

// 4. TRANSAKSI: DB tx + jurnal + usage_log
s.db.WithTx(ctx, func(tx pgx.Tx) error {
    // SELECT FOR UPDATE → validasi → update saldo
    // INSERT transaksi → journal.Post() → usageLog.Catat()
    // Dispatch event async
})

// 5. QUERY wajib scope tenant
WHERE bmt_id = @bmt_id AND cabang_id = @cabang_id  // ✅
WHERE id = @id  // ❌

// 6. UPDATE data via form, BUKAN langsung
formSvc.TerapkanForm(ctx, formID)  // ✅
db.Exec("UPDATE nasabah SET nama = $1...", nama, id)  // ❌

// 7. AUTODEBET GAGAL: partial + tunggakan
berhasil := min(rekening.Saldo, jadwal.NominalTarget)
sisa     := jadwal.NominalTarget - berhasil
if sisa > 0 { tunggakanRepo.Insert(ctx, tx, sisa) }

// 8. SDM & PAYROLL: selalu via module-vernon-hrm
hasil, _ := payroll.Hitung(ctx, payroll.Input{KontrakID: id, Periode: periode})  // ✅
db.Exec("SELECT gaji_pokok FROM sdm_kontrak WHERE id = $1", id)  // ❌ jangan kalkulasi sendiri
```

## Error Sentinel
```go
var (
    ErrSaldoTidakCukup           = errors.New("saldo tidak mencukupi")
    ErrRekeningBeku              = errors.New("rekening dalam status blokir")
    ErrRekeningTutup             = errors.New("rekening sudah ditutup")
    ErrTidakAdaSesiTeller        = errors.New("tidak ada sesi teller aktif")
    ErrSesiTellerSelisih         = errors.New("saldo fisik tidak sesuai, sesi ditolak")
    ErrFormTidakBisaDiubah       = errors.New("status form tidak mengizinkan perubahan")
    ErrApproverTidakBerwenang    = errors.New("role tidak berwenang menyetujui form ini")
    ErrKartuNFCTidakAktif        = errors.New("kartu NFC tidak aktif atau expired")
    ErrKartuNFCPINSalah          = errors.New("PIN kartu tidak sesuai")
    ErrFiturTidakAktif           = errors.New("fitur tidak tersedia di paket Anda")
    ErrJurnalTidakBalance        = errors.New("jurnal tidak balance: debit ≠ kredit")
    ErrSettingsNotFound          = errors.New("konfigurasi tidak ditemukan")
    ErrMidtransSignatureInvalid  = errors.New("signature Midtrans tidak valid")
    ErrDuplicateIdempotency      = errors.New("transaksi sudah diproses sebelumnya")
    ErrStokTidakCukup            = errors.New("stok produk tidak mencukupi")
    ErrIPKioskTidakDiizinkan     = errors.New("IP tidak terdaftar sebagai terminal kiosk")
    ErrListingTidakAktif         = errors.New("listing tidak aktif atau sudah expired")
)
```

## Standard Response
```go
type Response[T any] struct {
    Sukses bool      `json:"sukses"`
    Data   T         `json:"data,omitempty"`
    Error  *APIError `json:"error,omitempty"`
    Meta   *Meta     `json:"meta,omitempty"`
}
type APIError struct {
    Kode  string `json:"kode"`
    Pesan string `json:"pesan"`
}
```

## Flutter

```dart
// BLoC/Cubit per fitur
class SetoranCubit extends Cubit<SetoranState> {
  Future<void> prosesSetoran(SetoranRequest req) async {
    emit(const SetoranLoading());
    final result = await _useCase.call(req);
    result.fold(
      (f) => emit(SetoranGagal(f.pesan)),
      (t) => emit(SetoranBerhasil(t)),
    );
  }
}

// Model: gunakan freezed
@freezed
class Transaksi with _$Transaksi {
  const factory Transaksi({
    required String id,
    required int nominal,   // integer Rupiah
    required String jenis,
    required DateTime createdAt,
  }) = _Transaksi;
  factory Transaksi.fromJson(Map<String, dynamic> json) =>
      _$TransaksiFromJson(json);
}

// Tenant context auto-inject ke setiap request
class TenantInterceptor extends Interceptor {
  @override
  void onRequest(options, handler) {
    final session = getIt<SessionCubit>().state;
    options.headers['X-BMT-ID']    = session.bmtId;
    options.headers['X-Cabang-ID'] = session.cabangId;
    handler.next(options);
  }
}
```

## Database — Kolom Wajib
```sql
-- Tenant isolation (semua tabel operasional)
bmt_id    UUID NOT NULL REFERENCES bmt(id),
cabang_id UUID NOT NULL REFERENCES cabang(id),

-- Audit trail (semua tabel keuangan)
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
created_by UUID NOT NULL REFERENCES pengguna(id),
updated_by UUID NOT NULL REFERENCES pengguna(id),
deleted_at TIMESTAMPTZ,
is_voided  BOOLEAN NOT NULL DEFAULT FALSE,
voided_at  TIMESTAMPTZ, voided_by UUID, void_reason TEXT

-- Index wajib (prefix bmt_id)
CREATE INDEX idx_nasabah_bmt     ON nasabah(bmt_id, nomor_nasabah);
CREATE INDEX idx_rekening_nomor  ON rekening(bmt_id, nomor_rekening);
CREATE INDEX idx_transaksi_tgl   ON transaksi_rekening(bmt_id, cabang_id, created_at);
CREATE INDEX idx_jurnal_periode  ON jurnal(bmt_id, cabang_id, tanggal);
CREATE INDEX idx_audit_bmt_tgl   ON audit_log(bmt_id, created_at);
```

## Testing

```
Unit        → Kalkulasi, settings resolver, feature gate, form state machine, fraud rules
Integration → Handler + DB nyata (testcontainers-go)
E2E         → Setup BMT → tier/addon → nasabah → rekening → transaksi → laporan
              Autodebet partial → tunggakan → lunas
              OPOP → pesanan → komisi
              Listing → pendaftaran → approve → tampil di app
Widget      → Semua 7 Flutter app
Offline     → Teller offline → sync → conflict resolution
```

**Test wajib:**
```go
func TestFeatureGate_FiturTidakDiTier_Ditolak(t *testing.T)
func TestFeatureGate_AddOnAktif_Diizinkan(t *testing.T)
func TestSettings_TidakAdaHardcode_SelaluDariDB(t *testing.T)
func TestAutodebet_PartialDebitDanTunggakan(t *testing.T)
func TestJurnal_SemuaTransaksi_DoubleEntryBalance(t *testing.T)
func TestListing_PendaftaranSelfRegister_MenungguApprove(t *testing.T)
func TestListing_Premium_TampilDiAtas(t *testing.T)
func TestSaaSInvoice_GenerateBulanan_Benar(t *testing.T)
func TestOfflineSync_TransaksiTeller_BerhasilDiSync(t *testing.T)
func TestCrossTenant_QueryTanpaBMTID_Dilarang(t *testing.T)
```

**Coverage minimum:**
| Layer | Target |
|-------|--------|
| Domain logic & kalkulasi | 90% |
| Feature gate | 95% |
| Settings resolver | 95% |
| Tenant isolation | 95% |
| Service layer | 80% |
| Handler | 70% |
