import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';

class JadwalPage extends StatelessWidget {
  const JadwalPage({super.key});

  static const _jadwal = [
    {'hari': 'Senin', 'jam': '07:00 - 08:30', 'mapel': 'Matematika', 'kelas': 'MTS 7A', 'pengajar': 'Ust. Ahmad'},
    {'hari': 'Senin', 'jam': '08:30 - 10:00', 'mapel': 'Bahasa Arab', 'kelas': 'MTS 7A', 'pengajar': 'Ust. Hasan'},
    {'hari': 'Senin', 'jam': '10:15 - 11:45', 'mapel': 'Fiqih', 'kelas': 'MTS 7A', 'pengajar': 'Ust. Ibrahim'},
    {'hari': 'Selasa', 'jam': '07:00 - 08:30', 'mapel': 'IPA', 'kelas': 'MTS 7A', 'pengajar': 'Ust. Sholeh'},
    {'hari': 'Selasa', 'jam': '08:30 - 10:00', 'mapel': 'Bahasa Indonesia', 'kelas': 'MTS 7A', 'pengajar': 'Ust. Hendra'},
  ];

  static const _hariColors = {
    'Senin': Color(0xFF1565C0),
    'Selasa': Color(0xFF2E7D32),
    'Rabu': Color(0xFFE65100),
    'Kamis': Color(0xFF6A1B9A),
    'Jumat': Color(0xFF00838F),
    'Sabtu': Color(0xFF795548),
  };

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text(AppStrings.jadwalPelajaran),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {},
            tooltip: 'Tambah Jadwal',
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            // Filter row
            Row(
              children: [
                Expanded(
                  child: DropdownButtonFormField<String>(
                    value: 'MTS 7A',
                    decoration: const InputDecoration(labelText: 'Kelas'),
                    items: ['MTS 7A', 'MTS 7B', 'MA 10A', 'MA 10B']
                        .map((k) => DropdownMenuItem(value: k, child: Text(k)))
                        .toList(),
                    onChanged: (_) {},
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: DropdownButtonFormField<String>(
                    value: '2025/2026',
                    decoration: const InputDecoration(labelText: 'Tahun Ajaran'),
                    items: ['2024/2025', '2025/2026']
                        .map((t) => DropdownMenuItem(value: t, child: Text(t)))
                        .toList(),
                    onChanged: (_) {},
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Expanded(
              child: ListView.separated(
                itemCount: _jadwal.length,
                separatorBuilder: (_, __) => const SizedBox(height: 8),
                itemBuilder: (context, i) {
                  final j = _jadwal[i];
                  final hariColor =
                      _hariColors[j['hari']] ?? AppColors.primary;
                  return Card(
                    child: ListTile(
                      leading: Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 10, vertical: 8),
                        decoration: BoxDecoration(
                          color: hariColor.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(
                          j['hari']!.substring(0, 3),
                          style: TextStyle(
                            color: hariColor,
                            fontWeight: FontWeight.bold,
                            fontSize: 12,
                          ),
                        ),
                      ),
                      title: Text(j['mapel']!,
                          style: const TextStyle(fontWeight: FontWeight.w600)),
                      subtitle: Text(
                          '${j['jam']}  •  ${j['kelas']}  •  ${j['pengajar']}',
                          style: const TextStyle(fontSize: 12)),
                      trailing: IconButton(
                        icon: const Icon(Icons.edit_outlined,
                            size: 18, color: AppColors.textSecondary),
                        onPressed: () {},
                      ),
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
