# Domain: Inventaris & Aset

## Aset Tetap

```sql
CREATE TABLE inventaris_aset (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    kode_aset       VARCHAR(30) UNIQUE NOT NULL,
    nama            VARCHAR(255) NOT NULL,
    kategori        VARCHAR(30) NOT NULL,
    -- GEDUNG | KENDARAAN | PERALATAN | FURNITUR | ELEKTRONIK | LAINNYA
    lokasi          VARCHAR(100),
    tanggal_perolehan DATE NOT NULL,
    nilai_perolehan BIGINT NOT NULL,
    nilai_buku      BIGINT NOT NULL,       -- setelah penyusutan
    umur_ekonomis   SMALLINT,              -- tahun
    kondisi         VARCHAR(20) NOT NULL DEFAULT 'BAIK',
    -- BAIK | RUSAK_RINGAN | RUSAK_BERAT | TIDAK_AKTIF
    foto_url        TEXT,
    kode_akun       VARCHAR(10) NOT NULL DEFAULT '131',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Worker `DepresiasiAset` berjalan tgl 1 tahunan — hitung & posting depresiasi ke jurnal akuntansi.

---

## Peminjaman Aset / Ruang

```sql
CREATE TABLE inventaris_peminjaman (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id                      UUID NOT NULL REFERENCES bmt(id),
    aset_id                     UUID NOT NULL REFERENCES inventaris_aset(id),
    peminjam_id                 UUID NOT NULL,
    peminjam_tipe               VARCHAR(20) NOT NULL,
    -- SANTRI | PENGGUNA_PONDOK | EXTERNAL
    keperluan                   TEXT NOT NULL,
    tanggal_pinjam              TIMESTAMPTZ NOT NULL,
    tanggal_kembali_rencana     TIMESTAMPTZ NOT NULL,
    tanggal_kembali_aktual      TIMESTAMPTZ,
    status                      VARCHAR(20) NOT NULL DEFAULT 'DIPINJAM',
    -- DIPINJAM | DIKEMBALIKAN | TERLAMBAT
    kondisi_kembali             VARCHAR(20),
    catatan                     TEXT,
    disetujui_oleh              UUID REFERENCES pengguna_pondok(id),
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```
