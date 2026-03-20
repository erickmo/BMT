# 12 — Notifikasi & Komunikasi

> **Terakhir diperbarui:** 20 Maret 2026

## Channel (Feature: `KOMUNIKASI_WA`, `KOMUNIKASI_SMS`)

| Channel | Kegunaan | Provider (dari settings) |
|---------|----------|--------------------------|
| Push FCM | Notifikasi real-time in-app | `notifikasi.fcm_server_key` |
| WhatsApp Personal | Konfirmasi transaksi, reminder | `notifikasi.wa_provider` (fonnte/wablas) |
| WhatsApp Blast | Pengumuman massal | Same |
| SMS | OTP 2FA, H-1 angsuran | `notifikasi.sms_provider` (zenziva/twilio) |
| Email | Dokumen PDF, raport, slip | `notifikasi.email_smtp_*` |
| Bulletin Board | Pengumuman resmi pondok | In-app DB |

Semua provider dikonfigurasi di **BMT settings**, bukan hardcode.

## Template Notifikasi
```sql
notifikasi_template (bmt_id,  -- NULL = global platform
                     kode VARCHAR(50),
                     channel VARCHAR(20),  -- FCM|WA|SMS|EMAIL
                     judul VARCHAR(255),
                     isi TEXT,             -- {{nama}}, {{nominal}}, {{tanggal}}
                     is_aktif BOOLEAN,
                     UNIQUE (bmt_id, kode, channel))
```

## Bulletin Board Pondok
```sql
bulletin (bmt_id, cabang_id, judul, isi TEXT, foto_urls JSONB,
          target,              -- SEMUA|SANTRI|WALI|PENGAJAR|KARYAWAN
          target_kelas_id,
          target_asrama,
          pin_until TIMESTAMPTZ,
          is_aktif BOOLEAN,
          diterbitkan_oleh, diterbitkan_at)
```
Tampil di `app/nasabah` (tab pengumuman) dan `app/pondok`.

## Chat Internal Staf BMT
```sql
chat_room (bmt_id, nama, tipe,  -- DIRECT|GRUP
           anggota JSONB)
chat_pesan (room_id, pengirim_id, isi, lampiran_url,
            is_dibaca BOOLEAN, created_at)
```
Scope per BMT — staf satu BMT tidak bisa chat dengan BMT lain.
Tersedia di `app/management`.

## Event Notifikasi Otomatis

| Event | Channel | Trigger |
|-------|---------|---------|
| Transaksi setor/tarik | FCM + WA | Real-time |
| Angsuran H-3 | WA + Email | Worker harian |
| Angsuran H-1 | SMS | Worker harian |
| SPP H-3 | WA + Email | Worker harian |
| Autodebet gagal | FCM + WA | Saat worker |
| Tunggakan | Email | Setelah partial debit |
| Raport diterbitkan | FCM + WA | Setelah DITERBITKAN |
| Pesanan OPOP update | FCM | Perubahan status |
| Slip gaji tersedia | FCM | Setelah payroll |
| Surat izin disetujui | FCM + WA | Setelah approval |
| Fraud alert | FCM (manajer) | Real-time |
| OTP | SMS / Email | Login / transaksi besar |
| Listing langganan expired | Email (owner listing) | H-7 |
| SaaS invoice terbit | Email (ADMIN_BMT) | Awal bulan |
