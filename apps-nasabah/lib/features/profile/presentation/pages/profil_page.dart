import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';
import '../../../auth/presentation/bloc/auth_bloc.dart';
import '../bloc/profil_bloc.dart';

class ProfilPage extends StatefulWidget {
  const ProfilPage({super.key});

  @override
  State<ProfilPage> createState() => _ProfilPageState();
}

class _ProfilPageState extends State<ProfilPage> {
  @override
  void initState() {
    super.initState();
    context.read<ProfilBloc>().add(const LoadProfil());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(AppStrings.profil),
        actions: [
          TextButton(
            onPressed: () {
              showDialog(
                context: context,
                builder: (ctx) => AlertDialog(
                  title: const Text('Keluar'),
                  content: const Text('Apakah Anda yakin ingin keluar?'),
                  actions: [
                    TextButton(
                      onPressed: () => Navigator.pop(ctx),
                      child: const Text('Batal'),
                    ),
                    ElevatedButton(
                      style: ElevatedButton.styleFrom(
                        backgroundColor: AppColors.error,
                        minimumSize: Size.zero,
                        padding: const EdgeInsets.symmetric(
                          horizontal: AppSizes.md,
                          vertical: AppSizes.sm,
                        ),
                      ),
                      onPressed: () {
                        Navigator.pop(ctx);
                        context.read<AuthBloc>().add(const LogoutRequested());
                      },
                      child: const Text(AppStrings.logout),
                    ),
                  ],
                ),
              );
            },
            child: const Text(
              AppStrings.logout,
              style: TextStyle(color: Colors.white),
            ),
          ),
        ],
      ),
      body: BlocBuilder<ProfilBloc, ProfilState>(
        builder: (context, state) {
          if (state is ProfilLoading) {
            return const Center(child: CircularProgressIndicator());
          }

          if (state is ProfilError) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.error_outline,
                      size: 48, color: AppColors.error),
                  const SizedBox(height: AppSizes.md),
                  Text(state.message),
                  const SizedBox(height: AppSizes.md),
                  ElevatedButton(
                    onPressed: () =>
                        context.read<ProfilBloc>().add(const LoadProfil()),
                    child: const Text(AppStrings.retry),
                  ),
                ],
              ),
            );
          }

          if (state is ProfilLoaded) {
            final p = state.profil;
            return SingleChildScrollView(
              child: Column(
                children: [
                  // Header
                  Container(
                    width: double.infinity,
                    padding: const EdgeInsets.all(AppSizes.xl),
                    color: AppColors.primary,
                    child: Column(
                      children: [
                        CircleAvatar(
                          radius: 40,
                          backgroundColor: Colors.white24,
                          child: Text(
                            p.nama.isNotEmpty
                                ? p.nama[0].toUpperCase()
                                : 'N',
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 32,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                        const SizedBox(height: AppSizes.md),
                        Text(
                          p.nama,
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: AppSizes.xs),
                        Text(
                          p.nomorNasabah,
                          style: const TextStyle(
                            color: Colors.white70,
                            fontSize: 14,
                          ),
                        ),
                      ],
                    ),
                  ),

                  // Info list
                  Padding(
                    padding: const EdgeInsets.all(AppSizes.pagePadding),
                    child: Card(
                      child: Padding(
                        padding: const EdgeInsets.all(AppSizes.cardPadding),
                        child: Column(
                          children: [
                            _ProfilRow(
                              icon: Icons.badge_outlined,
                              label: AppStrings.nomorKTP,
                              value: p.nik ?? '-',
                            ),
                            const Divider(),
                            _ProfilRow(
                              icon: Icons.phone_outlined,
                              label: 'Telepon',
                              value: p.telepon ?? '-',
                            ),
                            const Divider(),
                            _ProfilRow(
                              icon: Icons.email_outlined,
                              label: 'Email',
                              value: p.email ?? '-',
                            ),
                            const Divider(),
                            _ProfilRow(
                              icon: Icons.location_on_outlined,
                              label: 'Alamat',
                              value: p.alamat ?? '-',
                            ),
                            if (p.tanggalLahir != null) ...[
                              const Divider(),
                              _ProfilRow(
                                icon: Icons.cake_outlined,
                                label: 'Tanggal Lahir',
                                value: formatTanggal(p.tanggalLahir!),
                              ),
                            ],
                          ],
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

class _ProfilRow extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;

  const _ProfilRow({
    required this.icon,
    required this.label,
    required this.value,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: AppSizes.sm),
      child: Row(
        children: [
          Icon(icon, color: AppColors.primary, size: AppSizes.iconMd),
          const SizedBox(width: AppSizes.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textSecondary,
                  ),
                ),
                Text(
                  value,
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
