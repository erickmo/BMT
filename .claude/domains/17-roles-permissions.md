# Domain: Role & Permission

## Tabel Role Lengkap

| Role | App | Scope | Kapabilitas |
|------|-----|-------|------------|
| `DEVELOPER` | apps-developer | Platform | `/dev/*`, pecahan uang, kontrak BMT, platform settings |
| `ADMIN_BMT` | apps-management | Seluruh BMT | Settings, jenis rekening, produk, pengguna |
| `MANAJER_BMT` | apps-management | Seluruh BMT | Laporan konsolidasi, approval besar |
| `AUDITOR_BMT` | apps-management | Seluruh BMT | Read-only |
| `MANAJER_CABANG` | apps-management | 1 cabang | Approval form, laporan, autodebet config |
| `KOMITE` | apps-management | 1 cabang | Approval pembiayaan |
| `ACCOUNT_OFFICER` | apps-management | 1 cabang | Pengajuan & monitoring pembiayaan |
| `FINANCE` | apps-management | 1 cabang | Jurnal manual, biaya operasional, approve slip gaji |
| `TELLER` | apps-teller | 1 cabang | Transaksi tunai, sesi kas (pecahan dari DB) |
| `NASABAH` | apps-nasabah | Data milik sendiri | E-banking + santri + belanja OPOP |
| `KASIR_MERCHANT` | apps-merchant | 1 merchant | Transaksi NFC kasir |
| `OWNER_MERCHANT` | apps-merchant | 1 merchant | Laporan penjualan, dashboard toko |
| `ADMIN_PONDOK` | apps-pondok | 1 cabang | Semua fitur pondok |
| `OPERATOR_PONDOK` | apps-pondok | 1 cabang | Input data santri, absensi, nilai |
| `BENDAHARA_PONDOK` | apps-pondok | 1 cabang | Tagihan, beasiswa, laporan keuangan pondok |
| `PETUGAS_UKS` | apps-pondok | 1 cabang | Input kunjungan UKS, health record |
| `PUSTAKAWAN` | apps-pondok | 1 cabang | Kelola buku, peminjaman perpustakaan |
| `PETUGAS_PPDB` | apps-pondok | 1 cabang | Proses pendaftaran santri baru, seleksi |
| `BK` | apps-pondok | 1 cabang | Konsultasi online, surat izin santri |
| `ALUMNI` | apps-nasabah | Data milik sendiri | Profil alumni, jaringan, belanja OPOP |

---

## Aturan Akses Khusus

- **Teller:** Semua tombol transaksi **disabled** tanpa sesi kas aktif
- **Kiosk (ceksaldo):** IP whitelist terminal — tidak ada login, tidak ada PIN
- **NFC transaksi:** Wajib `X-Idempotency-Key` header
- **ADMIN_PONDOK:** Satu-satunya yang bisa set beasiswa pembiayaan (`PUT /api/pembiayaan/:id/beasiswa`)
- **Developer:** Akses via `Developer-Token` header (bukan JWT biasa)
