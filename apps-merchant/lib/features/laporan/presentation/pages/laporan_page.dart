import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';

class LaporanPage extends StatelessWidget {
  const LaporanPage({super.key});

  static const _transaksi = [
    {'waktu': '08:15', 'nasabah': 'Muhammad Faqih', 'nominal': 15000, 'status': 'BERHASIL'},
    {'waktu': '09:02', 'nasabah': 'Ahmad Habibi', 'nominal': 7500, 'status': 'BERHASIL'},
    {'waktu': '10:35', 'nasabah': 'Siti Aminah', 'nominal': 25000, 'status': 'BERHASIL'},
    {'waktu': '11:10', 'nasabah': 'Ridwan Kamil', 'nominal': 12000, 'status': 'GAGAL'},
    {'waktu': '13:45', 'nasabah': 'Dewi Lestari', 'nominal': 8000, 'status': 'BERHASIL'},
  ];

  @override
  Widget build(BuildContext context) {
    final totalBerhasil = _transaksi
        .where((t) => t['status'] == 'BERHASIL')
        .fold<int>(0, (sum, t) => sum + (t['nominal'] as int));
    final jumlahBerhasil =
        _transaksi.where((t) => t['status'] == 'BERHASIL').length;

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text(AppStrings.laporan),
        automaticallyImplyLeading: false,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Row(
              children: [
                Expanded(
                  child: Card(
                    child: Padding(
                      padding: const EdgeInsets.all(20),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(AppStrings.totalHariIni,
                              style: TextStyle(
                                  color: AppColors.textSecondary, fontSize: 13)),
                          const SizedBox(height: 4),
                          Text(
                            formatRupiah(totalBerhasil),
                            style: const TextStyle(
                                fontSize: 24,
                                fontWeight: FontWeight.bold,
                                color: AppColors.success),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Card(
                    child: Padding(
                      padding: const EdgeInsets.all(20),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(AppStrings.jumlahTransaksi,
                              style: TextStyle(
                                  color: AppColors.textSecondary, fontSize: 13)),
                          const SizedBox(height: 4),
                          Text(
                            '$jumlahBerhasil transaksi',
                            style: const TextStyle(
                                fontSize: 24,
                                fontWeight: FontWeight.bold,
                                color: AppColors.primary),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Text(
              'Riwayat Transaksi — ${formatTanggal(DateTime.now())}',
              style: const TextStyle(
                  fontWeight: FontWeight.bold, color: AppColors.textPrimary),
            ),
            const SizedBox(height: 8),
            Expanded(
              child: ListView.separated(
                itemCount: _transaksi.length,
                separatorBuilder: (_, __) => const SizedBox(height: 6),
                itemBuilder: (context, i) {
                  final t = _transaksi[i];
                  final isOk = t['status'] == 'BERHASIL';
                  return Card(
                    child: ListTile(
                      leading: Icon(
                        isOk ? Icons.check_circle : Icons.cancel,
                        color: isOk ? AppColors.success : AppColors.error,
                        size: 28,
                      ),
                      title: Row(
                        children: [
                          Expanded(child: Text(t['nasabah'] as String,
                              style: const TextStyle(fontWeight: FontWeight.w600))),
                          Text(
                            formatRupiah(t['nominal'] as int),
                            style: TextStyle(
                                fontWeight: FontWeight.bold,
                                color: isOk ? AppColors.textPrimary : AppColors.textHint),
                          ),
                        ],
                      ),
                      subtitle: Text(t['waktu'] as String,
                          style: const TextStyle(fontSize: 12)),
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
