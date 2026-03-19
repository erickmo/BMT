import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';
import '../../../auth/presentation/bloc/auth_bloc.dart';
import '../../../rekening/domain/entities/rekening_entity.dart';
import '../bloc/home_bloc.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  bool _hideSaldo = false;

  @override
  void initState() {
    super.initState();
    context.read<HomeBloc>().add(const LoadDashboard());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      body: BlocBuilder<HomeBloc, HomeState>(
        builder: (context, state) {
          return RefreshIndicator(
            onRefresh: () async {
              context.read<HomeBloc>().add(const LoadDashboard());
            },
            child: CustomScrollView(
              slivers: [
                // App bar dengan saldo
                SliverAppBar(
                  expandedHeight: 200,
                  pinned: true,
                  backgroundColor: AppColors.primary,
                  actions: [
                    IconButton(
                      icon: const Icon(Icons.notifications_outlined,
                          color: Colors.white),
                      onPressed: () {},
                    ),
                    IconButton(
                      icon: const Icon(Icons.person_outline,
                          color: Colors.white),
                      onPressed: () => context.push('/profil'),
                    ),
                  ],
                  flexibleSpace: FlexibleSpaceBar(
                    background: _buildHeaderContent(context, state),
                  ),
                ),

                // Quick actions
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.all(AppSizes.pagePadding),
                    child: _buildQuickActions(context),
                  ),
                ),

                // Rekening summary
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                      horizontal: AppSizes.pagePadding,
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          AppStrings.rekening,
                          style: TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                        ),
                        TextButton(
                          onPressed: () => context.push('/rekening'),
                          child: const Text('Lihat Semua'),
                        ),
                      ],
                    ),
                  ),
                ),

                // Rekening list
                if (state is HomeLoaded)
                  SliverList(
                    delegate: SliverChildBuilderDelegate(
                      (context, index) {
                        final rek = state.dashboard.rekening[index];
                        return _RekeningCard(rekening: rek);
                      },
                      childCount: state.dashboard.rekening.length,
                    ),
                  ),

                if (state is HomeLoading)
                  const SliverToBoxAdapter(
                    child: Center(child: CircularProgressIndicator()),
                  ),

                if (state is HomeError)
                  SliverToBoxAdapter(
                    child: Center(
                      child: Padding(
                        padding: const EdgeInsets.all(AppSizes.xl),
                        child: Column(
                          children: [
                            Text(state.message,
                                style:
                                    const TextStyle(color: AppColors.error)),
                            const SizedBox(height: AppSizes.md),
                            ElevatedButton(
                              onPressed: () => context
                                  .read<HomeBloc>()
                                  .add(const LoadDashboard()),
                              child: const Text(AppStrings.retry),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),

                const SliverToBoxAdapter(child: SizedBox(height: AppSizes.xl)),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildHeaderContent(BuildContext context, HomeState state) {
    return Container(
      padding: const EdgeInsets.fromLTRB(
        AppSizes.pagePadding,
        kToolbarHeight + AppSizes.lg,
        AppSizes.pagePadding,
        AppSizes.lg,
      ),
      decoration: const BoxDecoration(
        gradient: LinearGradient(
          colors: [AppColors.primary, AppColors.primaryLight],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (state is HomeLoaded)
            Text(
              'Halo, ${state.dashboard.namaNasabah}',
              style: const TextStyle(
                color: Colors.white70,
                fontSize: 14,
              ),
            ),
          const SizedBox(height: AppSizes.xs),
          const Text(
            AppStrings.totalSaldo,
            style: TextStyle(color: Colors.white70, fontSize: 13),
          ),
          Row(
            children: [
              Text(
                state is HomeLoaded
                    ? (_hideSaldo
                        ? '••••••••'
                        : formatRupiah(state.dashboard.totalSaldo))
                    : '---',
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 26,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(width: AppSizes.sm),
              GestureDetector(
                onTap: () => setState(() => _hideSaldo = !_hideSaldo),
                child: Icon(
                  _hideSaldo ? Icons.visibility_off : Icons.visibility,
                  color: Colors.white70,
                  size: 20,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildQuickActions(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
      children: [
        _QuickAction(
          icon: Icons.swap_horiz,
          label: AppStrings.transaksi,
          onTap: () => context.push('/rekening'),
        ),
        _QuickAction(
          icon: Icons.history,
          label: 'Riwayat',
          onTap: () => context.push('/rekening'),
        ),
        _QuickAction(
          icon: Icons.person,
          label: AppStrings.profil,
          onTap: () => context.push('/profil'),
        ),
        _QuickAction(
          icon: Icons.logout,
          label: AppStrings.logout,
          onTap: () {
            context.read<AuthBloc>().add(const LogoutRequested());
          },
        ),
      ],
    );
  }
}

class _QuickAction extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  const _QuickAction({
    required this.icon,
    required this.label,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Column(
        children: [
          Container(
            width: 52,
            height: 52,
            decoration: BoxDecoration(
              color: AppColors.primary.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(AppSizes.radiusMd),
            ),
            child: Icon(icon, color: AppColors.primary),
          ),
          const SizedBox(height: AppSizes.xs),
          Text(
            label,
            style: const TextStyle(
              fontSize: 12,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }
}

class _RekeningCard extends StatelessWidget {
  final RekeningEntity rekening;

  const _RekeningCard({required this.rekening});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.fromLTRB(
        AppSizes.pagePadding,
        0,
        AppSizes.pagePadding,
        AppSizes.md,
      ),
      child: InkWell(
        borderRadius: BorderRadius.circular(AppSizes.radiusMd),
        onTap: () => context.push('/rekening/${rekening.id}'),
        child: Padding(
          padding: const EdgeInsets.all(AppSizes.cardPadding),
          child: Row(
            children: [
              Container(
                padding: const EdgeInsets.all(AppSizes.sm),
                decoration: BoxDecoration(
                  color: AppColors.primary.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(AppSizes.radiusSm),
                ),
                child: const Icon(
                  Icons.account_balance_wallet,
                  color: AppColors.primary,
                  size: 22,
                ),
              ),
              const SizedBox(width: AppSizes.md),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      rekening.jenisRekeningNama,
                      style: const TextStyle(fontWeight: FontWeight.w600),
                    ),
                    Text(
                      maskNomorRekening(rekening.nomorRekening),
                      style: const TextStyle(
                        color: AppColors.textSecondary,
                        fontSize: 13,
                      ),
                    ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    formatRupiah(rekening.saldo),
                    style: const TextStyle(
                      fontWeight: FontWeight.bold,
                      color: AppColors.primary,
                    ),
                  ),
                  Text(
                    rekening.status,
                    style: TextStyle(
                      fontSize: 11,
                      color: rekening.isAktif
                          ? AppColors.success
                          : AppColors.error,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
