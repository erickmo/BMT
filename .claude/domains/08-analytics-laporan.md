# Domain: Analytics & Laporan

## Dashboard Real-Time

Data dashboard di-cache Redis dengan TTL: `analytics.dashboard_cache_ttl_detik` (default: 30 detik).

```
GET /ws/dashboard   → WebSocket (memerlukan JWT)
GET /sse/dashboard  → SSE fallback

Payload yang di-stream:
- Jumlah transaksi hari ini (update tiap transaksi baru)
- Total setoran & penarikan hari ini
- Sesi teller aktif
- Alert fraud baru
- Pesanan OPOP baru masuk
```

```sql
-- Snapshot metrik harian (digenerate worker jam 00:05)
CREATE TABLE analytics_snapshot (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID REFERENCES cabang(id),  -- NULL = konsolidasi BMT
    tanggal         DATE NOT NULL,
    metrik          VARCHAR(50) NOT NULL,
    -- TOTAL_NASABAH | TOTAL_DPK | TOTAL_PEMBIAYAAN | TRANSAKSI_HARI_INI
    -- NPF_RATIO | KOLEKTIBILITAS_1..5 | OPOP_PENJUALAN | dll.
    nilai           NUMERIC NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (bmt_id, cabang_id, tanggal, metrik)
);
```

---

## Laporan Custom

User dapat memilih kolom, filter, dan rentang tanggal sendiri.
Semua laporan dapat di-export ke **PDF** dan **Excel/CSV**.

```sql
CREATE TABLE laporan_template (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id      UUID NOT NULL REFERENCES bmt(id),
    nama        VARCHAR(100) NOT NULL,
    domain      VARCHAR(30) NOT NULL,
    -- NASABAH | REKENING | TRANSAKSI | PEMBIAYAAN | ABSENSI | NILAI | OPOP
    kolom       JSONB NOT NULL,
    -- ["nomor_nasabah", "nama", "saldo", ...]
    filter      JSONB NOT NULL DEFAULT '{}',
    urutan      JSONB NOT NULL DEFAULT '[]',
    dibuat_oleh UUID NOT NULL REFERENCES pengguna(id),
    is_publik   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Format default export: `laporan.default_format` di settings.

---

## Analytics OPOP

```sql
CREATE TABLE opop_analytics_harian (
    tanggal             DATE NOT NULL,
    toko_id             UUID NOT NULL REFERENCES toko(id),
    bmt_id              UUID NOT NULL REFERENCES bmt(id),
    total_pesanan       INT NOT NULL DEFAULT 0,
    total_pendapatan    BIGINT NOT NULL DEFAULT 0,
    total_item_terjual  INT NOT NULL DEFAULT 0,
    pengunjung_unik     INT NOT NULL DEFAULT 0,
    PRIMARY KEY (tanggal, toko_id)
);
```

Worker `AnalyticsHarian` berjalan setiap 23:00 WIB untuk snapshot harian.
Worker `GenerateSnapshotAnalytics` berjalan 00:05 WIB untuk snapshot metrik kesiswaan.
