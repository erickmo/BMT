# Domain: Background Workers

Semua worker menggunakan **asynq** (Redis-backed). Jadwal dalam WIB kecuali disebutkan lain.

## Workers Per BMT

| Worker | Jadwal | Deskripsi |
|--------|--------|-----------|
| `AutodebetHarian` | 07:00 daily | Eksekusi jadwal autodebet (tanggal dari `rekening_autodebet_config`) |
| `AutodebetBulanan` | Tgl 1, 07:30 | Biaya admin rekening bulanan |
| `GenerateJadwalAutodebet` | Tgl 25, 08:00 | Generate jadwal autodebet bulan depan |
| `GenerateTagihanSPP` | Tgl 25, 08:30 | Generate tagihan SPP periode berikutnya |
| `UpdateKolektibilitas` | 00:05 daily | Hari tunggak → update kolektibilitas pembiayaan |
| `DistribusiBagiHasil` | Akhir bulan | Posting bagi hasil deposito |
| `ReminderAngsuran` | H-N (dari settings) | Email reminder angsuran pembiayaan |
| `ReminderSPP` | H-N (dari settings) | Email reminder SPP ke wali santri |
| `ClosingBulanan` | Akhir bulan | Jurnal penutup + laporan otomatis |
| `GenerateRaport` | Akhir semester (dari settings) | Generate raport digital semua santri |
| `GenerateSlipGaji` | Tgl 25, 09:00 | Generate slip gaji semua pegawai aktif |
| `EksekusiPayroll` | Tgl 1, 08:00 | Transfer gaji ke rekening BMT pegawai |
| `SinkronEMIS` | Sabtu, 02:00 | Sinkronisasi data santri dengan EMIS Kemenag |
| `SinkronDAPODIK` | Tgl 1, 03:00 | Sinkronisasi data siswa dengan DAPODIK |
| `ReminderPeminjamanBuku` | 08:00 daily | Reminder buku perpustakaan hampir jatuh tempo |
| `KirimNotifikasi` | Setiap 1 menit | Proses antrian notifikasi (FCM/WA/SMS/Email) |

## Workers Platform (Developer)

| Worker | Jadwal | Deskripsi |
|--------|--------|-----------|
| `CekMidtransPending` | Setiap 15 menit | Poll status transaksi Midtrans PENDING |
| `HitungUsageBulanan` | Tgl 1, 06:00 | Rekap usage_log bulan lalu |
| `CekKontrakExpiry` | 08:00 daily | Suspend BMT dengan kontrak expired |
| `ExpiredKartuNFC` | 08:00 daily | Set kartu NFC expired |
| `BackupDatabase` | 02:00 daily | Dump PostgreSQL → MinIO |
| `HitungKomisiOPOP` | Tgl 1, 07:00 | Rekap komisi OPOP bulan lalu → usage_log |
| `DepresiasiAset` | Tgl 1, tahunan | Hitung & posting depresiasi aset tetap |
| `GenerateSnapshotAnalytics` | 00:05 daily | Hitung & simpan snapshot metrik harian |
| `AnalyticsHarian` | 23:00 daily | Snapshot analytics OPOP & kesiswaan |
| `FraudDetection` | Real-time (event) | Cek setiap transaksi masuk ke fraud rules |
| `CleanupAuditLog` | Minggu, 03:00 | Hapus audit log melebihi retensi (dari settings) |
| `CleanupSesiExpired` | Setiap jam | Hapus sesi JWT yang sudah expired |
