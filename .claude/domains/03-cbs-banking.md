# Domain: CBS — Core Banking System

## Jenis Rekening (CRUD Management BMT)

```sql
CREATE TABLE jenis_rekening (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id                  UUID NOT NULL REFERENCES bmt(id),
    kode                    VARCHAR(20) NOT NULL,
    nama                    VARCHAR(100) NOT NULL,
    tipe_dasar              VARCHAR(30) NOT NULL,
    -- SIMPANAN_POKOK | SIMPANAN_WAJIB | SIMPANAN_SUKARELA | DEPOSITO | TABUNGAN_KHUSUS
    akad                    VARCHAR(30) NOT NULL,
    deskripsi               TEXT,
    setoran_awal_min        BIGINT NOT NULL DEFAULT 0,
    setoran_min             BIGINT NOT NULL DEFAULT 0,
    bisa_ditarik            BOOLEAN NOT NULL DEFAULT TRUE,
    syarat_penarikan        TEXT,
    nisbah_nasabah          SMALLINT,
    jangka_hari             SMALLINT,
    biaya_admin_bulanan     BIGINT NOT NULL DEFAULT 0,
    bisa_nfc                BOOLEAN NOT NULL DEFAULT FALSE,
    bisa_autodebet          BOOLEAN NOT NULL DEFAULT TRUE,
    biaya_admin_buka        BIGINT NOT NULL DEFAULT 0,
    is_aktif                BOOLEAN NOT NULL DEFAULT TRUE,
    urutan_tampil           SMALLINT NOT NULL DEFAULT 0,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by              UUID NOT NULL REFERENCES pengguna(id),
    updated_by              UUID NOT NULL REFERENCES pengguna(id),
    UNIQUE (bmt_id, kode)
);
```

---

## Rekening

```sql
CREATE TABLE rekening (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    cabang_id           UUID NOT NULL REFERENCES cabang(id),
    nasabah_id          UUID NOT NULL REFERENCES nasabah(id),
    jenis_rekening_id   UUID NOT NULL REFERENCES jenis_rekening(id),
    nomor_rekening      VARCHAR(40) UNIQUE NOT NULL,
    -- Format: {KODE_BMT}-{KODE_CAB}-{KODE_JENIS}-{SEQ:08d}
    saldo               BIGINT NOT NULL DEFAULT 0 CHECK (saldo >= 0),
    status              VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    -- AKTIF | BLOKIR | TUTUP
    alasan_blokir       TEXT,
    -- Snapshot saat buka (tidak berubah meski jenis_rekening diupdate)
    biaya_admin_bulanan BIGINT NOT NULL DEFAULT 0,
    nominal_deposito    BIGINT,
    nisbah_nasabah      SMALLINT,
    tanggal_buka        DATE NOT NULL,
    tanggal_jatuh_tempo DATE,
    tanggal_tutup       DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by_form_id  UUID REFERENCES form_pengajuan(id),
    updated_by_form_id  UUID REFERENCES form_pengajuan(id)
);
```

---

## Autodebet — Tanggal Diset Per Rekening

```sql
-- Konfigurasi autodebet per rekening (diset management BMT)
CREATE TABLE rekening_autodebet_config (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    rekening_id         UUID NOT NULL REFERENCES rekening(id),
    jenis               VARCHAR(30) NOT NULL,
    -- SIMPANAN_WAJIB | BIAYA_ADMIN_REKENING | ANGSURAN_PEMBIAYAAN | SPP_PONDOK
    tanggal_debet       SMALLINT NOT NULL,
    -- Tanggal dalam bulan (1-28) — diset management BMT
    -- Tanggal 29/30/31 otomatis disesuaikan ke akhir bulan
    is_aktif            BOOLEAN NOT NULL DEFAULT TRUE,
    referensi_id        UUID,              -- pembiayaan_id / tagihan_spp_id
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID NOT NULL REFERENCES pengguna(id)
);

CREATE TABLE jadwal_autodebet (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    rekening_id         UUID NOT NULL REFERENCES rekening(id),
    config_id           UUID NOT NULL REFERENCES rekening_autodebet_config(id),
    jenis               VARCHAR(30) NOT NULL,
    nominal_target      BIGINT NOT NULL,
    tanggal_jatuh_tempo DATE NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'MENUNGGU',
    -- MENUNGGU | SUKSES | PARTIAL | GAGAL
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tunggakan_autodebet (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    rekening_id         UUID NOT NULL REFERENCES rekening(id),
    jadwal_id           UUID NOT NULL REFERENCES jadwal_autodebet(id),
    jenis               VARCHAR(30) NOT NULL,
    nominal_target      BIGINT NOT NULL,
    nominal_terbayar    BIGINT NOT NULL DEFAULT 0,
    nominal_sisa        BIGINT NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'OUTSTANDING',
    -- OUTSTANDING | LUNAS
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Logika Autodebet Partial

```go
// Autodebet partial — jika saldo kurang, debit semampu saldo, sisa jadi tunggakan
berhasil := min(rekening.Saldo, jadwal.NominalTarget)
sisa     := jadwal.NominalTarget - berhasil
if sisa > 0 {
    tunggakanRepo.Insert(ctx, TunggakanAutodebet{
        NominalTarget:  jadwal.NominalTarget,
        NominalTerbayar: berhasil,
        NominalSisa:    sisa,
        Status:         "OUTSTANDING",
    })
}
```

---

## Pembiayaan & Beasiswa

```sql
CREATE TABLE pembiayaan (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id                  UUID NOT NULL REFERENCES bmt(id),
    cabang_id               UUID NOT NULL REFERENCES cabang(id),
    nasabah_id              UUID NOT NULL REFERENCES nasabah(id),
    produk_pembiayaan_id    UUID NOT NULL REFERENCES produk_pembiayaan(id),
    nomor_pembiayaan        VARCHAR(40) UNIQUE NOT NULL,
    akad                    VARCHAR(20) NOT NULL,
    -- MURABAHAH | MUSYARAKAH | MUDHARABAH | IJARAH | QARDH
    pokok                   BIGINT NOT NULL,
    margin_persen           NUMERIC(5,2),
    nisbah_nasabah          SMALLINT,
    jangka_bulan            SMALLINT NOT NULL,
    angsuran_per_bulan      BIGINT NOT NULL,
    total_kewajiban         BIGINT NOT NULL,
    -- Beasiswa (ditetapkan admin pondok)
    ada_beasiswa            BOOLEAN NOT NULL DEFAULT FALSE,
    beasiswa_persen         NUMERIC(5,2),
    beasiswa_nominal        BIGINT NOT NULL DEFAULT 0,
    beasiswa_sumber         VARCHAR(100),
    beasiswa_ditetapkan_oleh UUID REFERENCES pengguna_pondok(id),
    beasiswa_ditetapkan_at  TIMESTAMPTZ,
    -- Status
    status                  VARCHAR(20) NOT NULL DEFAULT 'PENGAJUAN',
    -- PENGAJUAN | DISETUJUI | AKTIF | LUNAS | MACET | DITOLAK
    kolektibilitas          SMALLINT NOT NULL DEFAULT 1,
    hari_tunggak            INT NOT NULL DEFAULT 0,
    saldo_pokok             BIGINT NOT NULL,
    saldo_margin            BIGINT NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by              UUID NOT NULL REFERENCES pengguna(id),
    updated_by              UUID NOT NULL REFERENCES pengguna(id),
    is_voided               BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE beasiswa_riwayat (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pembiayaan_id       UUID NOT NULL REFERENCES pembiayaan(id),
    persen_sebelum      NUMERIC(5,2),
    persen_sesudah      NUMERIC(5,2) NOT NULL,
    nominal_sebelum     BIGINT,
    nominal_sesudah     BIGINT NOT NULL,
    alasan              TEXT,
    ditetapkan_oleh     UUID NOT NULL REFERENCES pengguna_pondok(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Modul Akuntansi

Semua transaksi keuangan **wajib** memanggil `module-vernon-accounting`.

```go
journal.Post(ctx, Journal{
    BMTID:      bmtID,
    CabangID:   cabangID,
    Tanggal:    time.Now(),
    Keterangan: "Setoran tunai rekening ANNUR-KDR-SU-00000001",
    Referensi:  transaksiID.String(),
    Entries: []Entry{
        {KodeAkun: "101", Posisi: DEBIT,  Nominal: nominal},
        {KodeAkun: "202", Posisi: KREDIT, Nominal: nominal},
    },
})
```

### Chart of Accounts Default

```
1xx ASET      : 101 Kas, 102 Kas Bank, 111-114 Piutang/Pembiayaan
2xx KEWAJIBAN : 201 Simpanan Pokok/Wajib, 202 Sukarela, 203 Deposito, 211 Dana Sosial
3xx EKUITAS   : 301-302 Simpanan, 303 Cadangan, 304 SHU Ditahan
4xx PENDAPATAN: 401 Margin, 402-403 Bagi Hasil, 404 Ujrah, 405 Admin, 406 Komisi NFC, 407 Komisi OPOP
5xx BEBAN     : 501 Bagi Hasil, 502 Operasional, 503 Penyisihan, 504 Gaji, 505 Utilitas
6xx DANA SOSIAL:
    601 Penerimaan Donasi, 602 Penerimaan Wakaf, 603 Penerimaan Infaq, 604 Penerimaan Zakat
    611 Penyaluran Donasi, 612 Penyaluran Wakaf, 613 Penyaluran Infaq, 614 Penyaluran Zakat
```

---

## Checklist Syariah CBS

- [ ] Tidak ada riba — margin/nisbah/ujrah transparan & disepakati sebelum akad
- [ ] Ta'zir 100% masuk akun 211 — bukan pendapatan
- [ ] Bagi hasil dari realisasi pendapatan, bukan % nominal pokok
- [ ] Autodebet angsuran partial tetap menghasilkan jurnal syariah yang benar
- [ ] Biaya admin rekening: akad jelas saat buka rekening
