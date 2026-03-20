# 10 — Keuangan & Akuntansi

> **Terakhir diperbarui:** 20 Maret 2026

## module-vernon-accounting

Semua transaksi keuangan **wajib** posting jurnal. Jurnal tidak pernah dihapus — koreksi via reversal entry.

```go
journal.Post(ctx, Journal{
    BMTID: bmtID, CabangID: cabangID,
    Entries: []Entry{
        {KodeAkun: "101", Posisi: DEBIT,  Nominal: nominal},
        {KodeAkun: "202", Posisi: KREDIT, Nominal: nominal},
    },
}) // error jika Σdebit ≠ Σkredit
```

## Laporan yang Digenerate

| Laporan | Periode | Format |
|---------|---------|--------|
| Neraca | Bulanan/Tahunan | PDF, Excel |
| Laporan SHU | Bulanan/Tahunan | PDF, Excel |
| Laporan Arus Kas | Bulanan/Tahunan | PDF, Excel |
| Kolektibilitas | Harian/Bulanan | PDF, Excel |
| Bagi Hasil Deposito | Bulanan | PDF |
| Distribusi Dana Sosial | Bulanan/Tahunan | PDF |
| Transaksi Harian | Harian | PDF |
| Buku Besar | Bebas | PDF, Excel |
| Posisi DPK | Bulanan | PDF |

### Laporan RAT (Detail)
- Laporan Pertanggungjawaban Pengurus
- Neraca per 31 Desember
- Laporan SHU Tahunan + Rencana Pembagian SHU
- Perkembangan Anggota (masuk/keluar/aktif)
- Laporan Kolektibilitas Tahunan

### Laporan Custom
```sql
laporan_template (bmt_id, cabang_id, nama, domain,
                  -- TRANSAKSI|NASABAH|PEMBIAYAAN|KESISWAAN|ABSENSI|OPOP
                  kolom JSONB, filter JSONB, urutan JSONB, dibuat_oleh)
```
User pilih kolom, filter, rentang tanggal → generate on-demand.

## Dokumen Akad & Histori TTD
```sql
dokumen (bmt_id, cabang_id,
         jenis,           -- AKAD_SIMPANAN|AKAD_DEPOSITO|AKAD_MURABAHAH|...
                          -- SLIP_SETORAN|SLIP_PENARIKAN|KARTU_ANGSURAN|SURAT_KUASA
         referensi_id, referensi_tipe,
         nomor_dokumen VARCHAR(60) UNIQUE,
         file_url TEXT,
         ttd_pejabat_url TEXT,
         status_ttd,      -- MENUNGGU_TTD|SUDAH_TTD|TIDAK_PERLU_TTD
         email_terkirim BOOLEAN, email_tujuan)

-- Histori versi dokumen
dokumen_versi (dokumen_id, versi SMALLINT,
               file_url, ttd_pejabat_url,
               keterangan, created_by, created_at)
```
Format nomor: `{KODE_BMT}/{KODE_CAB}/{KODE_AKAD}/{TAHUN}/{SEQ:05d}`

## Kartu Anggota & Buku Tabungan Digital
Generate on-demand via chromedp dari template HTML per BMT:
- Kartu anggota: PDF A6 — foto, nama, nomor, QR code verifikasi
- Buku tabungan: PDF A5 — riwayat 3 bulan terakhir per rekening

```
GET /api/dokumen/kartu-anggota/:nasabah_id
GET /api/dokumen/buku-tabungan/:rekening_id
```

## Dashboard Real-Time
WebSocket/SSE untuk update live:
- Jumlah transaksi hari ini
- Total setoran & penarikan
- Sesi teller aktif
- Fraud alert baru
- Pesanan OPOP baru
