import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';

class SantriListPage extends StatelessWidget {
  const SantriListPage({super.key});

  static const _data = [
    {'nis': '2025-001', 'nama': 'Muhammad Faqih', 'kelas': 'MTS 7A', 'asrama': 'Al-Farabi', 'status': 'AKTIF'},
    {'nis': '2025-002', 'nama': 'Ahmad Habibi', 'kelas': 'MTS 7A', 'asrama': 'Al-Farabi', 'status': 'AKTIF'},
    {'nis': '2025-003', 'nama': 'Ridwan Kamil Jr', 'kelas': 'MA 10B', 'asrama': 'Ibnu Sina', 'status': 'AKTIF'},
    {'nis': '2024-089', 'nama': 'Zainul Arifin', 'kelas': 'MA 11A', 'asrama': 'Al-Ghazali', 'status': 'CUTI'},
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text(AppStrings.santri),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 16),
            child: ElevatedButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.add, size: 18),
              label: const Text(AppStrings.tambahSantri),
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.white,
                foregroundColor: AppColors.primary,
                minimumSize: const Size(0, 36),
              ),
            ),
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            TextField(
              decoration: const InputDecoration(
                hintText: AppStrings.cariSantri,
                prefixIcon: Icon(Icons.search),
              ),
            ),
            const SizedBox(height: 16),
            Expanded(
              child: ListView.separated(
                itemCount: _data.length,
                separatorBuilder: (_, __) => const SizedBox(height: 8),
                itemBuilder: (context, i) {
                  final s = _data[i];
                  final isAktif = s['status'] == 'AKTIF';
                  return Card(
                    child: ListTile(
                      leading: CircleAvatar(
                        backgroundColor: AppColors.primary.withOpacity(0.1),
                        child: Text(
                          s['nama']!.substring(0, 1),
                          style: const TextStyle(
                              color: AppColors.primary,
                              fontWeight: FontWeight.bold),
                        ),
                      ),
                      title: Text(s['nama']!,
                          style: const TextStyle(fontWeight: FontWeight.w600)),
                      subtitle: Text('${s['nis']}  •  ${s['kelas']}  •  ${s['asrama']}',
                          style: const TextStyle(fontSize: 12)),
                      trailing: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Container(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 8, vertical: 3),
                            decoration: BoxDecoration(
                              color: isAktif
                                  ? AppColors.statusHadir.withOpacity(0.1)
                                  : AppColors.warning.withOpacity(0.1),
                              borderRadius: BorderRadius.circular(4),
                            ),
                            child: Text(
                              s['status']!,
                              style: TextStyle(
                                  color: isAktif
                                      ? AppColors.statusHadir
                                      : AppColors.warning,
                                  fontSize: 11,
                                  fontWeight: FontWeight.w600),
                            ),
                          ),
                          const SizedBox(width: 8),
                          const Icon(Icons.chevron_right,
                              color: AppColors.textHint),
                        ],
                      ),
                      onTap: () {},
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
