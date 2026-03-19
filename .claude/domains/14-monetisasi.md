# Domain: Monetisasi Platform

## Model Revenue Platform

| Sumber | Mekanisme | Dikonfigurasi |
|--------|-----------|---------------|
| Biaya admin per transaksi | Flat + % per jenis transaksi | Developer per kontrak BMT |
| Komisi OPOP | % dari nilai transaksi pesanan | `monetisasi.komisi_opop_persen` di platform_settings |
| Fitur premium OPOP | Iklan/featured toko di marketplace | Paket iklan per BMT |
| White-label | Custom branding app per pondok | Kontrak tersendiri per BMT |

---

## Komisi OPOP

Setiap pesanan OPOP status `SELESAI` otomatis mencatat komisi:

```sql
CREATE TABLE opop_komisi (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pesanan_id      UUID NOT NULL REFERENCES pesanan(id) UNIQUE,
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    nilai_pesanan   BIGINT NOT NULL,
    persen_komisi   NUMERIC(5,2) NOT NULL,
    -- Snapshot dari kontrak saat transaksi
    nominal_komisi  BIGINT NOT NULL,
    -- = nilai_pesanan * persen / 100
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    -- PENDING | DITAGIH | LUNAS
    periode         CHAR(7) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Worker `HitungKomisiOPOP` berjalan tgl 1, 07:00 — rekap komisi bulanan ke `usage_log`.

---

## Iklan Premium OPOP

```sql
CREATE TABLE opop_iklan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    toko_id         UUID NOT NULL REFERENCES toko(id),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    tipe            VARCHAR(20) NOT NULL,
    -- FEATURED_TOKO | BANNER_PRODUK | TOP_SEARCH
    tanggal_mulai   DATE NOT NULL,
    tanggal_selesai DATE NOT NULL,
    nominal         BIGINT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## White-Label

Konfigurasi branding per BMT di `bmt_settings`:

```
whitelabel.nama_app           → "Santri Pay"
whitelabel.bundle_id_android  → "com.annur.santripay"
whitelabel.bundle_id_ios      → "com.annur.santripay"
whitelabel.primary_color      → "#1B5E20"
whitelabel.logo_url           → "https://..."
whitelabel.splash_url         → "https://..."
```

Juga disimpan di kolom `bmt.whitelabel` (JSONB) untuk akses cepat.

> **Catatan:** White-label membutuhkan **build pipeline terpisah** per BMT yang mengambil
> asset dari settings saat build Flutter (GitHub Actions matrix per BMT / Codemagic).

---

## Checklist Syariah Monetisasi

- [ ] Komisi OPOP: akad wakalah bil ujrah — transparan di awal antara platform dan pondok
- [ ] Biaya admin rekening: akad jelas saat buka rekening, bukan surprise fee
