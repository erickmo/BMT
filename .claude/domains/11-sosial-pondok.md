# Domain: Fitur Sosial Pondok

## Perpustakaan Digital

```sql
CREATE TABLE perpus_buku (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    judul           VARCHAR(255) NOT NULL,
    pengarang       VARCHAR(255),
    penerbit        VARCHAR(100),
    tahun           SMALLINT,
    isbn            VARCHAR(20),
    kategori        VARCHAR(50),
    jumlah_total    SMALLINT NOT NULL DEFAULT 1,
    jumlah_tersedia SMALLINT NOT NULL DEFAULT 1,
    cover_url       TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE perpus_peminjaman (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id                      UUID NOT NULL REFERENCES bmt(id),
    buku_id                     UUID NOT NULL REFERENCES perpus_buku(id),
    santri_id                   UUID NOT NULL REFERENCES pondok_santri(id),
    tanggal_pinjam              DATE NOT NULL,
    tanggal_kembali_rencana     DATE NOT NULL,
    tanggal_kembali_aktual      DATE,
    status                      VARCHAR(20) NOT NULL DEFAULT 'DIPINJAM',
    -- DIPINJAM | DIKEMBALIKAN | TERLAMBAT
    denda                       BIGINT NOT NULL DEFAULT 0,
    -- Denda masuk dana sosial (akun 611), bukan pendapatan
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Worker `ReminderPeminjamanBuku` berjalan 08:00 daily untuk H-1 jatuh tempo.

---

## Konsultasi Online

```sql
CREATE TABLE konsultasi_sesi (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    penanya_id      UUID NOT NULL,
    penanya_tipe    VARCHAR(20) NOT NULL,   -- SANTRI | WALI
    penjawab_id     UUID REFERENCES pengguna_pondok(id),
    topik           VARCHAR(30) NOT NULL,
    -- AKADEMIK | BK | KESEHATAN | KEUANGAN | UMUM
    judul           VARCHAR(255) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    -- OPEN | DIJAWAB | DITUTUP
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE konsultasi_pesan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sesi_id         UUID NOT NULL REFERENCES konsultasi_sesi(id),
    pengirim_id     UUID NOT NULL,
    pengirim_tipe   VARCHAR(20) NOT NULL,   -- SANTRI | WALI | PENGGUNA_PONDOK
    pesan           TEXT NOT NULL,
    file_url        TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Surat Izin Digital

```sql
CREATE TABLE surat_izin (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id                  UUID NOT NULL REFERENCES bmt(id),
    santri_id               UUID NOT NULL REFERENCES pondok_santri(id),
    jenis                   VARCHAR(20) NOT NULL,
    -- KELUAR | PULANG | SAKIT | LAINNYA
    keperluan               TEXT NOT NULL,
    tanggal_mulai           TIMESTAMPTZ NOT NULL,
    tanggal_kembali         TIMESTAMPTZ NOT NULL,
    tujuan                  VARCHAR(255),
    diajukan_oleh           VARCHAR(20) NOT NULL,  -- SANTRI | WALI
    nasabah_wali_id         UUID REFERENCES nasabah(id),
    status                  VARCHAR(20) NOT NULL DEFAULT 'MENUNGGU',
    -- MENUNGGU | DISETUJUI | DITOLAK | DIBATALKAN
    disetujui_oleh          UUID REFERENCES pengguna_pondok(id),
    alasan_tolak            TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Health Record Santri (UKS)

```sql
CREATE TABLE kesehatan_santri (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    santri_id       UUID NOT NULL REFERENCES pondok_santri(id),
    golongan_darah  CHAR(3),
    alergi          TEXT[],
    riwayat_penyakit TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (santri_id)
);

CREATE TABLE kesehatan_kunjungan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    santri_id       UUID NOT NULL REFERENCES pondok_santri(id),
    tanggal         TIMESTAMPTZ NOT NULL,
    keluhan         TEXT NOT NULL,
    diagnosa        TEXT,
    tindakan        TEXT,
    obat            TEXT,
    rujukan         BOOLEAN NOT NULL DEFAULT FALSE,
    dicatat_oleh    UUID NOT NULL REFERENCES pengguna_pondok(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Alumni Management

```sql
CREATE TABLE alumni (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    santri_id       UUID REFERENCES pondok_santri(id),
    nama_lengkap    VARCHAR(255) NOT NULL,
    angkatan        SMALLINT NOT NULL,
    tahun_lulus     SMALLINT NOT NULL,
    pekerjaan       VARCHAR(100),
    instansi        VARCHAR(255),
    kota_domisili   VARCHAR(100),
    telepon         VARCHAR(20),
    email           VARCHAR(255),
    linkedin_url    TEXT,
    foto_url        TEXT,
    is_verified     BOOLEAN NOT NULL DEFAULT FALSE,
    nasabah_id      UUID REFERENCES nasabah(id),
    -- Alumni bisa tetap jadi nasabah BMT
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```
