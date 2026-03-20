# 05 — CBS (Core Banking System)

> **Terakhir diperbarui:** 20 Maret 2026

## Feature Gate
Semua fitur CBS wajib dicek via `RequireFeature(kode)`:
- `CBS_DASAR` — nasabah, rekening sukarela, transaksi tunai (termasuk di tier FREE)
- `CBS_PEMBIAYAAN` — add-on
- `CBS_DEPOSITO` — add-on
- `CBS_AUTODEBET` — termasuk di BASIC+

## Jenis Rekening (CRUD Management BMT)
```sql
CREATE TABLE jenis_rekening (
    id, bmt_id, kode, nama,
    tipe_dasar VARCHAR(30),  -- SIMPANAN_POKOK|SIMPANAN_WAJIB|SIMPANAN_SUKARELA|DEPOSITO|TABUNGAN_KHUSUS
    akad VARCHAR(30),        -- WADIAH|WADIAH_YAD_DHAMANAH|MUDHARABAH
    setoran_awal_min BIGINT, setoran_min BIGINT,
    bisa_ditarik BOOLEAN,
    nisbah_nasabah SMALLINT, jangka_hari SMALLINT,
    biaya_admin_bulanan BIGINT,  -- per bulan, dikonfigurasi BMT
    bisa_nfc BOOLEAN,
    bisa_autodebet BOOLEAN,
    biaya_admin_buka BIGINT,
    is_aktif BOOLEAN, urutan_tampil SMALLINT,
    UNIQUE (bmt_id, kode)
);
```

## Rekening
- Nomor: `{KODE_BMT}-{KODE_CAB}-{KODE_JENIS}-{SEQ:08d}`
- `saldo BIGINT CHECK (saldo >= 0)`
- Status: `AKTIF | BEKU | TUTUP`
- `biaya_admin_bulanan` di-snapshot saat buka rekening
- Perubahan: wajib via form (→ docs/17-form-workflow.md)

## Pecahan Uang (Data, Bukan Konstanta)
```sql
CREATE TABLE pecahan_uang (
    id, nominal BIGINT, jenis VARCHAR(10),  -- LOGAM|KERTAS
    label VARCHAR(30), is_aktif BOOLEAN,
    urutan SMALLINT, berlaku_sejak DATE, ditarik_pada DATE
);
```
Teller app ambil dari DB saat buka sesi.

## Sesi Teller
- 1 sesi aktif per teller per hari
- Buka: input redenominasi per pecahan (dari DB) → `saldo_awal`
- Tutup: `saldo_akhir_fisik` harus = `saldo_akhir_sistem`
- Selisih ≠ 0 → status `DITOLAK`, teller wajib hitung ulang
- Toleransi selisih dari `settings: "sesi_teller.toleransi_selisih"`

## Autodebet

**Tanggal diset per rekening** di `rekening_autodebet_config.tanggal_debet` (1-28).

| Jenis | Sumber |
|-------|--------|
| `SIMPANAN_WAJIB` | Rekening sukarela |
| `BIAYA_ADMIN_REKENING` | Rekening itu sendiri |
| `ANGSURAN_PEMBIAYAAN` | Rekening simpanan |
| `SPP_PONDOK` | Rekening simpanan |

**Gagal:** partial debit + INSERT `tunggakan_autodebet` + notifikasi email.

## Pembiayaan (Feature: `CBS_PEMBIAYAAN`)
Akad: `MURABAHAH|MUSYARAKAH|MUDHARABAH|IJARAH|QARDH|RAHN`
Alur: `PENGAJUAN → ANALISIS → KOMITE → AKAD → PENCAIRAN → AKTIF → LUNAS`
Kolektibilitas OJK: 1=Lancar, 2=DPK(1-90hr), 3=KL(91-120), 4=Diragukan(121-180), 5=Macet(>180)

**Beasiswa** (ditetapkan `ADMIN_PONDOK`):
```sql
ada_beasiswa BOOLEAN,
beasiswa_persen NUMERIC(5,2),
beasiswa_nominal BIGINT,
beasiswa_sumber VARCHAR(100),
beasiswa_ditetapkan_oleh UUID REFERENCES pengguna_pondok(id)
```

## Transaksi Rekening
Jenis: `SETOR_TUNAI|SETOR_ONLINE|TARIK_TUNAI|TARIK_ONLINE|DEBIT_NFC|BAGI_HASIL|BONUS_WADIAH|KOREKSI_KREDIT|KOREKSI_DEBIT|TRANSFER_MASUK|TRANSFER_KELUAR`

Sumber: `TELLER|NASABAH_APP|NFC|SISTEM|FINANCE`

**Wajib setiap transaksi:**
1. `SELECT ... FOR UPDATE` (lock row)
2. Validasi domain
3. Update saldo rekening
4. INSERT transaksi_rekening (immutable)
5. `journal.Post(...)` — module-vernon-accounting
6. `usageLog.Catat(...)` — billing platform
7. Dispatch event async (notifikasi)

## NFC Santri (Feature: `CBS_DASAR`)
```sql
CREATE TABLE kartu_nfc (
    nasabah_id UUID UNIQUE,     -- 1 nasabah = 1 kartu
    rekening_id UUID,           -- bisa_nfc = true
    uid_nfc VARCHAR(50) UNIQUE,
    pin_hash VARCHAR(255),
    status VARCHAR(20),         -- AKTIF|BEKU|HILANG|EXPIRED
    limit_per_transaksi BIGINT,
    limit_harian BIGINT,
    tanggal_expired DATE
);
```
Tap → baca UID → input PIN → API validasi → debit rekening.

## Akuntansi — Chart of Accounts Default
```
1xx ASET      : 101 Kas, 102 Kas Bank, 111-114 Piutang/Pembiayaan
2xx KEWAJIBAN : 201 Simpanan Pokok/Wajib, 202 Sukarela, 203 Deposito, 211 Dana Sosial
3xx EKUITAS   : 301-302 Simpanan, 303 Cadangan, 304 SHU Ditahan
4xx PENDAPATAN: 401 Margin, 402-403 Bagi Hasil, 404 Ujrah, 405 Admin, 406 NFC, 407 OPOP
5xx BEBAN     : 501 Bagi Hasil, 502 Operasional, 503 Penyisihan, 504 Gaji, 505 Utilitas
```
