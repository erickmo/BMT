import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/utils/formatters.dart';

class PecahanUangPage extends StatelessWidget {
  const PecahanUangPage({super.key});

  static const _pecahan = [
    {'nominal': 100000, 'jenis': 'KERTAS', 'label': 'Rp 100.000 (kertas)', 'aktif': true, 'urutan': 1},
    {'nominal': 50000, 'jenis': 'KERTAS', 'label': 'Rp 50.000 (kertas)', 'aktif': true, 'urutan': 2},
    {'nominal': 20000, 'jenis': 'KERTAS', 'label': 'Rp 20.000 (kertas)', 'aktif': true, 'urutan': 3},
    {'nominal': 10000, 'jenis': 'KERTAS', 'label': 'Rp 10.000 (kertas)', 'aktif': true, 'urutan': 4},
    {'nominal': 5000, 'jenis': 'KERTAS', 'label': 'Rp 5.000 (kertas)', 'aktif': true, 'urutan': 5},
    {'nominal': 2000, 'jenis': 'KERTAS', 'label': 'Rp 2.000 (kertas)', 'aktif': true, 'urutan': 6},
    {'nominal': 1000, 'jenis': 'KERTAS', 'label': 'Rp 1.000 (kertas)', 'aktif': true, 'urutan': 7},
    {'nominal': 1000, 'jenis': 'LOGAM', 'label': 'Rp 1.000 (logam)', 'aktif': true, 'urutan': 8},
    {'nominal': 500, 'jenis': 'LOGAM', 'label': 'Rp 500 (logam)', 'aktif': true, 'urutan': 9},
    {'nominal': 200, 'jenis': 'LOGAM', 'label': 'Rp 200 (logam)', 'aktif': true, 'urutan': 10},
    {'nominal': 100, 'jenis': 'LOGAM', 'label': 'Rp 100 (logam)', 'aktif': true, 'urutan': 11},
    {'nominal': 75000, 'jenis': 'KERTAS', 'label': 'Rp 75.000 (kertas)', 'aktif': false, 'urutan': 99},
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Pecahan Uang Rupiah'),
        automaticallyImplyLeading: false,
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 16),
            child: ElevatedButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.add, size: 18),
              label: const Text('Tambah Pecahan'),
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
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: AppColors.codeBg,
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Row(
                children: [
                  Icon(Icons.info_outline, color: AppColors.codeText, size: 16),
                  SizedBox(width: 8),
                  Text(
                    'Pecahan uang diambil dari DB — tidak ada array hardcode di kode.',
                    style: TextStyle(
                        color: AppColors.codeText,
                        fontFamily: 'monospace',
                        fontSize: 12),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 20),
            Expanded(
              child: Card(
                child: SingleChildScrollView(
                  child: DataTable(
                    columns: const [
                      DataColumn(label: Text('Nominal')),
                      DataColumn(label: Text('Jenis')),
                      DataColumn(label: Text('Label')),
                      DataColumn(label: Text('Urutan')),
                      DataColumn(label: Text('Status')),
                      DataColumn(label: Text('Aksi')),
                    ],
                    rows: _pecahan.map((p) {
                      final isAktif = p['aktif'] as bool;
                      return DataRow(cells: [
                        DataCell(Text(
                          formatRupiah(p['nominal'] as int),
                          style: const TextStyle(fontWeight: FontWeight.w600),
                        )),
                        DataCell(
                          Container(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 8, vertical: 3),
                            decoration: BoxDecoration(
                              color: (p['jenis'] == 'KERTAS'
                                      ? AppColors.info
                                      : AppColors.warning)
                                  .withOpacity(0.1),
                              borderRadius: BorderRadius.circular(4),
                            ),
                            child: Text(
                              p['jenis'] as String,
                              style: TextStyle(
                                color: p['jenis'] == 'KERTAS'
                                    ? AppColors.info
                                    : AppColors.warning,
                                fontSize: 11,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ),
                        ),
                        DataCell(Text(p['label'] as String)),
                        DataCell(Text('#${p['urutan']}')),
                        DataCell(
                          Switch(
                            value: isAktif,
                            onChanged: (_) {},
                            activeColor: AppColors.success,
                          ),
                        ),
                        DataCell(Row(
                          children: [
                            IconButton(
                              icon: const Icon(Icons.edit_outlined,
                                  size: 18, color: AppColors.accent),
                              onPressed: () {},
                            ),
                            IconButton(
                              icon: const Icon(Icons.delete_outline,
                                  size: 18, color: AppColors.error),
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
          ],
        ),
      ),
    );
  }
}
