# Domain: ERP Pondok

## Administrasi Santri

```sql
CREATE TABLE pondok_santri (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    cabang_id           UUID NOT NULL REFERENCES cabang(id),
    nomor_induk_santri  VARCHAR(30) NOT NULL,
    nama_lengkap        VARCHAR(255) NOT NULL,
    nasabah_id          UUID REFERENCES nasabah(id),   -- 1 nasabah = 1 santri
    -- Kesiswaan
    tingkat             VARCHAR(20),       -- MTS | MA | S1 | TAHFIDZ | dll.
    kelas_id            UUID REFERENCES pondok_kelas(id),
    asrama              VARCHAR(100),
    kamar               VARCHAR(20),
    angkatan            SMALLINT,
    status_aktif        BOOLEAN NOT NULL DEFAULT TRUE,
    tanggal_masuk       DATE,
    tanggal_keluar      DATE,
    foto_url            TEXT,
    -- Wali
    nama_wali           VARCHAR(255),
    telepon_wali        VARCHAR(20),
    nasabah_wali_id     UUID REFERENCES nasabah(id),
    -- Biometrik
    fingerprint_template TEXT,             -- template sidik jari (encrypted)
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (bmt_id, nomor_induk_santri)
);

CREATE TABLE pondok_kelas (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id      UUID NOT NULL REFERENCES bmt(id),
    cabang_id   UUID NOT NULL REFERENCES cabang(id),
    nama        VARCHAR(50) NOT NULL,      -- "7A", "MA Kelas 1"
    tingkat     VARCHAR(20) NOT NULL,
    tahun_ajaran VARCHAR(10) NOT NULL,     -- "2025/2026"
    wali_kelas_id UUID REFERENCES pondok_pengajar(id),
    kapasitas   SMALLINT,
    UNIQUE (bmt_id, cabang_id, nama, tahun_ajaran)
);

CREATE TABLE pondok_pengajar (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    nip             VARCHAR(30),
    nama_lengkap    VARCHAR(255) NOT NULL,
    jabatan         VARCHAR(100),
    spesialisasi    VARCHAR(100),
    nasabah_id      UUID REFERENCES nasabah(id),
    fingerprint_template TEXT,
    status_aktif    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (bmt_id, nip)
);

CREATE TABLE pondok_karyawan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    nik_karyawan    VARCHAR(30),
    nama_lengkap    VARCHAR(255) NOT NULL,
    jabatan         VARCHAR(100),
    departemen      VARCHAR(100),
    nasabah_id      UUID REFERENCES nasabah(id),
    fingerprint_template TEXT,
    status_aktif    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Kurikulum & Akademik

```sql
CREATE TABLE pondok_mapel (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id      UUID NOT NULL REFERENCES bmt(id),
    cabang_id   UUID NOT NULL REFERENCES cabang(id),
    kode        VARCHAR(20) NOT NULL,
    nama        VARCHAR(100) NOT NULL,
    tingkat     VARCHAR(20) NOT NULL,
    is_aktif    BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (bmt_id, cabang_id, kode)
);

CREATE TABLE pondok_silabus (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    mapel_id        UUID NOT NULL REFERENCES pondok_mapel(id),
    tahun_ajaran    VARCHAR(10) NOT NULL,
    semester        SMALLINT NOT NULL,     -- 1 | 2
    deskripsi       TEXT,
    file_url        TEXT,                  -- MinIO
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID NOT NULL REFERENCES pengguna_pondok(id)
);

CREATE TABLE pondok_rpp (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    silabus_id      UUID NOT NULL REFERENCES pondok_silabus(id),
    pengajar_id     UUID NOT NULL REFERENCES pondok_pengajar(id),
    pertemuan_ke    SMALLINT NOT NULL,
    topik           VARCHAR(255) NOT NULL,
    tujuan          TEXT,
    materi_url      TEXT,                  -- MinIO
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Komponen penilaian per mapel per semester
CREATE TABLE pondok_komponen_nilai (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    mapel_id        UUID NOT NULL REFERENCES pondok_mapel(id),
    tahun_ajaran    VARCHAR(10) NOT NULL,
    semester        SMALLINT NOT NULL,
    nama            VARCHAR(50) NOT NULL,  -- "UH1", "UTS", "UAS", "Tugas"
    bobot_persen    SMALLINT NOT NULL,     -- % dari nilai akhir
    UNIQUE (bmt_id, mapel_id, tahun_ajaran, semester, nama)
);
```

---

## Jadwal

```sql
CREATE TABLE pondok_kalender (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id      UUID NOT NULL REFERENCES bmt(id),
    cabang_id   UUID NOT NULL REFERENCES cabang(id),
    tahun_ajaran VARCHAR(10) NOT NULL,
    tanggal     DATE NOT NULL,
    jenis       VARCHAR(20) NOT NULL,
    -- LIBUR | UJIAN | ACARA | HARI_EFEKTIF | LIBUR_NASIONAL
    keterangan  VARCHAR(255),
    UNIQUE (bmt_id, cabang_id, tanggal)
);

CREATE TABLE pondok_jadwal_pelajaran (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    kelas_id        UUID NOT NULL REFERENCES pondok_kelas(id),
    mapel_id        UUID NOT NULL REFERENCES pondok_mapel(id),
    pengajar_id     UUID NOT NULL REFERENCES pondok_pengajar(id),
    hari            SMALLINT NOT NULL,     -- 1=Senin ... 7=Minggu
    jam_mulai       TIME NOT NULL,
    jam_selesai     TIME NOT NULL,
    ruangan         VARCHAR(50),
    tahun_ajaran    VARCHAR(10) NOT NULL,
    semester        SMALLINT NOT NULL
);

CREATE TABLE pondok_jadwal_kegiatan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    nama            VARCHAR(255) NOT NULL,
    kategori        VARCHAR(30) NOT NULL,
    -- PENGAJIAN | OLAHRAGA | EKSTRA | RAPAT | ACARA | LAINNYA
    tanggal_mulai   TIMESTAMPTZ NOT NULL,
    tanggal_selesai TIMESTAMPTZ,
    lokasi          VARCHAR(100),
    peserta         VARCHAR(20) NOT NULL,
    -- SEMUA | SANTRI | PENGAJAR | KARYAWAN | TERTENTU
    deskripsi       TEXT,
    created_by      UUID NOT NULL REFERENCES pengguna_pondok(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pondok_jadwal_piket (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    santri_id       UUID NOT NULL REFERENCES pondok_santri(id),
    jenis_piket     VARCHAR(50) NOT NULL,
    hari            SMALLINT NOT NULL,
    lokasi          VARCHAR(100),
    periode_mulai   DATE NOT NULL,
    periode_selesai DATE NOT NULL
);

CREATE TABLE pondok_jadwal_shift (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    pengguna_id     UUID,
    jenis_pengguna  VARCHAR(20) NOT NULL,  -- PENGAJAR | KARYAWAN
    tanggal         DATE NOT NULL,
    jam_masuk       TIME NOT NULL,
    jam_keluar      TIME NOT NULL,
    keterangan      TEXT
);
```

---

## Absensi

```sql
CREATE TABLE pondok_absensi (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    subjek_id       UUID NOT NULL,         -- santri_id / pengajar_id / karyawan_id
    subjek_tipe     VARCHAR(20) NOT NULL,  -- SANTRI | PENGAJAR | KARYAWAN
    tanggal         DATE NOT NULL,
    sesi            VARCHAR(20),           -- PAGI | SIANG | MALAM | MAPEL_ID
    jadwal_id       UUID,
    status          VARCHAR(20) NOT NULL,
    -- HADIR | SAKIT | IZIN | ALFA | TERLAMBAT
    keterangan      TEXT,
    metode          VARCHAR(20) NOT NULL,  -- dari settings, bukan hardcode
    -- MANUAL | NFC | BIOMETRIK
    waktu_scan      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES pengguna_pondok(id)
);
```

**Metode absensi dari settings:**
```go
metodeDiizinkan := settings.ResolveJSON(ctx, bmtID, cabangID, "pondok.absensi_metode")
// → ["MANUAL", "NFC", "BIOMETRIK"] — dari DB, bukan konstanta
```

---

## Penilaian & Raport

```sql
CREATE TABLE pondok_nilai (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    santri_id           UUID NOT NULL REFERENCES pondok_santri(id),
    komponen_id         UUID NOT NULL REFERENCES pondok_komponen_nilai(id),
    nilai               NUMERIC(5,2) NOT NULL CHECK (nilai >= 0 AND nilai <= 100),
    catatan             TEXT,
    diinput_oleh        UUID NOT NULL REFERENCES pengguna_pondok(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (santri_id, komponen_id)
);

CREATE TABLE pondok_nilai_tahfidz (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    santri_id       UUID NOT NULL REFERENCES pondok_santri(id),
    surah           VARCHAR(100) NOT NULL,
    ayat_mulai      SMALLINT NOT NULL,
    ayat_selesai    SMALLINT NOT NULL,
    nilai           NUMERIC(5,2),
    status          VARCHAR(20) NOT NULL,  -- LULUS | MENGULANG | BELUM_DIUJI
    tanggal_ujian   DATE NOT NULL,
    penguji_id      UUID REFERENCES pondok_pengajar(id),
    catatan         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pondok_nilai_akhlak (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    santri_id       UUID NOT NULL REFERENCES pondok_santri(id),
    tanggal         DATE NOT NULL,
    jenis           VARCHAR(20) NOT NULL,  -- PELANGGARAN | PRESTASI
    kategori        VARCHAR(50) NOT NULL,
    deskripsi       TEXT NOT NULL,
    poin            SMALLINT NOT NULL,     -- positif atau negatif
    dicatat_oleh    UUID NOT NULL REFERENCES pengguna_pondok(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pondok_raport (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    santri_id       UUID NOT NULL REFERENCES pondok_santri(id),
    kelas_id        UUID NOT NULL REFERENCES pondok_kelas(id),
    tahun_ajaran    VARCHAR(10) NOT NULL,
    semester        SMALLINT NOT NULL,
    nilai_mapel     JSONB NOT NULL,        -- [{mapel_id, nama, nilai_akhir, predikat}, ...]
    nilai_tahfidz   JSONB,
    nilai_akhlak    JSONB,
    total_hadir     SMALLINT,
    total_sakit     SMALLINT,
    total_izin      SMALLINT,
    total_alfa      SMALLINT,
    peringkat       SMALLINT,
    catatan_wali_kelas TEXT,
    file_url        TEXT,                  -- MinIO — PDF raport
    status          VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    -- DRAFT | FINAL | DITERBITKAN
    diterbitkan_at  TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Tagihan & Keuangan Pondok

```sql
CREATE TABLE pondok_jenis_tagihan (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id      UUID NOT NULL REFERENCES bmt(id),
    kode        VARCHAR(20) NOT NULL,
    nama        VARCHAR(100) NOT NULL,
    nominal     BIGINT NOT NULL,
    frekuensi   VARCHAR(20) NOT NULL,      -- BULANAN | TAHUNAN | SEKALI | CUSTOM
    is_aktif    BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (bmt_id, kode)
);

CREATE TABLE tagihan_spp (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    cabang_id           UUID NOT NULL REFERENCES cabang(id),
    santri_id           UUID NOT NULL REFERENCES pondok_santri(id),
    jenis_tagihan_id    UUID NOT NULL REFERENCES pondok_jenis_tagihan(id),
    periode             CHAR(7) NOT NULL,  -- "2025-01"
    nominal             BIGINT NOT NULL,
    nominal_terbayar    BIGINT NOT NULL DEFAULT 0,
    nominal_sisa        BIGINT NOT NULL,
    beasiswa_persen     NUMERIC(5,2) NOT NULL DEFAULT 0,
    beasiswa_nominal    BIGINT NOT NULL DEFAULT 0,
    nominal_efektif     BIGINT NOT NULL,   -- = nominal - beasiswa_nominal
    status              VARCHAR(20) NOT NULL DEFAULT 'BELUM_BAYAR',
    tanggal_jatuh_tempo DATE NOT NULL,
    tanggal_lunas       DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```
