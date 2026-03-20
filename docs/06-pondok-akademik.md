# 06 — ERP Pondok: Akademik

> **Terakhir diperbarui:** 20 Maret 2026

## Feature Gates
- `PONDOK_ADMINISTRASI` — santri, guru, kelas (termasuk FREE)
- `PONDOK_RAPORT` — penilaian + raport digital (add-on)
- `PONDOK_TAHFIDZ` — manajemen hafalan (add-on)

## Administrasi

### Santri
```sql
pondok_santri (bmt_id, cabang_id, nomor_induk_santri UNIQUE per bmt,
               nasabah_id UUID REFERENCES nasabah(id),  -- nullable
               tingkat VARCHAR(20), kelas_id UUID, asrama VARCHAR(100), kamar VARCHAR(20),
               angkatan SMALLINT, nama_wali, telepon_wali,
               nasabah_wali_id UUID REFERENCES nasabah(id),
               fingerprint_template TEXT,  -- encrypted
               status_aktif BOOLEAN)
```

### Kelas & Pengajar
```sql
pondok_kelas (bmt_id, cabang_id, nama, tingkat, tahun_ajaran,
              wali_kelas_id, kapasitas)
pondok_pengajar (bmt_id, cabang_id, nip UNIQUE, nama_lengkap,
                 jabatan, spesialisasi, nasabah_id, fingerprint_template)
pondok_karyawan (bmt_id, cabang_id, nik_karyawan, nama_lengkap,
                 jabatan, departemen, nasabah_id, fingerprint_template)
```

## Kurikulum
```
pondok_mapel (kode, nama, tingkat)
  → pondok_silabus (per mapel per semester, file_url MinIO)
      → pondok_rpp (per pertemuan, materi_url MinIO)
          → pondok_komponen_nilai (UH1/UTS/UAS/Tugas — bobot % dari DB)
```

**Penting:** bobot komponen nilai dikonfigurasi per mapel per semester di DB — bukan hardcode.

## Jadwal

| Tabel | Isi |
|-------|-----|
| `pondok_kalender` | Kalender akademik (libur, ujian, acara, hari efektif) |
| `pondok_jadwal_pelajaran` | Per kelas per hari (mapel, guru, jam, ruangan) |
| `pondok_jadwal_kegiatan` | Pengajian, olahraga, acara — target: kelas/asrama/semua |
| `pondok_jadwal_piket` | Piket santri per hari |
| `pondok_jadwal_shift` | Shift pengajar & karyawan |

## Absensi

**Metode dari settings** (`pondok.absensi_metode`) — bukan konstanta:
- `MANUAL` — guru input via `app/pondok`
- `NFC` — santri tap kartu NFC di kelas
- `BIOMETRIK` — scan sidik jari via perangkat

```sql
pondok_absensi (subjek_id, subjek_tipe,       -- SANTRI|PENGAJAR|KARYAWAN
                tanggal, sesi VARCHAR(20),    -- PAGI|SIANG|MALAM|mapel_id
                jadwal_id,
                status,                       -- HADIR|SAKIT|IZIN|ALFA|TERLAMBAT
                metode,                       -- dari settings
                waktu_scan TIMESTAMPTZ)
```

## Penilaian & Raport (Feature: `PONDOK_RAPORT`)

```sql
pondok_nilai (santri_id, komponen_id, nilai NUMERIC(5,2), catatan)

-- Tahfidz (Feature: PONDOK_TAHFIDZ)
pondok_nilai_tahfidz (santri_id, surah, ayat_mulai, ayat_selesai,
                      nilai, status, penguji_id, tanggal_ujian)

-- Akhlak — poin prestasi/pelanggaran
pondok_nilai_akhlak (santri_id, tanggal, jenis,  -- PELANGGARAN|PRESTASI
                     kategori, deskripsi, poin SMALLINT)

-- Raport digital
pondok_raport (santri_id, kelas_id, tahun_ajaran, semester,
               nilai_mapel JSONB, nilai_tahfidz JSONB, nilai_akhlak JSONB,
               total_hadir, total_sakit, total_izin, total_alfa,
               peringkat SMALLINT, catatan_wali_kelas,
               file_url TEXT,      -- PDF MinIO
               status,             -- DRAFT|FINAL|DITERBITKAN
               diterbitkan_at)
```
Wali santri lihat raport di `app/nasabah` setelah status `DITERBITKAN`.

## Keuangan Pondok

```sql
pondok_jenis_tagihan (bmt_id, kode, nama, nominal, frekuensi)

tagihan_spp (santri_id, jenis_tagihan_id, periode CHAR(7),
             nominal, beasiswa_persen, beasiswa_nominal,
             nominal_efektif,    -- = nominal - beasiswa_nominal
             nominal_terbayar, nominal_sisa,
             status,             -- BELUM_BAYAR|SEBAGIAN|LUNAS
             tanggal_jatuh_tempo)
```
Beasiswa SPP ditetapkan `ADMIN_PONDOK` / `BENDAHARA_PONDOK`.
Pembayaran via: `AUTODEBET|TELLER_TUNAI|TELLER_REKENING|NASABAH_APP|NFC`
