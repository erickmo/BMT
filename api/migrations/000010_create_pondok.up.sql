CREATE TABLE pondok_santri (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id                  UUID NOT NULL REFERENCES bmt(id),
    cabang_id               UUID NOT NULL REFERENCES cabang(id),
    nomor_induk_santri      VARCHAR(30) NOT NULL,
    nama_lengkap            VARCHAR(255) NOT NULL,
    nasabah_id              UUID REFERENCES nasabah(id),
    tingkat                 VARCHAR(20),
    kelas_id                UUID,
    asrama                  VARCHAR(100),
    kamar                   VARCHAR(20),
    angkatan                SMALLINT,
    status_aktif            BOOLEAN NOT NULL DEFAULT TRUE,
    tanggal_masuk           DATE,
    tanggal_keluar          DATE,
    foto_url                TEXT,
    nama_wali               VARCHAR(255),
    telepon_wali            VARCHAR(20),
    nasabah_wali_id         UUID REFERENCES nasabah(id),
    fingerprint_template    TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (bmt_id, nomor_induk_santri)
);

CREATE TABLE pondok_kelas (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    nama            VARCHAR(50) NOT NULL,
    tingkat         VARCHAR(20) NOT NULL,
    tahun_ajaran    VARCHAR(10) NOT NULL,
    wali_kelas_id   UUID,
    kapasitas       SMALLINT,
    UNIQUE (bmt_id, cabang_id, nama, tahun_ajaran)
);

ALTER TABLE pondok_santri ADD CONSTRAINT fk_santri_kelas FOREIGN KEY (kelas_id) REFERENCES pondok_kelas(id);

CREATE TABLE pondok_pengajar (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id                  UUID NOT NULL REFERENCES bmt(id),
    cabang_id               UUID NOT NULL REFERENCES cabang(id),
    nip                     VARCHAR(30),
    nama_lengkap            VARCHAR(255) NOT NULL,
    jabatan                 VARCHAR(100),
    spesialisasi            VARCHAR(100),
    nasabah_id              UUID REFERENCES nasabah(id),
    fingerprint_template    TEXT,
    status_aktif            BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (bmt_id, nip)
);

ALTER TABLE pondok_kelas ADD CONSTRAINT fk_kelas_walikelas FOREIGN KEY (wali_kelas_id) REFERENCES pondok_pengajar(id);

CREATE TABLE pondok_karyawan (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id                  UUID NOT NULL REFERENCES bmt(id),
    cabang_id               UUID NOT NULL REFERENCES cabang(id),
    nik_karyawan            VARCHAR(30),
    nama_lengkap            VARCHAR(255) NOT NULL,
    jabatan                 VARCHAR(100),
    departemen              VARCHAR(100),
    nasabah_id              UUID REFERENCES nasabah(id),
    fingerprint_template    TEXT,
    status_aktif            BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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

CREATE TABLE pondok_jadwal_pelajaran (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    kelas_id        UUID NOT NULL REFERENCES pondok_kelas(id),
    mapel_id        UUID NOT NULL REFERENCES pondok_mapel(id),
    pengajar_id     UUID NOT NULL REFERENCES pondok_pengajar(id),
    hari            SMALLINT NOT NULL,
    jam_mulai       TIME NOT NULL,
    jam_selesai     TIME NOT NULL,
    ruangan         VARCHAR(50),
    tahun_ajaran    VARCHAR(10) NOT NULL,
    semester        SMALLINT NOT NULL
);

CREATE TABLE pondok_absensi (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    subjek_id       UUID NOT NULL,
    subjek_tipe     VARCHAR(20) NOT NULL,
    tanggal         DATE NOT NULL,
    sesi            VARCHAR(20),
    jadwal_id       UUID REFERENCES pondok_jadwal_pelajaran(id),
    status          VARCHAR(20) NOT NULL,
    keterangan      TEXT,
    metode          VARCHAR(20) NOT NULL,
    waktu_scan      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES pengguna_pondok(id)
);

CREATE TABLE pondok_komponen_nilai (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    mapel_id        UUID NOT NULL REFERENCES pondok_mapel(id),
    tahun_ajaran    VARCHAR(10) NOT NULL,
    semester        SMALLINT NOT NULL,
    nama            VARCHAR(50) NOT NULL,
    bobot_persen    SMALLINT NOT NULL,
    UNIQUE (bmt_id, mapel_id, tahun_ajaran, semester, nama)
);

CREATE TABLE pondok_nilai (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    santri_id       UUID NOT NULL REFERENCES pondok_santri(id),
    komponen_id     UUID NOT NULL REFERENCES pondok_komponen_nilai(id),
    nilai           NUMERIC(5,2) NOT NULL CHECK (nilai >= 0 AND nilai <= 100),
    catatan         TEXT,
    diinput_oleh    UUID NOT NULL REFERENCES pengguna_pondok(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
    status          VARCHAR(20) NOT NULL,
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
    jenis           VARCHAR(20) NOT NULL,
    kategori        VARCHAR(50) NOT NULL,
    deskripsi       TEXT NOT NULL,
    poin            SMALLINT NOT NULL,
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
    nilai_mapel     JSONB NOT NULL,
    nilai_tahfidz   JSONB,
    nilai_akhlak    JSONB,
    total_hadir     SMALLINT,
    total_sakit     SMALLINT,
    total_izin      SMALLINT,
    total_alfa      SMALLINT,
    peringkat       SMALLINT,
    catatan_wali_kelas TEXT,
    file_url        TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    diterbitkan_at  TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pondok_jenis_tagihan (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id      UUID NOT NULL REFERENCES bmt(id),
    kode        VARCHAR(20) NOT NULL,
    nama        VARCHAR(100) NOT NULL,
    nominal     BIGINT NOT NULL,
    frekuensi   VARCHAR(20) NOT NULL,
    is_aktif    BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (bmt_id, kode)
);

CREATE TABLE tagihan_spp (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    cabang_id           UUID NOT NULL REFERENCES cabang(id),
    santri_id           UUID NOT NULL REFERENCES pondok_santri(id),
    jenis_tagihan_id    UUID NOT NULL REFERENCES pondok_jenis_tagihan(id),
    periode             CHAR(7) NOT NULL,
    nominal             BIGINT NOT NULL,
    nominal_terbayar    BIGINT NOT NULL DEFAULT 0,
    nominal_sisa        BIGINT NOT NULL,
    beasiswa_persen     NUMERIC(5,2) NOT NULL DEFAULT 0,
    beasiswa_nominal    BIGINT NOT NULL DEFAULT 0,
    nominal_efektif     BIGINT NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'BELUM_BAYAR',
    tanggal_jatuh_tempo DATE NOT NULL,
    tanggal_lunas       DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pondok_santri_bmt ON pondok_santri(bmt_id, cabang_id, status_aktif);
CREATE INDEX idx_tagihan_spp_santri ON tagihan_spp(santri_id, periode);
CREATE INDEX idx_absensi_santri ON pondok_absensi(subjek_id, tanggal);
