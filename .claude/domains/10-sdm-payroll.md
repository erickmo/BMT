# Domain: SDM & Payroll

## Prinsip

- Gaji diproses via **payroll autodebet** ke rekening BMT milik karyawan/guru
- Absensi terintegrasi — rekap ketidakhadiran otomatis mengurangi gaji
- Slip gaji digital tersedia di `apps-nasabah` jika karyawan terdaftar sebagai nasabah
- Tanggal gajian dikonfigurasi: `sdm.tanggal_gajian` (default: 25)

---

## Kontrak Kerja

```sql
CREATE TABLE sdm_kontrak (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    pegawai_id      UUID NOT NULL,
    tipe_pegawai    VARCHAR(20) NOT NULL,  -- PENGAJAR | KARYAWAN
    nomor_kontrak   VARCHAR(40) UNIQUE NOT NULL,
    tipe_kontrak    VARCHAR(20) NOT NULL,
    -- TETAP | KONTRAK | HONORER | MAGANG
    tanggal_mulai   DATE NOT NULL,
    tanggal_selesai DATE,                  -- NULL = tidak berbatas
    gaji_pokok      BIGINT NOT NULL,
    tunjangan       JSONB NOT NULL DEFAULT '{}',
    -- {"transport": 200000, "makan": 150000, "jabatan": 500000}
    potongan_tetap  JSONB NOT NULL DEFAULT '{}',
    -- {"bpjs_kesehatan": 50000, "bpjs_ketenagakerjaan": 30000}
    rekening_gaji_id UUID REFERENCES rekening(id),
    -- Rekening BMT tujuan transfer gaji
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    dokumen_url     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Slip Gaji

```sql
CREATE TABLE sdm_slip_gaji (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    kontrak_id      UUID NOT NULL REFERENCES sdm_kontrak(id),
    periode         CHAR(7) NOT NULL,       -- "2025-01"
    -- Pendapatan
    gaji_pokok      BIGINT NOT NULL,
    tunjangan_total BIGINT NOT NULL DEFAULT 0,
    tunjangan_detail JSONB NOT NULL DEFAULT '{}',
    -- Potongan
    potongan_absensi BIGINT NOT NULL DEFAULT 0,
    -- = hari_alfa × (gaji_pokok / hari_kerja_bulan)
    potongan_tetap  BIGINT NOT NULL DEFAULT 0,
    potongan_lain   BIGINT NOT NULL DEFAULT 0,
    -- Total
    gaji_bersih     BIGINT NOT NULL,
    -- Rekap absensi
    hari_kerja      SMALLINT NOT NULL,
    hari_hadir      SMALLINT NOT NULL,
    hari_sakit      SMALLINT NOT NULL,
    hari_izin       SMALLINT NOT NULL,
    hari_alfa       SMALLINT NOT NULL,
    -- Status
    status          VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    -- DRAFT | DISETUJUI | DIBAYAR
    dibayar_at      TIMESTAMPTZ,
    transaksi_id    UUID REFERENCES transaksi_rekening(id),
    file_url        TEXT,                   -- MinIO — PDF slip gaji
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (kontrak_id, periode)
);
```

---

## Alur Payroll

```
Worker GenerateSlipGaji (tgl 25, 09:00):
  1. Ambil semua kontrak aktif
  2. Hitung absensi bulan ini → potongan_absensi
  3. Generate slip_gaji (status: DRAFT)
  4. Notifikasi ke Finance untuk review

Finance review & setujui slip → status: DISETUJUI

Worker EksekusiPayroll (tgl 1, 08:00):
  1. Untuk setiap slip DISETUJUI:
  2. Transfer dari kas BMT ke rekening_gaji pegawai
  3. INSERT transaksi_rekening (jenis: KREDIT_PAYROLL)
  4. POST jurnal: Debit 504 (Beban Gaji) / Kredit 101 (Kas)
  5. Update slip → status: DIBAYAR
  6. Push notifikasi + email slip ke pegawai
```

---

## Checklist Syariah SDM

- [ ] Gaji pegawai: tidak ada unsur riba — gaji flat, bukan % keuntungan
- [ ] Akad kerja jelas (kontrak tersimpan di DB + dokumen digital)
