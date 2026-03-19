# Domain: API Endpoints

## Developer

```
POST|GET|PUT  /dev/bmt / /dev/bmt/:id / /dev/bmt/:id/status
POST|GET|PUT  /dev/bmt/:id/kontrak / /dev/bmt/:id/cabang
POST          /dev/bmt/:id/pengguna/seed
GET|POST|PUT|DELETE /dev/pecahan-uang / /dev/pecahan-uang/:id
GET|POST|PUT  /dev/tarif-template / /dev/tarif-template/:id
GET|PUT       /dev/platform-settings / /dev/platform-settings/:kunci
GET           /dev/usage-log / /dev/usage-log/:bmt_id
GET           /dev/metrics / /dev/health/detail
POST          /dev/maintenance
```

## Auth

```
POST /auth/staf/login|refresh|logout
POST /auth/nasabah/login|refresh
POST /auth/pondok/login|refresh
POST /auth/merchant/login|refresh
```

## Platform (Management BMT)

```
GET|PUT       /platform/settings / /platform/settings/:kunci
GET|POST|PUT  /platform/jenis-rekening / /platform/jenis-rekening/:id
GET|POST|PUT  /platform/produk/simpanan|pembiayaan / /:id
GET|POST|PUT  /platform/pengguna / /platform/pengguna/:id
GET|POST|PUT  /platform/cabang / /platform/cabang/:id
GET|POST|PUT  /platform/merchant / /platform/merchant/:id
GET|POST|PUT  /platform/terminal-kiosk / /platform/terminal-kiosk/:id
GET|POST      /platform/pengguna-pondok / /:id
GET           /platform/laporan/konsolidasi|rat
GET           /platform/usage
```

## Teller

```
POST  /teller/sesi/buka              # Input redenominasi (pecahan dari DB)
GET   /teller/sesi/aktif
POST  /teller/sesi/tutup             # Ditolak jika ada selisih
GET   /teller/nasabah/cari
POST  /teller/rekening/:id/setor|tarik
POST  /teller/pembiayaan/:id/angsuran
POST  /teller/spp/:id/bayar
```

## Form Workflow

```
POST|GET  /form / /form/:id
PUT       /form/:id
POST      /form/:id/ajukan|setujui|tolak
GET       /form/:id/riwayat
```

## Management API

```
GET         /api/nasabah / /:id              # View only
GET         /api/rekening / /:id / /:id/transaksi|tunggakan
GET|POST|PUT /api/rekening/:id/autodebet-config
GET         /api/autodebet/jadwal|tunggakan
GET         /api/pembiayaan / /:id / /:id/jadwal
PUT         /api/pembiayaan/:id/analisis
POST        /api/pembiayaan/:id/pencairan|simulasi
PUT         /api/pembiayaan/:id/beasiswa     # Hanya ADMIN_PONDOK via pondok-app
POST|GET    /finance/jurnal|transaksi
GET         /api/laporan/neraca|shu|arus-kas|kolektibilitas|bagi-hasil-deposito
```

## Nasabah App

```
GET   /nasabah/profil|rekening|pembiayaan
GET   /nasabah/santri                  # Info santri + raport
GET   /nasabah/rekening/:id/transaksi
POST  /nasabah/rekening/:id/setor-online
GET   /nasabah/nfc/saldo|transaksi
POST  /nasabah/nfc/topup
GET   /nasabah/spp/tagihan
POST  /nasabah/spp/:id/bayar
# E-commerce
GET   /nasabah/shop/toko / /nasabah/shop/toko/:slug/produk
POST  /nasabah/shop/keranjang|pesanan
GET   /nasabah/shop/pesanan / /:id
POST  /nasabah/shop/ulasan
```

## NFC / Kiosk

```
POST  /nfc/transaksi                   # X-Idempotency-Key wajib
GET   /nfc/ceksaldo/:uid               # Kiosk — IP whitelist, tanpa PIN
```

## Pondok

```
GET|POST|PUT|DELETE /pondok/santri / /:id
GET|POST|PUT|DELETE /pondok/pengajar|karyawan / /:id
GET|POST|PUT        /pondok/kelas / /:id
GET|POST|PUT        /pondok/mapel|silabus|rpp / /:id
GET|POST|PUT        /pondok/komponen-nilai / /:id
GET|POST|PUT        /pondok/jadwal/pelajaran|kegiatan|piket|shift / /:id
GET|POST|PUT        /pondok/kalender / /:id
POST|GET            /pondok/absensi / /pondok/absensi/rekap
POST|GET|PUT        /pondok/nilai / /pondok/nilai/tahfidz|akhlak
POST|GET            /pondok/raport / /pondok/raport/:id/terbitkan
POST|GET|PUT        /pondok/jenis-tagihan / /:id
POST                /pondok/tagihan/generate
GET|PUT             /pondok/tagihan / /:id
POST                /pondok/tagihan/:id/beasiswa
POST                /pondok/pembiayaan
POST                /pondok/sinkron/dapodik|emis
GET                 /pondok/sinkron/log
```

## E-commerce

```
GET|POST|PUT  /shop/toko / /:id              # Owner toko (via pondok-app)
GET|POST|PUT  /shop/produk / /:id
PUT           /shop/produk/:id/stok
GET|PUT       /shop/pesanan / /:id/status
GET           /shop/laporan/penjualan
# OPOP marketplace (lintas pondok)
GET           /opop/toko
GET           /opop/produk
GET           /opop/toko/:slug
POST          /opop/pesanan                  # B2B: pondok pesan dari pondok lain
```

## Webhook

```
POST /webhook/midtrans                       # Verifikasi signature SHA512
```
