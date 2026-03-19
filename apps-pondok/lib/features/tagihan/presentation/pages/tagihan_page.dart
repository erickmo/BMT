import 'package:flutter/material.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';

class TagihanPage extends StatelessWidget {
  const TagihanPage({super.key});

  static const _tagihan = [
    {
      'santri': 'Muhammad Faqih',
      'nis': '2025-001',
      'periode': '2025-03',
      'nominal': 350000,
      'terbayar': 350000,
      'status': 'LUNAS',
      'beasiswa': 0,
    },
    {
      'santri': 'Ahmad Habibi',
      'nis': '2025-002',
      'periode': '2025-03',
      'nominal': 350000,
      'terbayar': 200000,
      'status': 'SEBAGIAN',
      'beasiswa': 0,
    },
    {
      'santri': 'Ridwan Kamil Jr',
      'nis': '2025-003',
      'periode': '2025-03',
      'nominal': 350000,
      'terbayar': 0,
      'status': 'BELUM_BAYAR',
      'beasiswa': 175000,
    },
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text(AppStrings.tagihan),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 16),
            child: ElevatedButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.generate_tokens, size: 18),
              label: const Text('Generate Tagihan'),
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.white,
                foregroundColor: AppColors.primary,
                minimumSize: const Size(0, 36),
              ),
            ),
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            // Period filter
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
                  child: DropdownButtonFormField<String>(
                    value: 'SPP Bulanan',
                    decoration: const InputDecoration(labelText: 'Jenis Tagihan'),
                    items: ['SPP Bulanan', 'Kitab', 'Ekstra']
                        .map((j) => DropdownMenuItem(value: j, child: Text(j)))
                        .toList(),
                    onChanged: (_) {},
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Expanded(
              child: ListView.separated(
                itemCount: _tagihan.length,
                separatorBuilder: (_, __) => const SizedBox(height: 8),
                itemBuilder: (context, i) {
                  final t = _tagihan[i];
                  final nominal = t['nominal'] as int;
                  final terbayar = t['terbayar'] as int;
                  final beasiswa = t['beasiswa'] as int;
                  final efektif = nominal - beasiswa;
                  final status = t['status'] as String;

                  final statusColor = status == 'LUNAS'
                      ? AppColors.success
                      : status == 'SEBAGIAN'
                          ? AppColors.warning
                          : AppColors.error;

                  return Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            children: [
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(t['santri'] as String,
                                        style: const TextStyle(
                                            fontWeight: FontWeight.bold,
                                            fontSize: 15)),
                                    Text(t['nis'] as String,
                                        style: const TextStyle(
                                            fontSize: 12,
                                            color: AppColors.textSecondary)),
                                  ],
                                ),
                              ),
                              Container(
                                padding: const EdgeInsets.symmetric(
                                    horizontal: 10, vertical: 4),
                                decoration: BoxDecoration(
                                  color: statusColor.withOpacity(0.1),
                                  borderRadius: BorderRadius.circular(6),
                                ),
                                child: Text(
                                  status.replaceAll('_', ' '),
                                  style: TextStyle(
                                      color: statusColor,
                                      fontWeight: FontWeight.w600,
                                      fontSize: 12),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 12),
                          Row(
                            children: [
                              _TagihanInfo(
                                  label: 'Nominal',
                                  value: formatRupiah(nominal)),
                              if (beasiswa > 0)
                                _TagihanInfo(
                                    label: 'Beasiswa',
                                    value: formatRupiah(beasiswa),
                                    valueColor: AppColors.success),
                              _TagihanInfo(
                                  label: 'Efektif',
                                  value: formatRupiah(efektif),
                                  isBold: true),
                              _TagihanInfo(
                                  label: 'Terbayar',
                                  value: formatRupiah(terbayar),
                                  valueColor: AppColors.primary),
                              const Spacer(),
                              if (status != 'LUNAS')
                                ElevatedButton.icon(
                                  onPressed: () {},
                                  icon: const Icon(Icons.payments, size: 16),
                                  label: const Text('Bayar'),
                                  style: ElevatedButton.styleFrom(
                                      minimumSize: const Size(0, 36),
                                      padding: const EdgeInsets.symmetric(
                                          horizontal: 12)),
                                ),
                              const SizedBox(width: 8),
                              OutlinedButton.icon(
                                onPressed: () {},
                                icon: const Icon(Icons.school, size: 16),
                                label: const Text('Beasiswa'),
                                style: OutlinedButton.styleFrom(
                                    minimumSize: const Size(0, 36),
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 12)),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _TagihanInfo extends StatelessWidget {
  final String label;
  final String value;
  final Color? valueColor;
  final bool isBold;

  const _TagihanInfo({
    required this.label,
    required this.value,
    this.valueColor,
    this.isBold = false,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(right: 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(label,
              style: const TextStyle(
                  fontSize: 11, color: AppColors.textSecondary)),
          Text(
            value,
            style: TextStyle(
              fontSize: 13,
              fontWeight:
                  isBold ? FontWeight.bold : FontWeight.normal,
              color: valueColor ?? AppColors.textPrimary,
            ),
          ),
        ],
      ),
    );
  }
}
