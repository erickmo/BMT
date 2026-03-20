# 07 — ERP Pondok: Operasional

> **Terakhir diperbarui:** 20 Maret 2026

## Feature Gates
- `PONDOK_PERPUS` — perpustakaan (add-on)
- `PONDOK_UKS` — health record (add-on)
- `PONDOK_PPDB` — penerimaan santri baru (add-on)

## Manajemen Asrama & Kamar
```sql
asrama (bmt_id, cabang_id, kode, nama, kapasitas, jenis_kelamin,
        pengurus_id, lantai, fasilitas JSONB)

kamar (asrama_id, nomor_kamar, kapasitas, is_aktif)

santri_kamar (santri_id, kamar_id,
              tanggal_masuk DATE, tanggal_keluar DATE,
              is_aktif BOOLEAN,
              UNIQUE aktif per santri)
-- History perpindahan kamar tersimpan
```
Wali santri lihat info kamar di `app/nasabah`. Kelola di `app/pondok`.

## Perpustakaan Digital (Feature: `PONDOK_PERPUS`)
```sql
perpus_buku (bmt_id, cabang_id, isbn, judul, pengarang, kategori,
             stok_total, stok_tersedia, foto_url)

perpus_peminjaman (buku_id, peminjam_id, peminjam_tipe,
                   tanggal_pinjam, tanggal_kembali_rencana,
                   tanggal_kembali_aktual,
                   status,   -- DIPINJAM|DIKEMBALIKAN|TERLAMBAT|HILANG
                   denda BIGINT)
-- Denda → akun 211 (dana sosial), bukan pendapatan
```

## Konsultasi Online
```sql
konsultasi_sesi (pemohon_id, pemohon_tipe,  -- SANTRI|WALI
                 konselor_id,
                 kategori,                  -- AKADEMIK|BK|KESEHATAN|KEUANGAN|UMUM
                 jadwal TIMESTAMPTZ, status,
                 catatan_privat TEXT)        -- hanya konselor bisa baca

konsultasi_pesan (sesi_id, pengirim_id, pengirim_tipe, isi, lampiran_url)
```

## Surat Izin Digital
```sql
surat_izin (santri_id,
            jenis,          -- IZIN_PULANG|IZIN_KELUAR|IZIN_SAKIT|IZIN_KEGIATAN
            keperluan, tanggal_keluar, tanggal_kembali,
            diajukan_oleh,  -- SANTRI|WALI
            status,         -- MENUNGGU|DISETUJUI|DITOLAK|SELESAI
            disetujui_oleh, alasan_tolak,
            waktu_keluar_aktual, waktu_kembali_aktual)
```

## Health Record UKS (Feature: `PONDOK_UKS`)
```sql
health_record (santri_id, tanggal,
               jenis_kunjungan, -- SAKIT|PEMERIKSAAN_RUTIN|KECELAKAAN|RUJUKAN
               keluhan, diagnosa, tindakan, obat_diberikan,
               petugas_uks_id,
               perlu_rujukan BOOLEAN, fasilitas_rujukan)
```

## Inventaris Aset
```sql
aset (bmt_id, cabang_id, kode_aset UNIQUE, nama,
      kategori,         -- GEDUNG|KENDARAAN|PERALATAN|FURNITUR|ELEKTRONIK|TANAH
      nilai_perolehan BIGINT, tanggal_perolehan DATE,
      umur_ekonomis SMALLINT, nilai_buku BIGINT,
      lokasi, kondisi,  -- BAIK|RUSAK_RINGAN|RUSAK_BERAT|TIDAK_LAYAK
      foto_url)

aset_peminjaman (aset_id, peminjam_id, peminjam_tipe,
                 keperluan, tanggal_pinjam,
                 tanggal_kembali_rencana, tanggal_kembali_aktual,
                 status, disetujui_oleh)
```
Worker `DepresiasiAset` jalan 1 Januari — hitung & posting jurnal depresiasi.

## PPDB — Penerimaan Santri Baru (Feature: `PONDOK_PPDB`)
```sql
ppdb_gelombang (bmt_id, cabang_id, nama, tahun_ajaran,
                tanggal_buka, tanggal_tutup, kuota, biaya_daftar,
                persyaratan JSONB)

ppdb_pendaftar (gelombang_id, nomor_pendaftaran UNIQUE,
                nama_lengkap, tanggal_lahir, asal_sekolah,
                nama_wali, telepon_wali, email_wali,
                dokumen_urls JSONB,
                status,       -- MENDAFTAR|SELEKSI|DITERIMA|DITOLAK|MENGUNDURKAN_DIRI
                status_bayar, midtrans_order_id,
                santri_id,    -- diisi setelah diterima
                nasabah_id)   -- diisi setelah buka rekening BMT
```

## Integrasi DAPODIK & EMIS (Feature: `INTEGRASI_DAPODIK` / `INTEGRASI_EMIS`)
```sql
sinkronisasi_eksternal (bmt_id, sumber,  -- DAPODIK|EMIS
                         jenis,          -- SANTRI|LEMBAGA|GURU
                         status, jumlah_record, berhasil, gagal,
                         error_detail JSONB)
```
Konfigurasi NPSN/NSM/token di BMT settings (`is_rahasia = true`).
Worker `SinkronDAPODIK` dan `SinkronEMIS` jalan sesuai `integrasi.sinkron_jadwal`.
