import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';

class NilaiPage extends StatelessWidget {
  const NilaiPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text(AppStrings.nilai),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Row(
              children: [
                Expanded(
                  child: DropdownButtonFormField<String>(
                    value: 'Matematika',
                    decoration: const InputDecoration(labelText: 'Mata Pelajaran'),
                    items: ['Matematika', 'Bahasa Arab', 'Fiqih', 'IPA']
                        .map((m) => DropdownMenuItem(value: m, child: Text(m)))
                        .toList(),
                    onChanged: (_) {},
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: DropdownButtonFormField<String>(
                    value: 'UTS',
                    decoration: const InputDecoration(labelText: 'Komponen'),
                    items: ['UH1', 'UH2', 'UTS', 'UAS']
                        .map((k) => DropdownMenuItem(value: k, child: Text(k)))
                        .toList(),
                    onChanged: (_) {},
                  ),
                ),
                const SizedBox(width: 12),
                ElevatedButton.icon(
                  onPressed: () {},
                  icon: const Icon(Icons.save, size: 18),
                  label: const Text('Simpan Semua'),
                  style: ElevatedButton.styleFrom(
                      minimumSize: const Size(0, 48)),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Expanded(
              child: Card(
                child: SingleChildScrollView(
                  child: DataTable(
                    columns: const [
                      DataColumn(label: Text('No.')),
                      DataColumn(label: Text('Nama Santri')),
                      DataColumn(label: Text('NIS')),
                      DataColumn(label: Text('Nilai (0-100)')),
                      DataColumn(label: Text('Predikat')),
                    ],
                    rows: [
                      _buildRow('1', 'Muhammad Faqih', '2025-001', '85'),
                      _buildRow('2', 'Ahmad Habibi', '2025-002', '78'),
                      _buildRow('3', 'Ridwan Kamil Jr', '2025-003', '92'),
                      _buildRow('4', 'Zainul Arifin', '2024-089', null),
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
      String no, String nama, String nis, String? nilai) {
    final nilaiNum = double.tryParse(nilai ?? '');
    final predikat = nilaiNum == null
        ? '-'
        : nilaiNum >= 90
            ? 'A'
            : nilaiNum >= 80
                ? 'B'
                : nilaiNum >= 70
                    ? 'C'
                    : 'D';
    final predikatColor = predikat == 'A'
        ? AppColors.success
        : predikat == 'B'
            ? AppColors.primary
            : predikat == 'C'
                ? AppColors.warning
                : AppColors.error;

    return DataRow(cells: [
      DataCell(Text(no)),
      DataCell(Text(nama)),
      DataCell(Text(nis,
          style: const TextStyle(fontSize: 12, color: AppColors.textSecondary))),
      DataCell(
        SizedBox(
          width: 120,
          child: TextFormField(
            initialValue: nilai ?? '',
            keyboardType: TextInputType.number,
            decoration: const InputDecoration(
              hintText: '0-100',
              contentPadding:
                  EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            ),
          ),
        ),
      ),
      DataCell(
        nilai != null
            ? Container(
                padding: const EdgeInsets.symmetric(
                    horizontal: 10, vertical: 4),
                decoration: BoxDecoration(
                  color: predikatColor.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(6),
                ),
                child: Text(
                  predikat,
                  style: TextStyle(
                      color: predikatColor, fontWeight: FontWeight.bold),
                ),
              )
            : const Text('-'),
      ),
    ]);
  }
}
