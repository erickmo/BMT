# 15 — Background Workers (asynq)

> **Terakhir diperbarui:** 20 Maret 2026

Semua jadwal dari **settings DB** — tidak ada cron hardcode di kode.

## CBS
| Worker | Jadwal (dari settings) | Scope |
|--------|----------------------|-------|
| `AutodebetHarian` | `autodebet.jam_eksekusi` (07:00) | Per BMT |
| `AutodebetBulanan` | Tgl 1, 07:30 | Per BMT |
| `GenerateJadwalAutodebet` | Tgl 25, 08:00 | Per BMT |
| `UpdateKolektibilitas` | 00:05 WIB daily | Per BMT |
| `DistribusiBagiHasil` | Akhir bulan | Per BMT |
| `ClosingBulanan` | Akhir bulan | Per BMT |
| `HitungZakat` | Akhir tahun (dari settings) | Per BMT |
| `ReminderAngsuran` | H-N dari `notifikasi.reminder_angsuran_hari` | Per BMT |

## Pondok
| Worker | Jadwal | Scope |
|--------|--------|-------|
| `GenerateTagihanSPP` | Tgl 25, 08:30 | Per BMT |
| `ReminderSPP` | H-N dari `notifikasi.reminder_spp_hari` | Per BMT |
| `PayrollBulanan` | Tgl dari `sdm.tanggal_gajian`, 06:00 | Per BMT |
| `GenerateRaport` | Akhir semester (dari settings) | Per BMT |
| `ReminderPerpus` | H-1 kembali | Per BMT |
| `SinkronDAPODIK` | Dari `integrasi.sinkron_jadwal` | Per BMT |
| `SinkronEMIS` | Dari `integrasi.sinkron_jadwal` | Per BMT |

## Platform
| Worker | Jadwal | Scope |
|--------|--------|-------|
| `CekMidtransPending` | Setiap 15 menit | Platform |
| `HitungUsageBulanan` | Tgl 1, 06:00 | Platform |
| `HitungKomisiOPOP` | Akhir bulan | Platform |
| `GenerateSaaSInvoice` | Tgl dari settings | Platform |
| `CekSaaSExpiry` | 08:00 WIB daily | Platform |
| `CekListingExpiry` | 08:00 WIB daily | Platform |
| `ReminderListingExpiry` | H-7 expiry | Platform |
| `CekKontrakExpiry` | 08:00 WIB daily | Platform |
| `ExpiredKartuNFC` | 08:00 WIB daily | Platform |
| `DepresiasiAset` | 1 Januari, 07:00 | Per BMT |
| `AnalyticsHarian` | 23:00 WIB daily | Platform |
| `FraudDetection` | **Real-time** per transaksi | Platform |
| `BackupDatabase` | 02:00 WIB daily | Platform |

## Offline Sync
| Worker | Trigger | Deskripsi |
|--------|---------|-----------|
| `SyncOfflineTransaksi` | Device online | Proses antrian transaksi teller offline |
| `SyncOfflineAbsensi` | Device online | Proses antrian absensi offline |
| `ResolveOfflineConflict` | Setelah sync | Flag konflik ke manajer |
