# Domain: Donasi, Wakaf & Infaq

## Prinsip Syariah

- **Donasi** — hibah sukarela tanpa imbalan
- **Wakaf** — aset yang dibekukan manfaatnya untuk kepentingan umum, dikelola BMT
- **Infaq/Shadaqah** — pengeluaran sukarela di jalan Allah (termasuk dana ta'zir)
- Dana sosial **harus terpisah** dari operasional BMT (rekening & akun akuntansi tersendiri)
- Dana ta'zir 100% masuk akun 211 (Dana Sosial) — bukan pendapatan BMT

---

## Skema Tabel

```sql
CREATE TABLE program_donasi (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    nama            VARCHAR(255) NOT NULL,
    deskripsi       TEXT,
    tipe            VARCHAR(20) NOT NULL,  -- DONASI | WAKAF | INFAQ | ZAKAT
    target_nominal  BIGINT,               -- NULL = tidak ada target
    terkumpul       BIGINT NOT NULL DEFAULT 0,
    tanggal_mulai   DATE NOT NULL,
    tanggal_selesai DATE,
    foto_url        TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    rekening_id     UUID NOT NULL REFERENCES rekening(id),
    -- Rekening khusus program (bisa beda per program)
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE transaksi_donasi (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    program_id      UUID NOT NULL REFERENCES program_donasi(id),
    nasabah_id      UUID REFERENCES nasabah(id),  -- NULL jika anonim
    nominal         BIGINT NOT NULL CHECK (nominal > 0),
    is_anonim       BOOLEAN NOT NULL DEFAULT FALSE,
    pesan           TEXT,
    metode          VARCHAR(30) NOT NULL,
    -- MIDTRANS | REKENING_BMT | NFC
    midtrans_order_id VARCHAR(100) UNIQUE,
    rekening_id     UUID REFERENCES rekening(id),
    idempotency_key UUID UNIQUE,
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    -- PENDING | SETTLEMENT | EXPIRE | CANCEL
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Aset wakaf yang dikelola BMT
CREATE TABLE aset_wakaf (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    nama            VARCHAR(255) NOT NULL,
    deskripsi       TEXT,
    jenis           VARCHAR(30) NOT NULL,
    -- TANAH | BANGUNAN | UANG | KENDARAAN | LAINNYA
    nilai_awal      BIGINT NOT NULL,
    wakif           VARCHAR(255),          -- nama pemberi wakaf
    nazhir          VARCHAR(255),          -- pengelola wakaf
    peruntukan      TEXT,                  -- tujuan wakaf
    dokumen_url     TEXT,                  -- akta wakaf (MinIO)
    tanggal_wakaf   DATE NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    -- AKTIF | DIKELOLA | SELESAI
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Hasil usaha wakaf produktif
CREATE TABLE hasil_wakaf (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aset_id     UUID NOT NULL REFERENCES aset_wakaf(id),
    periode     CHAR(7) NOT NULL,          -- "2025-01"
    pendapatan  BIGINT NOT NULL,
    beban       BIGINT NOT NULL DEFAULT 0,
    hasil_bersih BIGINT NOT NULL,
    distribusi  JSONB,                     -- ke mana hasil disalurkan
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Akun Akuntansi Dana Sosial

```
6xx  DANA SOSIAL & ZAKAT
  601  Penerimaan Donasi
  602  Penerimaan Wakaf
  603  Penerimaan Infaq/Shadaqah
  604  Penerimaan Zakat
  611  Penyaluran Donasi
  612  Penyaluran Wakaf
  613  Penyaluran Infaq/Shadaqah
  614  Penyaluran Zakat
```

---

## Checklist Syariah

- [ ] Donasi & wakaf: akad jelas, dana **tidak bercampur** dengan operasional BMT
- [ ] Wakaf produktif: hasil usaha dibagikan sesuai peruntukan wakaf, bukan ke BMT
- [ ] Infaq/shadaqah: penyaluran tercatat lengkap dengan mustahiq
- [ ] Dana wakaf tidak bercampur dengan modal BMT
- [ ] Denda perpustakaan: masuk dana sosial (akun 611), bukan pendapatan
