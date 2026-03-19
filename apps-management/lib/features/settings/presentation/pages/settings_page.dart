import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';

class SettingsPage extends StatelessWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text(AppStrings.settingsBMT),
        automaticallyImplyLeading: false,
      ),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Settings categories sidebar
            SizedBox(
              width: 240,
              child: Card(
                child: Column(
                  children: [
                    _SettingsCategory(
                        icon: Icons.access_time,
                        label: 'Operasional',
                        selected: true),
                    _SettingsCategory(
                        icon: Icons.schedule,
                        label: 'Autodebet',
                        selected: false),
                    _SettingsCategory(
                        icon: Icons.approval,
                        label: 'Approval',
                        selected: false),
                    _SettingsCategory(
                        icon: Icons.notifications,
                        label: 'Notifikasi',
                        selected: false),
                    _SettingsCategory(
                        icon: Icons.payment,
                        label: 'Midtrans',
                        selected: false),
                  ],
                ),
              ),
            ),
            const SizedBox(width: 20),
            // Settings form
            Expanded(
              child: Card(
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Pengaturan Operasional',
                        style: Theme.of(context)
                            .textTheme
                            .titleMedium
                            ?.copyWith(fontWeight: FontWeight.bold),
                      ),
                      const Divider(height: 32),
                      _SettingField(
                        kunci: 'operasional.jam_buka',
                        label: 'Jam Buka',
                        nilai: '08:00',
                        hint: 'Format HH:mm',
                      ),
                      const SizedBox(height: 16),
                      _SettingField(
                        kunci: 'operasional.jam_tutup',
                        label: 'Jam Tutup',
                        nilai: '16:00',
                        hint: 'Format HH:mm',
                      ),
                      const SizedBox(height: 16),
                      _SettingField(
                        kunci: 'operasional.zona_waktu',
                        label: 'Zona Waktu',
                        nilai: 'Asia/Jakarta',
                        hint: 'Contoh: Asia/Jakarta',
                      ),
                      const SizedBox(height: 16),
                      _SettingField(
                        kunci: 'operasional.hari_kerja',
                        label: 'Hari Kerja',
                        nilai: '[1,2,3,4,5]',
                        hint: 'Array hari (1=Senin, 7=Minggu)',
                      ),
                      const SizedBox(height: 32),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.end,
                        children: [
                          OutlinedButton(
                            onPressed: () {},
                            child: const Text('Reset'),
                          ),
                          const SizedBox(width: 12),
                          ElevatedButton(
                            onPressed: () {},
                            child: const Text(AppStrings.save),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _SettingsCategory extends StatelessWidget {
  final IconData icon;
  final String label;
  final bool selected;

  const _SettingsCategory({
    required this.icon,
    required this.label,
    required this.selected,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Icon(
        icon,
        color: selected ? AppColors.primary : AppColors.textSecondary,
        size: 20,
      ),
      title: Text(
        label,
        style: TextStyle(
          color: selected ? AppColors.primary : AppColors.textPrimary,
          fontWeight: selected ? FontWeight.w600 : FontWeight.normal,
          fontSize: 14,
        ),
      ),
      selected: selected,
      selectedTileColor: AppColors.primary.withOpacity(0.08),
      onTap: () {},
    );
  }
}

class _SettingField extends StatelessWidget {
  final String kunci;
  final String label;
  final String nilai;
  final String hint;

  const _SettingField({
    required this.kunci,
    required this.label,
    required this.nilai,
    required this.hint,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(
          width: 260,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label,
                style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary),
              ),
              Text(
                kunci,
                style: const TextStyle(
                    fontSize: 11,
                    fontFamily: 'monospace',
                    color: AppColors.textSecondary),
              ),
            ],
          ),
        ),
        Expanded(
          child: TextFormField(
            initialValue: nilai,
            decoration: InputDecoration(hintText: hint),
          ),
        ),
      ],
    );
  }
}
