# Domain: Testing

## Strategi Testing

```
Unit Test        → Kalkulasi: murabahah, bagi hasil, autodebet partial, beasiswa
                   Settings resolver (3 level override)
                   Form state machine, jurnal balance
                   NFC: PIN, limit, idempotency
                   E-commerce: stok, harga, checkout

Integration Test → Handler + DB nyata (testcontainers-go)

E2E Test         → Developer setup BMT → pecahan → kontrak → cabang
                   Teller buka sesi (pecahan dari DB) → transaksi → tutup
                   Autodebet: partial + tunggakan → pelunasan saat saldo masuk
                   SPP: generate tagihan → beasiswa → autodebet partial → tunggakan
                   Pondok: santri → kelas → jadwal → absensi (3 metode) → nilai → raport
                   E-commerce: tambah produk → wali beli → bayar rekening BMT → pesanan
                   OPOP: pondok A order dari pondok B → B2B checkout

Widget Test      → Semua 7 Flutter app
```

## Test Wajib

```go
func TestSettings_TidakAdaHardcode_SelaluDariDB(t *testing.T)
func TestPecahanUang_DariDB_BukanKonstanta(t *testing.T)
func TestAutodebet_TanggalDariRekeningConfig_Benar(t *testing.T)
func TestAutodebet_SaldoKurang_PartialDebitDanTunggakan(t *testing.T)
func TestBeasiswa_TagihanSPP_NominalEfektifBenar(t *testing.T)
func TestBeasiswa_Pembiayaan_SaldoAngsuranDikurangi(t *testing.T)
func TestAbsensi_MetodeDariSettings_BukanKonstanta(t *testing.T)
func TestRaport_NilaiTertimbang_KomponenDariDB(t *testing.T)
func TestOPOP_B2BPesanan_LintasBMT_Berhasil(t *testing.T)
func TestEcommerce_BayarRekeningBMT_SaldoTerpotong(t *testing.T)
func TestJurnal_SemuaTransaksi_DoubleEntryBalance(t *testing.T)
func TestCrossTenant_QueryTanpaBMTID_Dilarang(t *testing.T)
```

## Stack Testing

```
Go  : go testing + testify + testcontainers-go
      - testcontainers untuk PostgreSQL + Redis nyata (bukan mock)
      - Tidak ada mock database — integration test pakai DB asli

Flutter: widget test + integration test
```
