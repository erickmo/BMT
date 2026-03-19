import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';
import '../bloc/rekening_bloc.dart';

class RekeningDetailPage extends StatefulWidget {
  final String rekeningId;

  const RekeningDetailPage({super.key, required this.rekeningId});

  @override
  State<RekeningDetailPage> createState() => _RekeningDetailPageState();
}

class _RekeningDetailPageState extends State<RekeningDetailPage> {
  @override
  void initState() {
    super.initState();
    context.read<RekeningBloc>().add(LoadDetailRekening(widget.rekeningId));
    context.read<RekeningBloc>().add(
          LoadRiwayatTransaksi(rekeningId: widget.rekeningId),
        );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text(AppStrings.detailRekening)),
      body: BlocBuilder<RekeningBloc, RekeningState>(
        builder: (context, state) {
          if (state is RekeningLoading) {
            return const Center(child: CircularProgressIndicator());
          }

          if (state is RekeningError) {
            return Center(
              child: Text(
                state.message,
                style: const TextStyle(color: AppColors.error),
              ),
            );
          }

          if (state is DetailRekeningLoaded) {
            final rek = state.rekening;
            return RefreshIndicator(
              onRefresh: () async {
                context.read<RekeningBloc>()
                  ..add(LoadDetailRekening(widget.rekeningId))
                  ..add(
                    LoadRiwayatTransaksi(rekeningId: widget.rekeningId),
                  );
              },
              child: ListView(
                children: [
                  // Saldo card
                  Container(
                    margin: const EdgeInsets.all(AppSizes.pagePadding),
                    padding: const EdgeInsets.all(AppSizes.lg),
                    decoration: BoxDecoration(
                      gradient: const LinearGradient(
                        colors: [AppColors.primary, AppColors.primaryLight],
                        begin: Alignment.topLeft,
                        end: Alignment.bottomRight,
                      ),
                      borderRadius:
                          BorderRadius.circular(AppSizes.radiusLg),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          rek.jenisRekeningNama,
                          style: const TextStyle(
                            color: Colors.white70,
                            fontSize: 14,
                          ),
                        ),
                        const SizedBox(height: AppSizes.xs),
                        Text(
                          rek.nomorRekening,
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 16,
                            fontWeight: FontWeight.w500,
                            letterSpacing: 1,
                          ),
                        ),
                        const SizedBox(height: AppSizes.lg),
                        const Text(
                          AppStrings.saldo,
                          style: TextStyle(color: Colors.white70, fontSize: 13),
                        ),
                        Text(
                          formatRupiah(rek.saldo),
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 28,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
                  ),

                  // Info rekening
                  Padding(
                    padding:
                        const EdgeInsets.symmetric(horizontal: AppSizes.pagePadding),
                    child: Card(
                      child: Padding(
                        padding: const EdgeInsets.all(AppSizes.cardPadding),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Informasi Rekening',
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                fontSize: 15,
                              ),
                            ),
                            const SizedBox(height: AppSizes.md),
                            _InfoRow(
                              label: AppStrings.statusRekening,
                              value: rek.status,
                              valueColor: rek.isAktif
                                  ? AppColors.success
                                  : AppColors.error,
                            ),
                            _InfoRow(
                              label: 'Tanggal Buka',
                              value: formatTanggal(rek.tanggalBuka),
                            ),
                            if (rek.biayaAdminBulanan > 0)
                              _InfoRow(
                                label: 'Biaya Admin/Bulan',
                                value: formatRupiah(rek.biayaAdminBulanan),
                              ),
                          ],
                        ),
                      ),
                    ),
                  ),

                  // Riwayat transaksi
                  const Padding(
                    padding: EdgeInsets.fromLTRB(
                      AppSizes.pagePadding,
                      AppSizes.lg,
                      AppSizes.pagePadding,
                      AppSizes.sm,
                    ),
                    child: Text(
                      AppStrings.riwayatTransaksi,
                      style: TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                  ),

                  if (state.isLoadingMore)
                    const Padding(
                      padding: EdgeInsets.all(AppSizes.md),
                      child: Center(child: CircularProgressIndicator()),
                    ),

                  if (state.transaksi.isEmpty && !state.isLoadingMore)
                    const Padding(
                      padding: EdgeInsets.all(AppSizes.lg),
                      child: Center(
                        child: Text(
                          'Belum ada transaksi',
                          style: TextStyle(color: AppColors.textSecondary),
                        ),
                      ),
                    ),

                  ...state.transaksi.map(
                    (tx) => ListTile(
                      leading: CircleAvatar(
                        backgroundColor: tx.isKredit
                            ? AppColors.success.withValues(alpha: 0.15)
                            : AppColors.error.withValues(alpha: 0.15),
                        child: Icon(
                          tx.isKredit ? Icons.arrow_downward : Icons.arrow_upward,
                          color:
                              tx.isKredit ? AppColors.success : AppColors.error,
                          size: 20,
                        ),
                      ),
                      title: Text(
                        tx.keterangan.isNotEmpty ? tx.keterangan : tx.tipe,
                        style: const TextStyle(fontSize: 14),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      subtitle: Text(
                        formatTanggalWaktu(tx.tanggal),
                        style: const TextStyle(
                          fontSize: 12,
                          color: AppColors.textSecondary,
                        ),
                      ),
                      trailing: Text(
                        '${tx.isKredit ? '+' : '-'}${formatRupiah(tx.nominal)}',
                        style: TextStyle(
                          fontWeight: FontWeight.w600,
                          color:
                              tx.isKredit ? AppColors.success : AppColors.error,
                          fontSize: 14,
                        ),
                      ),
                    ),
                  ),

                  const SizedBox(height: AppSizes.xl),
                ],
              ),
            );
          }

          return const SizedBox.shrink();
        },
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  final String label;
  final String value;
  final Color? valueColor;

  const _InfoRow({
    required this.label,
    required this.value,
    this.valueColor,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: AppSizes.xs),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: const TextStyle(
              color: AppColors.textSecondary,
              fontSize: 13,
            ),
          ),
          Text(
            value,
            style: TextStyle(
              fontWeight: FontWeight.w500,
              fontSize: 13,
              color: valueColor,
            ),
          ),
        ],
      ),
    );
  }
}
