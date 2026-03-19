import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';

class DashboardPage extends StatelessWidget {
  const DashboardPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text(AppStrings.navDashboard),
        automaticallyImplyLeading: false,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Ringkasan BMT',
              style: Theme.of(context).textTheme.titleLarge?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
            ),
            const SizedBox(height: 20),
            // KPI Cards
            GridView.count(
              crossAxisCount:
                  MediaQuery.of(context).size.width > 1200 ? 4 : 2,
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              crossAxisSpacing: 16,
              mainAxisSpacing: 16,
              childAspectRatio: 1.8,
              children: const [
                _KpiCard(
                  label: AppStrings.totalNasabah,
                  value: '1.248',
                  icon: Icons.people,
                  color: AppColors.primary,
                  trend: '+12 bulan ini',
                ),
                _KpiCard(
                  label: AppStrings.totalDPK,
                  value: 'Rp 4,2M',
                  icon: Icons.account_balance_wallet,
                  color: AppColors.success,
                  trend: '+5.2% bulan ini',
                ),
                _KpiCard(
                  label: AppStrings.totalPembiayaan,
                  value: 'Rp 2,8M',
                  icon: Icons.receipt_long,
                  color: AppColors.secondary,
                  trend: '+3.1% bulan ini',
                ),
                _KpiCard(
                  label: AppStrings.npfRatio,
                  value: '2.4%',
                  icon: Icons.warning_amber,
                  color: AppColors.warning,
                  trend: 'Maret 2025',
                ),
              ],
            ),
            const SizedBox(height: 32),
            Text(
              'Aktivitas Hari Ini',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(child: _ActivityCard(
                  title: 'Form Menunggu Approval',
                  value: '5',
                  icon: Icons.pending_actions,
                  color: AppColors.warning,
                )),
                const SizedBox(width: 16),
                Expanded(child: _ActivityCard(
                  title: 'Transaksi Hari Ini',
                  value: '128',
                  icon: Icons.swap_horiz,
                  color: AppColors.primary,
                )),
                const SizedBox(width: 16),
                Expanded(child: _ActivityCard(
                  title: 'Autodebet Berhasil',
                  value: '43',
                  icon: Icons.check_circle_outline,
                  color: AppColors.success,
                )),
                const SizedBox(width: 16),
                Expanded(child: _ActivityCard(
                  title: 'Tunggakan Baru',
                  value: '3',
                  icon: Icons.error_outline,
                  color: AppColors.error,
                )),
              ],
            ),
            const SizedBox(height: 32),
            Text(
              'Form Pengajuan Terbaru',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
            ),
            const SizedBox(height: 12),
            _RecentFormsCard(),
          ],
        ),
      ),
    );
  }
}

class _KpiCard extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;
  final Color color;
  final String trend;

  const _KpiCard({
    required this.label,
    required this.value,
    required this.icon,
    required this.color,
    required this.trend,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: color.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Icon(icon, color: color, size: 20),
                ),
                const Spacer(),
                Text(
                  trend,
                  style: TextStyle(
                    color: AppColors.textSecondary,
                    fontSize: 11,
                  ),
                ),
              ],
            ),
            const Spacer(),
            Text(
              value,
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: AppColors.textPrimary,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              label,
              style: TextStyle(
                color: AppColors.textSecondary,
                fontSize: 12,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ActivityCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;
  final Color color;

  const _ActivityCard({
    required this.title,
    required this.value,
    required this.icon,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Icon(icon, color: color, size: 32),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    value,
                    style: const TextStyle(
                      fontSize: 22,
                      fontWeight: FontWeight.bold,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  Text(
                    title,
                    style: const TextStyle(
                      fontSize: 12,
                      color: AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _RecentFormsCard extends StatelessWidget {
  final _forms = const [
    {'jenis': 'Daftar Nasabah', 'nama': 'Ahmad Fauzi', 'waktu': '10 menit lalu', 'status': 'MENUNGGU'},
    {'jenis': 'Buka Rekening', 'nama': 'Siti Rahayu', 'waktu': '25 menit lalu', 'status': 'MENUNGGU'},
    {'jenis': 'Pembiayaan', 'nama': 'Budi Santoso', 'waktu': '1 jam lalu', 'status': 'DISETUJUI'},
    {'jenis': 'Tutup Rekening', 'nama': 'Dewi Lestari', 'waktu': '2 jam lalu', 'status': 'DITOLAK'},
  ];

  const _RecentFormsCard();

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: _forms.map((form) {
            final statusColor = form['status'] == 'MENUNGGU'
                ? AppColors.warning
                : form['status'] == 'DISETUJUI'
                    ? AppColors.success
                    : AppColors.error;
            return ListTile(
              leading: CircleAvatar(
                backgroundColor: AppColors.surfaceVariant,
                child: const Icon(Icons.description_outlined,
                    color: AppColors.primary),
              ),
              title: Text(form['jenis']!,
                  style: const TextStyle(fontWeight: FontWeight.w600)),
              subtitle: Text(form['nama']!),
              trailing: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Container(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                    decoration: BoxDecoration(
                      color: statusColor.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      form['status']!,
                      style: TextStyle(
                          color: statusColor,
                          fontSize: 11,
                          fontWeight: FontWeight.w600),
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(form['waktu']!,
                      style: const TextStyle(
                          fontSize: 11, color: AppColors.textHint)),
                ],
              ),
            );
          }).toList(),
        ),
      ),
    );
  }
}
