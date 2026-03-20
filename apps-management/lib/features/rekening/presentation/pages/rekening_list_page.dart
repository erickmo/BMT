import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';

class RekeningListPage extends StatelessWidget {
  const RekeningListPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text(AppStrings.rekening),
        automaticallyImplyLeading: false,
      ),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Expanded(
                  flex: 3,
                  child: TextField(
                    decoration: InputDecoration(
                      hintText: 'Cari rekening...',
                      prefixIcon: Icon(Icons.search),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                OutlinedButton.icon(
                  onPressed: () {},
                  icon: const Icon(Icons.filter_list, size: 16),
                  label: const Text('Filter'),
                  style: OutlinedButton.styleFrom(
                    minimumSize: const Size(0, 44),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 20),
            Expanded(
              child: Card(
                child: SingleChildScrollView(
                  child: DataTable(
                    columnSpacing: 24,
                    columns: const [
                      DataColumn(label: Text('No. Rekening')),
                      DataColumn(label: Text('Nasabah')),
                      DataColumn(label: Text('Jenis')),
                      DataColumn(label: Text('Saldo')),
                      DataColumn(label: Text('Status')),
                      DataColumn(label: Text('Autodebet')),
                      DataColumn(label: Text('Aksi')),
                    ],
                    rows: [
                      _buildRow(
                        'ANNUR-KDR-SU-00000001',
                        'Ahmad Fauzi',
                        'Simpanan Sukarela',
                        1500000,
                        'AKTIF',
                        true,
                      ),
                      _buildRow(
                        'ANNUR-KDR-SW-00000001',
                        'Ahmad Fauzi',
                        'Simpanan Wajib',
                        500000,
                        'AKTIF',
                        true,
                      ),
                      _buildRow(
                        'ANNUR-KDR-SU-00000002',
                        'Siti Rahayu',
                        'Simpanan Sukarela',
                        2800000,
                        'AKTIF',
                        false,
                      ),
                      _buildRow(
                        'ANNUR-KDR-SU-00000003',
                        'Budi Santoso',
                        'Simpanan Sukarela',
                        0,
                        'BLOKIR',
                        false,
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

  DataRow _buildRow(
    String nomor,
    String nama,
    String jenis,
    int saldo,
    String status,
    bool autodebet,
  ) {
    final isAktif = status == 'AKTIF';
    return DataRow(cells: [
      DataCell(Text(nomor,
          style: const TextStyle(
              fontFamily: 'monospace',
              fontWeight: FontWeight.w600,
              color: AppColors.primary,
              fontSize: 12))),
      DataCell(Text(nama)),
      DataCell(Text(jenis,
          style: const TextStyle(fontSize: 12))),
      DataCell(Text(
        formatRupiah(saldo),
        style: const TextStyle(fontWeight: FontWeight.w600),
      )),
      DataCell(
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
          decoration: BoxDecoration(
            color: isAktif
                ? AppColors.statusAktif.withOpacity(0.1)
                : AppColors.statusBlokir.withOpacity(0.1),
            borderRadius: BorderRadius.circular(4),
          ),
          child: Text(
            status,
            style: TextStyle(
              color: isAktif ? AppColors.statusAktif : AppColors.statusBlokir,
              fontSize: 11,
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
      ),
      DataCell(
        Icon(
          autodebet ? Icons.check_circle : Icons.cancel,
          color: autodebet ? AppColors.success : AppColors.divider,
          size: 18,
        ),
      ),
      DataCell(
        Row(
          children: [
            IconButton(
              icon: const Icon(Icons.visibility_outlined,
                  size: 18, color: AppColors.primary),
              onPressed: () {},
              tooltip: 'Detail',
            ),
            IconButton(
              icon: const Icon(Icons.schedule,
                  size: 18, color: AppColors.secondary),
              onPressed: () {},
              tooltip: 'Autodebet',
            ),
          ],
        ),
      ),
    ]);
  }
}
