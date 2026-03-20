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
    _NavDef(
      icon: Icons.swap_horiz_rounded,
      activeIcon: Icons.swap_horiz_rounded,
      label: 'Transaksi',
    ),
    _NavDef(
      icon: Icons.calendar_today_outlined,
      activeIcon: Icons.calendar_today_rounded,
      label: 'Sesi',
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
          // ── Elegant Sidebar ───────────────────────────────────────────────
          _TellerSidebar(
            items: _navItems,
            selectedIndex: _selectedIndex,
            onItemSelected: (i) => setState(() => _selectedIndex = i),
            teller: localStorage.namaTeller ?? 'Teller',
            role: localStorage.role ?? '',
            onLogout: () =>
                context.read<AuthBloc>().add(const LogoutRequested()),
          ),
          // ── Main Content ──────────────────────────────────────────────────
          Expanded(child: _buildBody(context)),
        ],
      ),
    );
  }
}

class _NavDef {
  final IconData icon;
  final IconData activeIcon;
  final String label;
  const _NavDef({
    required this.icon,
    required this.activeIcon,
    required this.label,
  });
}

class _TellerSidebar extends StatelessWidget {
  final List<_NavDef> items;
  final int selectedIndex;
  final ValueChanged<int> onItemSelected;
  final String teller;
  final String role;
  final VoidCallback onLogout;

  const _TellerSidebar({
    required this.items,
    required this.selectedIndex,
    required this.onItemSelected,
    required this.teller,
    required this.role,
    required this.onLogout,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 240,
      decoration: const BoxDecoration(
        gradient: LinearGradient(
          colors: [Color(0xFF052E16), Color(0xFF0A3D20)],
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
        ),
      ),
      child: Column(
        children: [
          // ── Header ───────────────────────────────────────────────────────
          Container(
            padding: const EdgeInsets.fromLTRB(20, 32, 20, 24),
            child: Column(
              children: [
                // Logo
                Container(
                  width: 52,
                  height: 52,
                  decoration: BoxDecoration(
                    gradient: const LinearGradient(
                      colors: [Color(0xFF1A7A4A), Color(0xFF2EA878)],
                      begin: Alignment.topLeft,
                      end: Alignment.bottomRight,
                    ),
                    borderRadius: BorderRadius.circular(14),
                    boxShadow: [
                      BoxShadow(
                        color: const Color(0xFF1A7A4A).withOpacity(0.4),
                        blurRadius: 12,
                        offset: const Offset(0, 4),
                      ),
                    ],
                  ),
                  child: const Icon(Icons.point_of_sale_rounded,
                      color: Colors.white, size: 26),
                ),
                const SizedBox(height: 16),
                // App name
                const Text(
                  AppStrings.appName,
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 15,
                    fontWeight: FontWeight.bold,
                    letterSpacing: 0.3,
                  ),
                ),
                const SizedBox(height: 2),
                const Text(
                  'Teller Portal',
                  style: TextStyle(
                    color: Color(0xFF6EE7B7),
                    fontSize: 11,
                    letterSpacing: 0.5,
                  ),
                ),
              ],
            ),
          ),

          // ── Divider ───────────────────────────────────────────────────────
          Container(
            height: 1,
            margin: const EdgeInsets.symmetric(horizontal: 20),
            color: Colors.white.withOpacity(0.08),
          ),
          const SizedBox(height: 8),

          // ── Nav Items ─────────────────────────────────────────────────────
          Expanded(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12),
              child: Column(
                children: List.generate(items.length, (i) {
                  final item = items[i];
                  final selected = selectedIndex == i;
                  return _SidebarNavItem(
                    icon: selected ? item.activeIcon : item.icon,
                    label: item.label,
                    selected: selected,
                    onTap: () => onItemSelected(i),
                  );
                }),
              ),
            ),
          ),

          // ── Divider ───────────────────────────────────────────────────────
          Container(
            height: 1,
            margin: const EdgeInsets.symmetric(horizontal: 20),
            color: Colors.white.withOpacity(0.08),
          ),

          // ── User + Logout ─────────────────────────────────────────────────
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 16, 20, 24),
            child: Column(
              children: [
                // User info
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.05),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                        color: Colors.white.withOpacity(0.08), width: 1),
                  ),
                  child: Row(
                    children: [
                      Container(
                        width: 36,
                        height: 36,
                        decoration: BoxDecoration(
                          color: const Color(0xFF1A7A4A).withOpacity(0.4),
                          borderRadius: BorderRadius.circular(10),
                        ),
                        child: Center(
                          child: Text(
                            teller.isNotEmpty
                                ? teller[0].toUpperCase()
                                : 'T',
                            style: const TextStyle(
                              color: Color(0xFF6EE7B7),
                              fontWeight: FontWeight.bold,
                              fontSize: 16,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 10),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              teller,
                              style: const TextStyle(
                                color: Colors.white,
                                fontSize: 13,
                                fontWeight: FontWeight.w600,
                              ),
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                            ),
                            Text(
                              role,
                              style: TextStyle(
                                color: Colors.white.withOpacity(0.5),
                                fontSize: 11,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 8),
                // Logout
                _SidebarNavItem(
                  icon: Icons.logout_rounded,
                  label: AppStrings.logout,
                  selected: false,
                  onTap: onLogout,
                  isDestructive: true,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _SidebarNavItem extends StatelessWidget {
  final IconData icon;
  final String label;
  final bool selected;
  final VoidCallback onTap;
  final bool isDestructive;

  const _SidebarNavItem({
    required this.icon,
    required this.label,
    required this.selected,
    required this.onTap,
    this.isDestructive = false,
  });

  @override
  Widget build(BuildContext context) {
    Color iconColor;
    Color textColor;

    if (isDestructive) {
      iconColor = const Color(0xFFFCA5A5);
      textColor = const Color(0xFFFCA5A5);
    } else if (selected) {
      iconColor = const Color(0xFF4ADE80);
      textColor = Colors.white;
    } else {
      iconColor = const Color(0xFF6EE7B7);
      textColor = const Color(0xFF6EE7B7);
    }

    return Container(
      margin: const EdgeInsets.symmetric(vertical: 2),
      decoration: BoxDecoration(
        color: selected ? const Color(0xFF1A7A4A).withOpacity(0.25) : null,
        borderRadius: BorderRadius.circular(10),
        border: selected
            ? Border.all(
                color: const Color(0xFF1A7A4A).withOpacity(0.4), width: 1)
            : null,
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(10),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
            child: Row(
              children: [
                if (selected)
                  Container(
                    width: 3,
                    height: 20,
                    decoration: BoxDecoration(
                      color: const Color(0xFF4ADE80),
                      borderRadius: BorderRadius.circular(2),
                    ),
                    margin: const EdgeInsets.only(right: 10),
                  )
                else
                  const SizedBox(width: 13),
                Icon(icon, color: iconColor, size: 18),
                const SizedBox(width: 10),
                Text(
                  label,
                  style: TextStyle(
                    color: textColor,
                    fontSize: 13,
                    fontWeight:
                        selected ? FontWeight.w600 : FontWeight.normal,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
