# 09 — E-commerce OPOP

> **Terakhir diperbarui:** 20 Maret 2026

Feature Gate: `ECOMMERCE_OPOP`

## Model Bisnis
- **B2C** — Wali santri beli dari toko pondok anaknya
- **B2B (OPOP)** — Pondok A beli dari Pondok B (marketplace lintas pondok)
- **Seller** — Pondok & BMT/koperasi pondok
- **Buyer** — Wali santri (B2C), pondok lain (B2B)

## Toko & Produk
```sql
toko (bmt_id, cabang_id, nama, slug UNIQUE,
      kategori_toko,  -- PONDOK|BMT_KOPERASI
      is_opop BOOLEAN,  -- tampil di marketplace lintas pondok
      status, rating NUMERIC(3,2))

produk (bmt_id, toko_id, nama, slug, kategori,
        harga BIGINT, harga_b2b BIGINT,
        stok INT, satuan, berat_gram,
        foto_urls JSONB,
        is_opop BOOLEAN, is_aktif BOOLEAN,
        rating, total_terjual)
```

## Pesanan
```sql
pesanan (buyer_tipe,         -- WALI_SANTRI|PONDOK
         nasabah_id,         -- jika WALI_SANTRI
         bmt_buyer_id,       -- jika PONDOK (B2B)
         toko_id, bmt_seller_id,
         nomor_pesanan UNIQUE,
         status,             -- MENUNGGU_PEMBAYARAN|DIBAYAR|DIPROSES|DIKIRIM|SELESAI|DIBATALKAN
         subtotal, ongkir, total,
         alamat_kirim JSONB, kurir, nomor_resi,
         metode_bayar,       -- MIDTRANS|REKENING_BMT|NFC
         catatan)

pesanan_item (pesanan_id, produk_id,
              nama_produk, harga,  -- snapshot saat pesan
              jumlah, subtotal)

pembayaran_pesanan (pesanan_id, metode, nominal, status,
                    midtrans_order_id, rekening_id, kartu_nfc_id,
                    idempotency_key UNIQUE, settled_at)

ulasan_produk (produk_id, pesanan_id, nasabah_id,
               rating SMALLINT CHECK (1-5),
               komentar, foto_urls JSONB,
               UNIQUE per pesanan+produk)
```

## Monetisasi OPOP
```sql
-- Komisi per pesanan selesai
komisi_opop (bmt_seller_id, pesanan_id,
             gmv BIGINT,
             persen_komisi NUMERIC(5,2),  -- snapshot dari kontrak
             nominal_komisi BIGINT,
             periode CHAR(7),
             status)  -- PENDING|TAGIH|LUNAS

-- Iklan/slot premium
opop_iklan (toko_id, bmt_id,
            jenis,          -- BANNER|PRODUK_UNGGULAN|TOKO_FEATURED
            tanggal_mulai, tanggal_selesai,
            biaya BIGINT, status)

-- Analytics harian (dihitung worker malam)
opop_analytics_harian (toko_id, tanggal DATE,
                        total_pesanan INT, gmv BIGINT,
                        total_pembeli INT, rating_avg NUMERIC(3,2),
                        UNIQUE per toko+tanggal)
```

## OPOP Marketplace Lintas Pondok
- `GET /opop/toko` — semua toko `is_opop = true` lintas BMT
- `GET /opop/produk` — semua produk `is_opop = true`
- `POST /opop/pesanan` — B2B order antar pondok
- Search: PostgreSQL FTS (fase 1), Elasticsearch (fase 2)
