import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/storage/local_storage.dart';
import '../../../auth/presentation/bloc/auth_bloc.dart';
import '../../../sesi/presentation/bloc/sesi_bloc.dart';
import '../../../transaksi/presentation/bloc/transaksi_bloc.dart';
import '../../../transaksi/presentation/pages/transaksi_page.dart';
import '../../../sesi/presentation/pages/sesi_page.dart';
import '../../../../injection_container.dart';

class DashboardPage extends StatefulWidget {
  const DashboardPage({super.key});

  @override
  State<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends State<DashboardPage> {
  int _selectedIndex = 0;

  static const _navItems = [
    NavigationRailDestination(
      icon: Icon(Icons.swap_horiz),
      label: Text('Transaksi'),
    ),
    NavigationRailDestination(
      icon: Icon(Icons.calendar_today),
      label: Text('Sesi'),
    ),
  ];

  Widget _buildBody(BuildContext context) {
    switch (_selectedIndex) {
      case 0:
        return BlocProvider(
          create: (_) => sl<TransaksiBloc>(),
          child: const TransaksiPage(),
        );
      case 1:
        return BlocProvider(
          create: (_) => sl<SesiBloc>(),
          child: const SesiPage(),
        );
      default:
        return const SizedBox.shrink();
    }
  }

  @override
  Widget build(BuildContext context) {
    final localStorage = sl<LocalStorage>();
    return Scaffold(
      body: Row(
        children: [
          // Sidebar
          NavigationRail(
            backgroundColor: AppColors.sidebarBg,
            selectedIndex: _selectedIndex,
            onDestinationSelected: (i) => setState(() => _selectedIndex = i),
            labelType: NavigationRailLabelType.all,
            selectedIconTheme:
                const IconThemeData(color: AppColors.primary),
            unselectedIconTheme:
                const IconThemeData(color: AppColors.textOnSidebar),
            selectedLabelTextStyle:
                const TextStyle(color: AppColors.primary, fontSize: 12),
            unselectedLabelTextStyle: const TextStyle(
                color: AppColors.textOnSidebar, fontSize: 12),
            leading: Padding(
              padding: const EdgeInsets.symmetric(vertical: AppSizes.lg),
              child: Column(
                children: [
                  const Icon(
                    Icons.point_of_sale,
                    color: AppColors.primary,
                    size: 36,
                  ),
                  const SizedBox(height: AppSizes.sm),
                  Text(
                    localStorage.namaTeller ?? 'Teller',
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 12,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  Text(
                    localStorage.role ?? '',
                    style: const TextStyle(
                      color: Colors.white54,
                      fontSize: 11,
                    ),
                  ),
                ],
              ),
            ),
            trailing: Padding(
              padding: const EdgeInsets.only(bottom: AppSizes.lg),
              child: IconButton(
                icon: const Icon(Icons.logout, color: Colors.white54),
                tooltip: AppStrings.logout,
                onPressed: () {
                  context.read<AuthBloc>().add(const LogoutRequested());
                },
              ),
            ),
            destinations: _navItems,
          ),
          const VerticalDivider(thickness: 1, width: 1),
          // Main content
          Expanded(child: _buildBody(context)),
        ],
      ),
    );
  }
}
