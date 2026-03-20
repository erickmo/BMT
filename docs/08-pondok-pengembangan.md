# 08 — ERP Pondok: Pengembangan Santri, SDM & Event

> **Terakhir diperbarui:** 20 Maret 2026

## Pengembangan Santri

### Portfolio
```sql
santri_portfolio (santri_id, bmt_id, judul, deskripsi,
                  kategori,  -- KARYA|SERTIFIKAT|PRESTASI|PENGHARGAAN
                  file_url TEXT, tanggal DATE,
                  is_publik BOOLEAN)
```

### Tracking Hafalan (Timeline Progress)
```sql
-- Berbeda dari pondok_nilai_tahfidz (per ujian)
-- Ini tracking kemajuan hafalan dari waktu ke waktu
santri_hafalan_progress (santri_id, bmt_id, tanggal DATE,
                          total_juz NUMERIC(4,1),  -- misal: 3.5
                          total_surah SMALLINT, catatan, dicatat_oleh UUID)
-- Menghasilkan grafik progress di app/nasabah
```

### Program Hafalan Target
```sql
hafalan_program (bmt_id, cabang_id, nama, target_juz NUMERIC(4,1),
                 periode_mulai, periode_selesai, untuk_tingkat, is_aktif)

santri_hafalan_target (santri_id, program_id,
                       target_personal NUMERIC(4,1),
                       progress_terakhir NUMERIC(4,1),
                       status)  -- AKTIF|TERCAPAI|TIDAK_TERCAPAI
```

### Ekstrakurikuler & Organisasi
```sql
ekstra_kegiatan (bmt_id, cabang_id, nama, deskripsi,
                 pembina_id, jadwal JSONB, is_aktif)

santri_ekstra (santri_id, ekstra_id,
               peran VARCHAR(50),  -- ANGGOTA|PENGURUS|KETUA
               tanggal_bergabung, tanggal_keluar, is_aktif)
```

### Alumni & Bimbingan Karir
```sql
alumni (santri_id, nama_lengkap, angkatan, tahun_lulus,
        email, telepon, pekerjaan, instansi, domisili,
        nasabah_id,  -- alumni bisa jadi nasabah BMT
        is_aktif_jaringan BOOLEAN)

-- Alumni sebagai mentor
alumni_mentor (alumni_id, bidang_keahlian, tersedia_konsultasi, deskripsi)

mentoring_sesi (santri_id, alumni_id, topik, jadwal TIMESTAMPTZ,
                status, catatan)  -- MENUNGGU|TERJADWAL|SELESAI
```

---

## Event Pondok

```sql
pondok_event (bmt_id, cabang_id, nama, deskripsi,
              kategori,    -- SEMINAR|WISUDA|HAFLAH|LOMBA|OLAHRAGA|LAINNYA
              tanggal_mulai, tanggal_selesai, lokasi,
              kapasitas SMALLINT,
              harga_tiket BIGINT DEFAULT 0,  -- 0 = gratis
              foto_url, status)              -- DRAFT|PUBLIKASI|SELESAI|BATAL

event_registrasi (event_id, registran_id, registran_tipe,
                  nomor_tiket VARCHAR(30) UNIQUE,
                  status_bayar,   -- BELUM_BAYAR|LUNAS|GRATIS
                  midtrans_order_id, rekening_id,
                  hadir BOOLEAN DEFAULT FALSE)
```

### Program Beasiswa Eksternal
```sql
beasiswa_eksternal_program (bmt_id, nama_program, donatur_nama,
                             donatur_kontak, total_dana BIGINT,
                             dana_tersalurkan BIGINT,
                             tanggal_mulai, tanggal_selesai, status)

beasiswa_eksternal_penerima (program_id, santri_id,
                              nominal_per_periode BIGINT,
                              tagihan_spp_id, pembiayaan_id)
```

---

## SDM & Penggajian (Feature: `SDM_PAYROLL`) — via module-vernon-hrm

```sql
-- Data kontrak disimpan di module-vernon-hrm/contract
-- Diakses via: hrm.GetContract(ctx, contractID)
sdm_kontrak (bmt_id, cabang_id, karyawan_id | pengajar_id,
             tipe_kontrak,       -- TETAP|KONTRAK|HONORER|PARUH_WAKTU
             gaji_pokok BIGINT,
             tunjangan JSONB,    -- {transportasi, makan, jabatan, ...}
             potongan_per_hari_absen BIGINT,
             rekening_gaji_id UUID,  -- rekening BMT tujuan gaji
             tanggal_mulai, tanggal_selesai, is_aktif)

sdm_slip_gaji (kontrak_id, periode CHAR(7),
               gaji_pokok, tunjangan_total, tunjangan_detail JSONB,
               potongan_absen, hari_absen SMALLINT,
               potongan_lain, potongan_detail JSONB,
               gaji_bersih,
               status_transfer,  -- PENDING|DIPROSES|BERHASIL|GAGAL
               transaksi_id,
               file_url)         -- PDF slip di MinIO
```

**Worker `PayrollBulanan`** (tanggal dari `settings: "sdm.tanggal_gajian"`):
1. Ambil semua kontrak aktif
2. Hitung rekap absensi → hitung potongan
3. Generate `sdm_slip_gaji`
4. Transfer via autodebet ke `rekening_gaji_id`
5. Generate PDF → push FCM ke nasabah/karyawan

---

## Donasi, Wakaf & Zakat

```sql
donasi_program (bmt_id, cabang_id, judul, jenis,
                -- DONASI|INFAQ|SHADAQAH|WAKAF_UANG
                target_nominal, terkumpul, tanggal_mulai, tanggal_selesai, status)

donasi_transaksi (program_id, nasabah_id, nama_donatur,
                  nominal, metode,  -- REKENING_BMT|MIDTRANS|NFC|TUNAI
                  rekening_id, is_anonim, idempotency_key UNIQUE)

wakaf_aset (bmt_id, jenis_aset,  -- UANG|TANAH|BANGUNAN|USAHA
            nilai_aset, nama_wakif, mauquf_alaih TEXT, status)
```

### Prinsip Syariah Keuangan
- Donasi/infaq — tidak ada imbalan materi
- Wakaf produktif — BMT sebagai nazir, dana tidak bercampur modal BMT
- Ta'zir — 100% masuk akun 211 (dana sosial)
- Denda perpustakaan — masuk dana sosial, bukan pendapatan
- Zakat maal — dihitung worker akhir tahun jika keuntungan ≥ nisab
