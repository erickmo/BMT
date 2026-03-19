import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';

class FormListPage extends StatefulWidget {
  const FormListPage({super.key});

  @override
  State<FormListPage> createState() => _FormListPageState();
}

class _FormListPageState extends State<FormListPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text(AppStrings.formPengajuan),
        automaticallyImplyLeading: false,
        bottom: TabBar(
          controller: _tabController,
          labelColor: Colors.white,
          unselectedLabelColor: Colors.white60,
          indicatorColor: Colors.white,
          tabs: const [
            Tab(text: 'Menunggu'),
            Tab(text: 'Disetujui'),
            Tab(text: 'Ditolak'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _FormTabContent(status: 'MENUNGGU'),
          _FormTabContent(status: 'DISETUJUI'),
          _FormTabContent(status: 'DITOLAK'),
        ],
      ),
    );
  }
}

class _FormTabContent extends StatelessWidget {
  final String status;

  static const _forms = [
    {
      'id': 'F-001',
      'jenis': 'FORM_DAFTAR_NASABAH',
      'jenis_label': 'Daftar Nasabah',
      'nama': 'Ahmad Fauzi',
      'diajukan_oleh': 'Teller Siti',
      'created_at': '2025-03-19 09:15',
      'status': 'MENUNGGU',
    },
    {
      'id': 'F-002',
      'jenis': 'FORM_BUKA_REKENING',
      'jenis_label': 'Buka Rekening',
      'nama': 'Dewi Lestari',
      'diajukan_oleh': 'Teller Budi',
      'created_at': '2025-03-19 10:30',
      'status': 'MENUNGGU',
    },
    {
      'id': 'F-003',
      'jenis': 'FORM_DAFTAR_NASABAH',
      'jenis_label': 'Daftar Nasabah',
      'nama': 'Rudi Hartono',
      'diajukan_oleh': 'Teller Siti',
      'created_at': '2025-03-18 14:00',
      'status': 'DISETUJUI',
    },
    {
      'id': 'F-004',
      'jenis': 'FORM_TUTUP_REKENING',
      'jenis_label': 'Tutup Rekening',
      'nama': 'Budi Santoso',
      'diajukan_oleh': 'Teller Budi',
      'created_at': '2025-03-18 11:00',
      'status': 'DITOLAK',
    },
  ];

  const _FormTabContent({required this.status});

  @override
  Widget build(BuildContext context) {
    final filtered =
        _forms.where((f) => f['status'] == status).toList();

    return Padding(
      padding: const EdgeInsets.all(24),
      child: filtered.isEmpty
          ? const Center(
              child: Text(
                AppStrings.noData,
                style: TextStyle(color: AppColors.textSecondary),
              ),
            )
          : ListView.separated(
              itemCount: filtered.length,
              separatorBuilder: (_, __) => const SizedBox(height: 12),
              itemBuilder: (context, i) {
                final form = filtered[i];
                return _FormCard(form: form);
              },
            ),
    );
  }
}

class _FormCard extends StatelessWidget {
  final Map<String, String> form;

  const _FormCard({required this.form});

  @override
  Widget build(BuildContext context) {
    final isPending = form['status'] == 'MENUNGGU';
    final isApproved = form['status'] == 'DISETUJUI';

    final statusColor = isPending
        ? AppColors.warning
        : isApproved
            ? AppColors.success
            : AppColors.error;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: AppColors.primary.withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(Icons.description_outlined,
                  color: AppColors.primary, size: 24),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Text(
                        form['jenis_label']!,
                        style: const TextStyle(
                          fontWeight: FontWeight.bold,
                          fontSize: 15,
                          color: AppColors.textPrimary,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 6, vertical: 2),
                        decoration: BoxDecoration(
                          color: statusColor.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text(
                          form['status']!,
                          style: TextStyle(
                              color: statusColor,
                              fontSize: 11,
                              fontWeight: FontWeight.w600),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 4),
                  Text(
                    'Nasabah: ${form['nama']}',
                    style: const TextStyle(color: AppColors.textSecondary),
                  ),
                  Text(
                    'Diajukan oleh: ${form['diajukan_oleh']}  •  ${form['created_at']}',
                    style: const TextStyle(
                        color: AppColors.textHint, fontSize: 12),
                  ),
                ],
              ),
            ),
            if (isPending) ...[
              const SizedBox(width: 16),
              Column(
                children: [
                  ElevatedButton(
                    onPressed: () => _showApproveDialog(context),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: AppColors.success,
                      minimumSize: const Size(100, 36),
                      padding: const EdgeInsets.symmetric(horizontal: 12),
                    ),
                    child: const Text(AppStrings.setujui,
                        style: TextStyle(fontSize: 13)),
                  ),
                  const SizedBox(height: 8),
                  OutlinedButton(
                    onPressed: () => _showRejectDialog(context),
                    style: OutlinedButton.styleFrom(
                      side: const BorderSide(color: AppColors.error),
                      foregroundColor: AppColors.error,
                      minimumSize: const Size(100, 36),
                      padding: const EdgeInsets.symmetric(horizontal: 12),
                    ),
                    child: const Text(AppStrings.tolak,
                        style: TextStyle(fontSize: 13)),
                  ),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }

  void _showApproveDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (_) => AlertDialog(
        title: const Text('Setujui Form'),
        content: Text('Yakin ingin menyetujui form ${form['jenis_label']} untuk ${form['nama']}?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Batal'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context),
            style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.success),
            child: const Text('Setujui'),
          ),
        ],
      ),
    );
  }

  void _showRejectDialog(BuildContext context) {
    final reasonCtrl = TextEditingController();
    showDialog(
      context: context,
      builder: (_) => AlertDialog(
        title: const Text('Tolak Form'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Tolak form ${form['jenis_label']} untuk ${form['nama']}?'),
            const SizedBox(height: 16),
            TextField(
              controller: reasonCtrl,
              decoration: const InputDecoration(
                labelText: 'Alasan penolakan',
                hintText: 'Masukkan alasan...',
              ),
              maxLines: 3,
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Batal'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context),
            style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.error),
            child: const Text('Tolak'),
          ),
        ],
      ),
    );
  }
}
