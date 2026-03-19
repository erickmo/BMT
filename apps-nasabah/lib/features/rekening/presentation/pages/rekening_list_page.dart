import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';
import '../bloc/rekening_bloc.dart';

class RekeningListPage extends StatefulWidget {
  const RekeningListPage({super.key});

  @override
  State<RekeningListPage> createState() => _RekeningListPageState();
}

class _RekeningListPageState extends State<RekeningListPage> {
  @override
  void initState() {
    super.initState();
    context.read<RekeningBloc>().add(const LoadDaftarRekening());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(AppStrings.daftarRekening),
      ),
      body: BlocBuilder<RekeningBloc, RekeningState>(
        builder: (context, state) {
          if (state is RekeningLoading) {
            return const Center(child: CircularProgressIndicator());
          }

          if (state is RekeningError) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.error_outline,
                      size: 48, color: AppColors.error),
                  const SizedBox(height: AppSizes.md),
                  Text(state.message,
                      textAlign: TextAlign.center,
                      style:
                          const TextStyle(color: AppColors.textSecondary)),
                  const SizedBox(height: AppSizes.md),
                  ElevatedButton(
                    onPressed: () => context
                        .read<RekeningBloc>()
                        .add(const LoadDaftarRekening()),
                    child: const Text(AppStrings.retry),
                  ),
                ],
              ),
            );
          }

          if (state is DaftarRekeningLoaded) {
            if (state.rekening.isEmpty) {
              return const Center(
                child: Text(
                  'Belum ada rekening',
                  style: TextStyle(color: AppColors.textSecondary),
                ),
              );
            }

            return RefreshIndicator(
              onRefresh: () async {
                context.read<RekeningBloc>().add(const LoadDaftarRekening());
              },
              child: ListView.builder(
                padding: const EdgeInsets.all(AppSizes.pagePadding),
                itemCount: state.rekening.length,
                itemBuilder: (context, index) {
                  final rek = state.rekening[index];
                  return Card(
                    margin: const EdgeInsets.only(bottom: AppSizes.md),
                    child: InkWell(
                      borderRadius:
                          BorderRadius.circular(AppSizes.radiusMd),
                      onTap: () => context.push('/rekening/${rek.id}'),
                      child: Padding(
                        padding: const EdgeInsets.all(AppSizes.cardPadding),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Container(
                                  padding: const EdgeInsets.all(AppSizes.sm),
                                  decoration: BoxDecoration(
                                    color:
                                        AppColors.primary.withValues(alpha: 0.1),
                                    borderRadius: BorderRadius.circular(
                                        AppSizes.radiusSm),
                                  ),
                                  child: const Icon(
                                    Icons.account_balance_wallet,
                                    color: AppColors.primary,
                                  ),
                                ),
                                const SizedBox(width: AppSizes.md),
                                Expanded(
                                  child: Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      Text(
                                        rek.jenisRekeningNama,
                                        style: const TextStyle(
                                          fontWeight: FontWeight.w600,
                                          fontSize: 15,
                                        ),
                                      ),
                                      Text(
                                        maskNomorRekening(rek.nomorRekening),
                                        style: const TextStyle(
                                          color: AppColors.textSecondary,
                                          fontSize: 13,
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                                Container(
                                  padding: const EdgeInsets.symmetric(
                                    horizontal: AppSizes.sm,
                                    vertical: AppSizes.xs,
                                  ),
                                  decoration: BoxDecoration(
                                    color: rek.isAktif
                                        ? AppColors.success.withValues(alpha: 0.15)
                                        : AppColors.error.withValues(alpha: 0.15),
                                    borderRadius: BorderRadius.circular(
                                        AppSizes.radiusFull),
                                  ),
                                  child: Text(
                                    rek.status,
                                    style: TextStyle(
                                      fontSize: 11,
                                      fontWeight: FontWeight.w600,
                                      color: rek.isAktif
                                          ? AppColors.success
                                          : AppColors.error,
                                    ),
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: AppSizes.md),
                            const Text(
                              AppStrings.saldo,
                              style: TextStyle(
                                fontSize: 12,
                                color: AppColors.textSecondary,
                              ),
                            ),
                            Text(
                              formatRupiah(rek.saldo),
                              style: const TextStyle(
                                fontSize: 22,
                                fontWeight: FontWeight.bold,
                                color: AppColors.primary,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  );
                },
              ),
            );
          }

          return const SizedBox.shrink();
        },
      ),
    );
  }
}
