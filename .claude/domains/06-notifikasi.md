# Domain: Notifikasi & Komunikasi

## Prinsip

Semua channel notifikasi dikonfigurasi **via settings BMT** — tidak ada provider yang hardcode di kode.

```go
// Settings menentukan provider, bukan kode
waProvider := settings.Resolve(ctx, bmtID, cabangID, "notifikasi.wa_provider")
// → "fonnte" | "wablas" | "whacenter"
```

---

## Channel yang Didukung

| Channel | Kegunaan | Konfigurasi |
|---------|----------|-------------|
| **FCM Push** | Notifikasi real-time in-app | `notifikasi.fcm_server_key` di platform_settings |
| **WhatsApp Personal** | Konfirmasi transaksi, OTP, reminder | `notifikasi.wa_provider`, `notifikasi.wa_token` di bmt_settings |
| **WhatsApp Blast** | Pengumuman massal ke semua wali/santri | — |
| **SMS** | OTP 2FA, reminder darurat | `notifikasi.sms_provider`, `notifikasi.sms_apikey` |
| **Email** | Dokumen PDF, laporan, slip | `notifikasi.email_smtp_*` di bmt_settings |
| **Bulletin Board** | Pengumuman resmi pondok dalam app | Tersimpan di DB |

---

## Skema Tabel

```sql
-- Template notifikasi (bisa dikustomisasi per BMT)
CREATE TABLE notifikasi_template (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id      UUID REFERENCES bmt(id),   -- NULL = template global platform
    kode        VARCHAR(50) NOT NULL,
    -- TRANSAKSI_SETOR | ANGSURAN_JATUH_TEMPO | OTP_LOGIN | dll.
    channel     VARCHAR(20) NOT NULL,       -- FCM | WA | SMS | EMAIL
    judul       VARCHAR(255),
    isi         TEXT NOT NULL,
    -- Mendukung variabel: {{nama}}, {{nominal}}, {{tanggal}}, dll.
    is_aktif    BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (bmt_id, kode, channel)
);

-- Antrian notifikasi (diproses worker setiap 1 menit)
CREATE TABLE notifikasi_antrian (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    channel         VARCHAR(20) NOT NULL,   -- FCM | WHATSAPP | SMS | EMAIL
    tujuan          VARCHAR(255) NOT NULL,  -- token FCM / no. WA / no. HP / email
    subjek          VARCHAR(255),
    pesan           TEXT NOT NULL,
    data_ekstra     JSONB,                  -- payload FCM / attachment email
    status          VARCHAR(20) NOT NULL DEFAULT 'MENUNGGU',
    -- MENUNGGU | TERKIRIM | GAGAL
    percobaan       SMALLINT NOT NULL DEFAULT 0,
    error_terakhir  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    dikirim_at      TIMESTAMPTZ
);

-- Log pengiriman notifikasi
CREATE TABLE notifikasi_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    template_kode   VARCHAR(50) NOT NULL,
    channel         VARCHAR(20) NOT NULL,
    tujuan          VARCHAR(255) NOT NULL,
    isi_terkirim    TEXT NOT NULL,
    status          VARCHAR(20) NOT NULL,   -- TERKIRIM | GAGAL | PENDING
    error_message   TEXT,
    referensi_id    UUID,
    referensi_tipe  VARCHAR(30),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Bulletin Board (Pengumuman In-App)

```sql
CREATE TABLE pengumuman (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    judul           VARCHAR(255) NOT NULL,
    isi             TEXT NOT NULL,
    tipe            VARCHAR(20) NOT NULL,
    -- SEMUA | SANTRI | WALI | PENGAJAR | KARYAWAN | KELAS | ASRAMA
    target_id       UUID,
    -- kelas_id atau asrama jika tipe KELAS/ASRAMA
    file_url        TEXT,
    is_pinned       BOOLEAN NOT NULL DEFAULT FALSE,
    tanggal_mulai   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    tanggal_selesai TIMESTAMPTZ,
    dibuat_oleh     UUID NOT NULL REFERENCES pengguna_pondok(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pengumuman_baca (
    pengumuman_id   UUID NOT NULL REFERENCES pengumuman(id),
    nasabah_id      UUID REFERENCES nasabah(id),
    pengguna_id     UUID REFERENCES pengguna_pondok(id),
    dibaca_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (pengumuman_id, COALESCE(nasabah_id, pengguna_id))
);
```

---

## Worker: KirimNotifikasi

- Berjalan setiap **1 menit**
- Proses antrian `notifikasi_antrian` status `MENUNGGU`
- Retry max 3x dengan backoff eksponensial (dari settings)
- Provider dipilih berdasarkan settings BMT, bukan hardcode
