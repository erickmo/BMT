# 14 — API Endpoint

> **Terakhir diperbarui:** 20 Maret 2026

## Auth Header per Grup

| Grup | Header Wajib |
|------|-------------|
| `/dev/*` | `Developer-Token: <secret>` |
| `/platform/*` | `Authorization: Bearer <jwt>` + `X-BMT-ID` |
| `/teller/*` | `Authorization` + `X-BMT-ID` + `X-Cabang-ID` |
| `/api/*` | `Authorization` + `X-BMT-ID` + `X-Cabang-ID` |
| `/pondok/*` | `Authorization (jwt_pondok)` + `X-BMT-ID` + `X-Cabang-ID` |
| `/nasabah/*` | `Authorization (jwt_nasabah)` |
| `/nfc/*` | `X-API-Key: <merchant_api_key>` |
| `/listing/daftar` | Publik (tanpa auth) |
| `/webhook/*` | Tanpa auth — verifikasi via signature |

---

## Developer
```
POST|GET|PUT     /dev/bmt / /:id / /:id/status
POST             /dev/bmt/:id/langganan
PUT              /dev/bmt/:id/langganan/:id    # Upgrade/downgrade tier
POST|DELETE      /dev/bmt/:id/addon / /:id
POST|GET|PUT     /dev/bmt/:id/cabang / /:id
POST             /dev/bmt/:id/seed

GET|POST|PUT     /dev/saas/paket / /:id
GET|POST|PUT     /dev/saas/addon / /:id

GET|POST|PUT|DELETE /dev/pecahan-uang / /:id
GET|POST|PUT     /dev/tarif-template / /:id
GET|PUT          /dev/platform-settings / /:kunci

GET              /dev/invoice / /:bmt_id
POST             /dev/invoice/generate
GET              /dev/revenue                  # MRR, ARR, churn

GET              /dev/listing/pendaftar        # Antrian approve
POST             /dev/listing/pendaftar/:id/approve|reject
GET|POST|PUT     /dev/listing/kategori / /:id
GET              /dev/listing/aktif
PUT              /dev/listing/:id/suspend
GET              /dev/listing/revenue

GET              /dev/usage-log / /:bmt_id
GET              /dev/metrics / /dev/health/detail
POST             /dev/maintenance
```

## Auth
```
POST  /auth/staf/login|refresh|logout
POST  /auth/nasabah/login|refresh
POST  /auth/pondok/login|refresh
POST  /auth/merchant/login|refresh
POST  /auth/listing/login|refresh
GET   /auth/sesi
DELETE /auth/sesi/:id | /auth/sesi  # Logout device / semua
```

## Platform (Management BMT)
```
GET|PUT          /platform/settings / /:kunci
GET|POST|PUT     /platform/jenis-rekening / /:id
GET|POST|PUT     /platform/produk/simpanan|pembiayaan / /:id
GET|POST|PUT     /platform/pengguna / /:id
GET|POST|PUT     /platform/cabang / /:id
GET|POST|PUT     /platform/merchant / /:id
GET|POST|PUT     /platform/terminal-kiosk / /:id
GET|POST         /platform/pengguna-pondok / /:id
GET|POST|PUT|DELETE /platform/ttd / /:id
GET              /platform/laporan/konsolidasi|rat
GET              /platform/usage
GET              /platform/saas/langganan      # Info paket aktif BMT ini
GET|POST         /platform/chat/room / /:id/pesan
```

## Teller
```
POST  /teller/sesi/buka|tutup
GET   /teller/sesi/aktif
GET   /teller/nasabah/cari
GET   /teller/rekening/:nomor
POST  /teller/rekening/:id/setor|tarik         # X-Idempotency-Key wajib
POST  /teller/pembiayaan/:id/angsuran
POST  /teller/spp/:id/bayar
GET   /teller/dokumen/:id/cetak
```

## Form Workflow
```
POST|GET  /form / /:id
PUT       /form/:id                    # Hanya DRAFT
POST      /form/:id/ajukan|setujui|tolak
GET       /form/:id/riwayat
```

## Management API
```
GET  /api/nasabah / /:id               # View only
GET  /api/rekening / /:id / /:id/transaksi|tunggakan
GET|POST|PUT /api/rekening/:id/autodebet-config
GET  /api/autodebet/jadwal|tunggakan
POST /api/pembiayaan
GET  /api/pembiayaan / /:id / /:id/jadwal
PUT  /api/pembiayaan/:id/analisis
POST /api/pembiayaan/:id/pencairan|simulasi
POST|GET /finance/jurnal|transaksi
GET  /api/laporan/neraca|shu|arus-kas|kolektibilitas|bagi-hasil-deposito|transaksi-harian|rat
GET|POST|PUT /api/laporan/template / /:id
POST /api/laporan/generate
GET  /api/dokumen/kartu-anggota/:id|buku-tabungan/:id
POST /api/dokumen/:id/upload-ttd|kirim-email
GET  /api/dokumen/:id/download|versi
GET  /ws/dashboard | /sse/dashboard    # Real-time
GET|PUT /api/fraud/alert / /:id
```

## Nasabah App
```
GET  /nasabah/profil|rekening|pembiayaan
GET  /nasabah/santri                   # + kamar, hafalan progress
GET  /nasabah/rekening/:id/transaksi|tunggakan
POST /nasabah/rekening/:id/setor-online
GET  /nasabah/nfc/saldo|transaksi
POST /nasabah/nfc/topup
GET  /nasabah/spp/tagihan
POST /nasabah/spp/tagihan/:id/bayar
GET|POST /nasabah/donasi|donasi/bayar
GET  /nasabah/shop/toko|produk|pesanan
POST /nasabah/shop/pesanan / /:id/cancel
GET  /nasabah/dokumen/:id/download
GET  /nasabah/bulletin
GET  /nasabah/listing                  # ?kategori=&radius=
GET  /nasabah/listing/:id
POST /nasabah/listing/:id/ulasan
GET  /nasabah/listing/kategori
```

## NFC / Kiosk
```
POST /nfc/transaksi                    # X-Idempotency-Key wajib
GET  /nfc/ceksaldo/:uid                # IP whitelist, tanpa PIN
```

## Pondok
```
GET|POST|PUT|DELETE /pondok/santri / /:id
GET|POST|PUT|DELETE /pondok/pengajar|karyawan / /:id
GET|POST|PUT /pondok/kelas / /:id
GET|POST|PUT /pondok/asrama|kamar / /:id
PUT          /pondok/santri/:id/kamar
GET|POST|PUT /pondok/mapel|silabus|rpp / /:id
GET|POST|PUT /pondok/komponen-nilai / /:id
GET|POST|PUT /pondok/jadwal/pelajaran|kegiatan|piket|shift / /:id
GET|POST|PUT /pondok/kalender / /:id
POST|GET     /pondok/absensi / /pondok/absensi/rekap
POST|GET|PUT /pondok/nilai / /pondok/nilai/tahfidz|akhlak
GET|POST     /pondok/raport / /:id/terbitkan
GET|POST|PUT /pondok/jenis-tagihan / /:id
POST         /pondok/tagihan/generate
GET|PUT      /pondok/tagihan / /:id
POST         /pondok/tagihan/:id/beasiswa
POST         /pondok/pembiayaan
GET|POST|PUT /pondok/sdm/kontrak / /:id
GET          /pondok/sdm/slip-gaji / /:id
GET|POST|PUT|DELETE /pondok/perpus/buku / /:id
GET|POST|PUT /pondok/perpus/peminjaman / /:id
GET|POST|PUT /pondok/konsultasi / /:id/pesan
GET|POST|PUT /pondok/surat-izin / /:id/setujui|tolak
GET|POST     /pondok/health-record
GET|POST|PUT /pondok/alumni / /:id
GET|POST|PUT /pondok/aset / /:id/pinjam
GET|POST|PUT /pondok/ppdb/gelombang|pendaftar / /:id/terima|tolak
POST         /pondok/sinkron/dapodik|emis
GET|POST|PUT /pondok/toko|produk / /:id / /:id/stok
GET|PUT      /pondok/pesanan / /:id/status
GET|POST|PUT /pondok/santri/:id/portfolio|hafalan
GET|POST|PUT /pondok/ekstra / /:id
GET|POST|PUT /pondok/event / /:id
POST         /pondok/event/:id/registrasi
```

## Listing (Owner)
```
POST /listing/daftar                   # Publik — self-register
GET  /listing/status/:id               # Cek status pendaftaran
GET|PUT /listing/profil                # Setelah approved
POST    /listing/langganan             # Upgrade ke premium
GET     /listing/ulasan
```

## E-commerce
```
GET  /opop/toko / /:slug
GET  /opop/produk
POST /opop/pesanan                     # B2B
```

## Webhook
```
POST /webhook/midtrans                 # Verifikasi SHA512
```
