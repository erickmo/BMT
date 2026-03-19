import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/utils/formatters.dart';

class UsageLogPage extends StatelessWidget {
  const UsageLogPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Usage Log'),
        automaticallyImplyLeading: false,
      ),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          children: [
            Row(
              children: [
                Expanded(
                  child: DropdownButtonFormField<String>(
                    value: '2025-03',
                    decoration: const InputDecoration(labelText: 'Periode'),
                    items: ['2025-01', '2025-02', '2025-03']
                        .map((p) => DropdownMenuItem(
                            value: p, child: Text(formatPeriode(p))))
                        .toList(),
                    onChanged: (_) {},
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: const TextField(
                    decoration: InputDecoration(
                        hintText: 'Filter BMT...', prefixIcon: Icon(Icons.search)),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Expanded(
              child: Card(
                child: SingleChildScrollView(
                  child: DataTable(
                    columns: const [
                      DataColumn(label: Text('BMT')),
                      DataColumn(label: Text('Transaksi')),
                      DataColumn(label: Text('Setor/Tarik')),
                      DataColumn(label: Text('Pembiayaan')),
                      DataColumn(label: Text('NFC')),
                      DataColumn(label: Text('OPOP')),
                      DataColumn(label: Text('Total Fee')),
                    ],
                    rows: const [
                      DataRow(cells: [
                        DataCell(Text('BMT An-Nur',
                            style: TextStyle(fontWeight: FontWeight.w600))),
                        DataCell(Text('3.248')),
                        DataCell(Text('2.100')),
                        DataCell(Text('145')),
                        DataCell(Text('870')),
                        DataCell(Text('133')),
                        DataCell(Text('Rp 1.250.000',
                            style: TextStyle(
                                fontWeight: FontWeight.w600,
                                color: AppColors.accent))),
                      ]),
                      DataRow(cells: [
                        DataCell(Text('BMT Salam',
                            style: TextStyle(fontWeight: FontWeight.w600))),
                        DataCell(Text('1.890')),
                        DataCell(Text('1.200')),
                        DataCell(Text('87')),
                        DataCell(Text('540')),
                        DataCell(Text('63')),
                        DataCell(Text('Rp 750.000',
                            style: TextStyle(
                                fontWeight: FontWeight.w600,
                                color: AppColors.accent))),
                      ]),
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
