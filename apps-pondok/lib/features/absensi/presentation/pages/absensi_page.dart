import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';

class AbsensiPage extends StatefulWidget {
  const AbsensiPage({super.key});

  @override
  State<AbsensiPage> createState() => _AbsensiPageState();
}

class _AbsensiPageState extends State<AbsensiPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  DateTime _selectedDate = DateTime.now();

  static const _santriList = [
    {'id': '1', 'nama': 'Muhammad Faqih', 'nis': '2025-001'},
    {'id': '2', 'nama': 'Ahmad Habibi', 'nis': '2025-002'},
    {'id': '3', 'nama': 'Ridwan Kamil Jr', 'nis': '2025-003'},
    {'id': '4', 'nama': 'Zainul Arifin', 'nis': '2024-089'},
  ];

  final Map<String, String> _statusMap = {};

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    for (final s in _santriList) {
      _statusMap[s['id']!] = 'HADIR';
    }
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Color _statusColor(String status) {
    switch (status) {
      case 'HADIR':
        return AppColors.statusHadir;
      case 'SAKIT':
        return AppColors.statusSakit;
      case 'IZIN':
        return AppColors.statusIzin;
      case 'ALFA':
        return AppColors.statusAlfa;
      default:
        return AppColors.textHint;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text(AppStrings.absensi),
        bottom: TabBar(
          controller: _tabController,
          labelColor: Colors.white,
          unselectedLabelColor: Colors.white60,
          indicatorColor: Colors.white,
          tabs: const [
            Tab(text: 'Input Absensi'),
            Tab(text: 'Rekap'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _InputAbsensiTab(
            santriList: _santriList,
            statusMap: _statusMap,
            selectedDate: _selectedDate,
            onDateChanged: (d) => setState(() => _selectedDate = d),
            onStatusChanged: (id, status) =>
                setState(() => _statusMap[id] = status),
            statusColor: _statusColor,
          ),
          _RekapAbsensiTab(),
        ],
      ),
      floatingActionButton: _tabController.index == 0
          ? FloatingActionButton.extended(
              onPressed: () {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(
                      content: Text('Absensi berhasil disimpan'),
                      backgroundColor: AppColors.success),
                );
              },
              icon: const Icon(Icons.save),
              label: const Text('Simpan Absensi'),
            )
          : null,
    );
  }
}

class _InputAbsensiTab extends StatelessWidget {
  final List<Map<String, String>> santriList;
  final Map<String, String> statusMap;
  final DateTime selectedDate;
  final ValueChanged<DateTime> onDateChanged;
  final Function(String, String) onStatusChanged;
  final Color Function(String) statusColor;

  const _InputAbsensiTab({
    required this.santriList,
    required this.statusMap,
    required this.selectedDate,
    required this.onDateChanged,
    required this.onStatusChanged,
    required this.statusColor,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        // Date selector
        Container(
          color: AppColors.surface,
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              const Icon(Icons.calendar_today, color: AppColors.primary),
              const SizedBox(width: 8),
              Text(
                formatTanggal(selectedDate),
                style: const TextStyle(
                    fontWeight: FontWeight.w600, fontSize: 16),
              ),
              const SizedBox(width: 12),
              OutlinedButton(
                onPressed: () async {
                  final picked = await showDatePicker(
                    context: context,
                    initialDate: selectedDate,
                    firstDate: DateTime(2024),
                    lastDate: DateTime.now(),
                  );
                  if (picked != null) onDateChanged(picked);
                },
                child: const Text('Ganti Tanggal'),
              ),
              const Spacer(),
              // Quick summary
              _SummaryChip(label: 'Hadir', count: statusMap.values.where((s) => s == 'HADIR').length, color: AppColors.statusHadir),
              const SizedBox(width: 8),
              _SummaryChip(label: 'Sakit', count: statusMap.values.where((s) => s == 'SAKIT').length, color: AppColors.statusSakit),
              const SizedBox(width: 8),
              _SummaryChip(label: 'Izin', count: statusMap.values.where((s) => s == 'IZIN').length, color: AppColors.statusIzin),
              const SizedBox(width: 8),
              _SummaryChip(label: 'Alfa', count: statusMap.values.where((s) => s == 'ALFA').length, color: AppColors.statusAlfa),
            ],
          ),
        ),
        const Divider(height: 1),
        Expanded(
          child: ListView.separated(
            padding: const EdgeInsets.all(16),
            itemCount: santriList.length,
            separatorBuilder: (_, __) => const SizedBox(height: 8),
            itemBuilder: (context, i) {
              final s = santriList[i];
              final currentStatus = statusMap[s['id']] ?? 'HADIR';
              return Card(
                child: Padding(
                  padding: const EdgeInsets.symmetric(
                      horizontal: 16, vertical: 8),
                  child: Row(
                    children: [
                      CircleAvatar(
                        backgroundColor:
                            AppColors.primary.withOpacity(0.1),
                        child: Text(
                          s['nama']!.substring(0, 1),
                          style: const TextStyle(
                              color: AppColors.primary,
                              fontWeight: FontWeight.bold),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(s['nama']!,
                                style: const TextStyle(
                                    fontWeight: FontWeight.w600)),
                            Text(s['nis']!,
                                style: const TextStyle(
                                    fontSize: 12,
                                    color: AppColors.textSecondary)),
                          ],
                        ),
                      ),
                      // Status buttons
                      for (final status in ['HADIR', 'SAKIT', 'IZIN', 'ALFA'])
                        Padding(
                          padding: const EdgeInsets.only(left: 6),
                          child: GestureDetector(
                            onTap: () => onStatusChanged(s['id']!, status),
                            child: Container(
                              padding: const EdgeInsets.symmetric(
                                  horizontal: 10, vertical: 6),
                              decoration: BoxDecoration(
                                color: currentStatus == status
                                    ? statusColor(status)
                                    : statusColor(status).withOpacity(0.1),
                                borderRadius: BorderRadius.circular(6),
                              ),
                              child: Text(
                                status,
                                style: TextStyle(
                                  color: currentStatus == status
                                      ? Colors.white
                                      : statusColor(status),
                                  fontSize: 11,
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                            ),
                          ),
                        ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}

class _SummaryChip extends StatelessWidget {
  final String label;
  final int count;
  final Color color;

  const _SummaryChip({required this.label, required this.count, required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color.withOpacity(0.3)),
      ),
      child: Text(
        '$label: $count',
        style: TextStyle(color: color, fontSize: 12, fontWeight: FontWeight.w600),
      ),
    );
  }
}

class _RekapAbsensiTab extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  const Text('Rekap Absensi',
                      style: TextStyle(
                          fontWeight: FontWeight.bold, fontSize: 16)),
                  const Spacer(),
                  OutlinedButton.icon(
                    onPressed: () {},
                    icon: const Icon(Icons.file_download, size: 16),
                    label: const Text('Export'),
                  ),
                ],
              ),
              const Divider(height: 24),
              DataTable(
                columns: const [
                  DataColumn(label: Text('Santri')),
                  DataColumn(label: Text('Hadir')),
                  DataColumn(label: Text('Sakit')),
                  DataColumn(label: Text('Izin')),
                  DataColumn(label: Text('Alfa')),
                  DataColumn(label: Text('%')),
                ],
                rows: const [
                  DataRow(cells: [
                    DataCell(Text('Muhammad Faqih')),
                    DataCell(Text('20')),
                    DataCell(Text('1')),
                    DataCell(Text('0')),
                    DataCell(Text('0')),
                    DataCell(Text('95.2%', style: TextStyle(color: AppColors.success, fontWeight: FontWeight.w600))),
                  ]),
                  DataRow(cells: [
                    DataCell(Text('Ahmad Habibi')),
                    DataCell(Text('18')),
                    DataCell(Text('2')),
                    DataCell(Text('1')),
                    DataCell(Text('0')),
                    DataCell(Text('85.7%', style: TextStyle(color: AppColors.warning, fontWeight: FontWeight.w600))),
                  ]),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
