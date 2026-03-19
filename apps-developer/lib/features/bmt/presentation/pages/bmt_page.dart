import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';

class BmtPage extends StatelessWidget {
  const BmtPage({super.key});

  static const _bmtData = [
    {
      'id': 'bmt-001',
      'kode': 'ANNUR',
      'nama': 'BMT An-Nur Kediri',
      'pic': 'Ahmad Ghazali',
      'telepon': '0354-123456',
      'cabang': 3,
      'status': 'AKTIF',
      'kontrak': '2025-12-31',
    },
    {
      'id': 'bmt-002',
      'kode': 'SALAM',
      'nama': 'BMT Salam Jombang',
      'pic': 'Budi Santoso',
      'telepon': '0321-654321',
      'cabang': 2,
      'status': 'AKTIF',
      'kontrak': '2024-06-30',
    },
    {
      'id': 'bmt-003',
      'kode': 'BADR',
      'nama': 'BMT Badr Malang',
      'pic': 'Eko Wahyudi',
      'telepon': '0341-777888',
      'cabang': 1,
      'status': 'SUSPEND',
      'kontrak': '2023-12-31',
    },
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Manajemen BMT'),
        automaticallyImplyLeading: false,
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 16),
            child: ElevatedButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.add, size: 18),
              label: const Text('Daftarkan BMT'),
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
          children: [
            const TextField(
              decoration: InputDecoration(
                hintText: 'Cari BMT...',
                prefixIcon: Icon(Icons.search),
              ),
            ),
            const SizedBox(height: 16),
            Expanded(
              child: Card(
                child: SingleChildScrollView(
                  child: DataTable(
                    columnSpacing: 20,
                    columns: const [
                      DataColumn(label: Text('Kode')),
                      DataColumn(label: Text('Nama BMT')),
                      DataColumn(label: Text('PIC')),
                      DataColumn(label: Text('Cabang')),
                      DataColumn(label: Text('Kontrak s/d')),
                      DataColumn(label: Text('Status')),
                      DataColumn(label: Text('Aksi')),
                    ],
                    rows: _bmtData.map((bmt) {
                      final isAktif = bmt['status'] == 'AKTIF';
                      final kontrakDate = DateTime.tryParse(bmt['kontrak'] as String);
                      final isExpired = kontrakDate != null &&
                          kontrakDate.isBefore(DateTime.now());
                      return DataRow(cells: [
                        DataCell(
                          Container(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 8, vertical: 4),
                            decoration: BoxDecoration(
                              color: AppColors.codeBg,
                              borderRadius: BorderRadius.circular(4),
                            ),
                            child: Text(
                              bmt['kode'] as String,
                              style: const TextStyle(
                                  fontFamily: 'monospace',
                                  color: AppColors.codeText,
                                  fontSize: 12),
                            ),
                          ),
                        ),
                        DataCell(Text(bmt['nama'] as String,
                            style: const TextStyle(fontWeight: FontWeight.w600))),
                        DataCell(Text(bmt['pic'] as String)),
                        DataCell(Text('${bmt['cabang']} cabang')),
                        DataCell(
                          Text(
                            bmt['kontrak'] as String,
                            style: TextStyle(
                                color: isExpired ? AppColors.error : AppColors.textPrimary,
                                fontWeight: isExpired ? FontWeight.bold : FontWeight.normal),
                          ),
                        ),
                        DataCell(
                          Container(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 8, vertical: 3),
                            decoration: BoxDecoration(
                              color: isAktif
                                  ? AppColors.success.withOpacity(0.1)
                                  : AppColors.error.withOpacity(0.1),
                              borderRadius: BorderRadius.circular(4),
                            ),
                            child: Text(
                              bmt['status'] as String,
                              style: TextStyle(
                                  color: isAktif
                                      ? AppColors.success
                                      : AppColors.error,
                                  fontSize: 11,
                                  fontWeight: FontWeight.w600),
                            ),
                          ),
                        ),
                        DataCell(Row(
                          children: [
                            IconButton(
                              icon: const Icon(Icons.edit_outlined,
                                  size: 18, color: AppColors.accent),
                              onPressed: () {},
                              tooltip: 'Edit',
                            ),
                            IconButton(
                              icon: const Icon(Icons.corporate_fare,
                                  size: 18, color: AppColors.primary),
                              onPressed: () {},
                              tooltip: 'Cabang',
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
