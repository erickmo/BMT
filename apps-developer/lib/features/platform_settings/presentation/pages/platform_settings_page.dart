import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';

class PlatformSettingsPage extends StatelessWidget {
  const PlatformSettingsPage({super.key});

  static const _settings = [
    {'kunci': 'pecahan_uang.sumber', 'nilai': 'DB', 'tipe': 'string', 'rahasia': false},
    {'kunci': 'platform.midtrans_env', 'nilai': 'sandbox', 'tipe': 'string', 'rahasia': false},
    {'kunci': 'platform.maintenance_mode', 'nilai': 'false', 'tipe': 'bool', 'rahasia': false},
    {'kunci': 'platform.min_app_version.nasabah', 'nilai': '2.0.0', 'tipe': 'string', 'rahasia': false},
    {'kunci': 'platform.min_app_version.teller', 'nilai': '1.5.0', 'tipe': 'string', 'rahasia': false},
    {'kunci': 'platform.rate_limit_rpm', 'nilai': '300', 'tipe': 'int', 'rahasia': false},
    {'kunci': 'platform.midtrans_server_key', 'nilai': '•••••••••••••', 'tipe': 'string', 'rahasia': true},
    {'kunci': 'monetisasi.komisi_opop_persen', 'nilai': '2.5', 'tipe': 'float', 'rahasia': false},
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Platform Settings'),
        automaticallyImplyLeading: false,
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 16),
            child: ElevatedButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.add, size: 18),
              label: const Text('Tambah Setting'),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.accent,
                foregroundColor: Colors.black87,
                minimumSize: const Size(0, 36),
              ),
            ),
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Card(
          child: SingleChildScrollView(
            child: DataTable(
              columnSpacing: 20,
              columns: const [
                DataColumn(label: Text('Kunci')),
                DataColumn(label: Text('Nilai')),
                DataColumn(label: Text('Tipe')),
                DataColumn(label: Text('Rahasia')),
                DataColumn(label: Text('Aksi')),
              ],
              rows: _settings.map((s) {
                final isRahasia = s['rahasia'] as bool;
                return DataRow(cells: [
                  DataCell(
                    SelectableText(
                      s['kunci'] as String,
                      style: const TextStyle(
                          fontFamily: 'monospace',
                          fontSize: 12,
                          color: AppColors.accent),
                    ),
                  ),
                  DataCell(
                    isRahasia
                        ? const Text('••••••••••',
                            style: TextStyle(color: AppColors.textHint))
                        : Text(
                            s['nilai'] as String,
                            style: const TextStyle(fontFamily: 'monospace'),
                          ),
                  ),
                  DataCell(
                    Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 6, vertical: 2),
                      decoration: BoxDecoration(
                        color: AppColors.surfaceVariant,
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: Text(
                        s['tipe'] as String,
                        style: const TextStyle(
                            fontSize: 11,
                            fontFamily: 'monospace',
                            color: AppColors.textSecondary),
                      ),
                    ),
                  ),
                  DataCell(
                    isRahasia
                        ? const Icon(Icons.lock, size: 16, color: AppColors.warning)
                        : const Icon(Icons.lock_open,
                            size: 16, color: AppColors.textHint),
                  ),
                  DataCell(Row(
                    children: [
                      IconButton(
                        icon: const Icon(Icons.edit_outlined,
                            size: 18, color: AppColors.accent),
                        onPressed: () {},
                      ),
                    ],
                  )),
                ]);
              }).toList(),
            ),
          ),
        ),
      ),
    );
  }
}
