-- Inventaris
CREATE TABLE inventaris_aset (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    kode_aset       VARCHAR(30) UNIQUE NOT NULL,
    nama            VARCHAR(255) NOT NULL,
    kategori        VARCHAR(30) NOT NULL,
    lokasi          VARCHAR(100),
    tanggal_perolehan DATE NOT NULL,
    nilai_perolehan BIGINT NOT NULL,
    nilai_buku      BIGINT NOT NULL,
    umur_ekonomis   SMALLINT,
    kondisi         VARCHAR(20) NOT NULL DEFAULT 'BAIK',
    foto_url        TEXT,
    kode_akun       VARCHAR(10) NOT NULL DEFAULT '131',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE inventaris_peminjaman (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    aset_id         UUID NOT NULL REFERENCES inventaris_aset(id),
    peminjam_id     UUID NOT NULL,
    peminjam_tipe   VARCHAR(20) NOT NULL,
    keperluan       TEXT NOT NULL,
    tanggal_pinjam  TIMESTAMPTZ NOT NULL,
    tanggal_kembali_rencana TIMESTAMPTZ NOT NULL,
    tanggal_kembali_aktual  TIMESTAMPTZ,
    status          VARCHAR(20) NOT NULL DEFAULT 'DIPINJAM',
    kondisi_kembali VARCHAR(20),
    catatan         TEXT,
    disetujui_oleh  UUID REFERENCES pengguna_pondok(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- PPDB
CREATE TABLE ppdb_pendaftar (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    cabang_id           UUID NOT NULL REFERENCES cabang(id),
    tahun_ajaran        VARCHAR(10) NOT NULL,
    nama_lengkap        VARCHAR(255) NOT NULL,
    nik                 VARCHAR(16),
    tanggal_lahir       DATE,
    nama_wali           VARCHAR(255) NOT NULL,
    telepon_wali        VARCHAR(20) NOT NULL,
    email_wali          VARCHAR(255),
    pilihan_tingkat     VARCHAR(20),
    status              VARCHAR(20) NOT NULL DEFAULT 'DAFTAR',
    nomor_pendaftaran   VARCHAR(30) UNIQUE NOT NULL,
    dokumen             JSONB NOT NULL DEFAULT '{}',
    catatan             TEXT,
    nasabah_id          UUID REFERENCES nasabah(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Perpustakaan
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
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    buku_id         UUID NOT NULL REFERENCES perpus_buku(id),
    santri_id       UUID NOT NULL REFERENCES pondok_santri(id),
    tanggal_pinjam  DATE NOT NULL,
    tanggal_kembali_rencana DATE NOT NULL,
    tanggal_kembali_aktual  DATE,
    status          VARCHAR(20) NOT NULL DEFAULT 'DIPINJAM',
    denda           BIGINT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Surat Izin
CREATE TABLE surat_izin (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    santri_id       UUID NOT NULL REFERENCES pondok_santri(id),
    jenis           VARCHAR(20) NOT NULL,
    keperluan       TEXT NOT NULL,
    tanggal_mulai   TIMESTAMPTZ NOT NULL,
    tanggal_kembali TIMESTAMPTZ NOT NULL,
    tujuan          VARCHAR(255),
    diajukan_oleh   VARCHAR(20) NOT NULL,
    nasabah_wali_id UUID REFERENCES nasabah(id),
    status          VARCHAR(20) NOT NULL DEFAULT 'MENUNGGU',
    disetujui_oleh  UUID REFERENCES pengguna_pondok(id),
    alasan_tolak    TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Health Record
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

-- Alumni
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
    foto_url        TEXT,
    is_verified     BOOLEAN NOT NULL DEFAULT FALSE,
    nasabah_id      UUID REFERENCES nasabah(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ppdb_bmt ON ppdb_pendaftar(bmt_id, tahun_ajaran, status);
CREATE INDEX idx_perpus_peminjaman ON perpus_peminjaman(santri_id, status);
