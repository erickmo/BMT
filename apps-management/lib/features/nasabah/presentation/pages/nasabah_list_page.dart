import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';

class NasabahListPage extends StatefulWidget {
  const NasabahListPage({super.key});

  @override
  State<NasabahListPage> createState() => _NasabahListPageState();
}

class _NasabahListPageState extends State<NasabahListPage> {
  final _searchCtrl = TextEditingController();

  @override
  void dispose() {
    _searchCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text(AppStrings.nasabah),
        automaticallyImplyLeading: false,
      ),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  flex: 3,
                  child: TextField(
                    controller: _searchCtrl,
                    decoration: const InputDecoration(
                      hintText: AppStrings.cariNasabah,
                      prefixIcon: Icon(Icons.search),
                    ),
                    onChanged: (v) => setState(() {}),
                  ),
                ),
                const SizedBox(width: 12),
                OutlinedButton.icon(
                  onPressed: () {},
                  icon: const Icon(Icons.filter_list, size: 16),
                  label: const Text('Filter'),
                  style: OutlinedButton.styleFrom(minimumSize: const Size(0, 44)),
                ),
                const SizedBox(width: 12),
                ElevatedButton.icon(
                  onPressed: () {},
                  icon: const Icon(Icons.file_download, size: 16),
                  label: const Text('Export'),
                  style: ElevatedButton.styleFrom(minimumSize: const Size(0, 44)),
                ),
              ],
            ),
            const SizedBox(height: 20),
            Expanded(
              child: Card(
                child: _NasabahDataTable(),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _NasabahDataTable extends StatelessWidget {
  final _data = const [
    {
      'nomor': 'BMT-001-00001',
      'nama': 'Ahmad Fauzi',
      'nik': '3512010101010001',
      'telepon': '08123456789',
      'rekening': '3',
      'status': 'AKTIF',
    },
    {
      'nomor': 'BMT-001-00002',
      'nama': 'Siti Rahayu',
      'nik': '3512010101010002',
      'telepon': '08234567890',
      'rekening': '2',
      'status': 'AKTIF',
    },
    {
      'nomor': 'BMT-001-00003',
      'nama': 'Budi Santoso',
      'nik': '3512010101010003',
      'telepon': '08345678901',
      'rekening': '1',
      'status': 'BLOKIR',
    },
  ];

  const _NasabahDataTable();

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: DataTable(
        columnSpacing: 24,
        columns: const [
          DataColumn(label: Text('No. Nasabah')),
          DataColumn(label: Text('Nama')),
          DataColumn(label: Text('NIK')),
          DataColumn(label: Text('Telepon')),
          DataColumn(label: Text('Rekening')),
          DataColumn(label: Text('Status')),
          DataColumn(label: Text('Aksi')),
        ],
        rows: _data.map((n) {
          final isAktif = n['status'] == 'AKTIF';
          return DataRow(cells: [
            DataCell(Text(n['nomor']!,
                style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    color: AppColors.primary))),
            DataCell(Text(n['nama']!)),
            DataCell(Text(n['nik']!)),
            DataCell(Text(n['telepon']!)),
            DataCell(Text(n['rekening']!)),
            DataCell(
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                decoration: BoxDecoration(
                  color: isAktif
                      ? AppColors.statusAktif.withOpacity(0.1)
                      : AppColors.statusBlokir.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(4),
                ),
                child: Text(
                  n['status']!,
                  style: TextStyle(
                    color: isAktif
                        ? AppColors.statusAktif
                        : AppColors.statusBlokir,
                    fontSize: 11,
                    fontWeight: FontWeight.w600,
                  ),
                ),
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
                ],
              ),
            ),
          ]);
        }).toList(),
      ),
    );
  }
}
