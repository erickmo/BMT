import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../core/constants/app_colors.dart';
import '../../auth/domain/entities/user_entity.dart';
import '../../auth/presentation/bloc/auth_bloc.dart';
import '../../dashboard/presentation/pages/dashboard_page.dart';
import '../../nasabah/presentation/pages/nasabah_list_page.dart';
import '../../rekening/presentation/pages/rekening_list_page.dart';
import '../../laporan/presentation/pages/laporan_page.dart';
import '../../akuntansi/presentation/pages/akuntansi_page.dart';

// ─── Menu definition ──────────────────────────────────────────────────────────

class _Menu {
  final String label;
  final IconData icon;
  final Widget page;
  const _Menu({required this.label, required this.icon, required this.page});
}

const _menus = [
  _Menu(
    label: 'Dashboard',
    icon: Icons.dashboard_outlined,
    page: DashboardPage(),
  ),
  _Menu(
    label: 'Rekening',
    icon: Icons.account_balance_wallet_outlined,
    page: RekeningListPage(),
  ),
  _Menu(
    label: 'Nasabah',
    icon: Icons.people_outline_rounded,
    page: NasabahListPage(),
  ),
  _Menu(
    label: 'Laporan',
    icon: Icons.bar_chart_rounded,
    page: LaporanPage(),
  ),
  _Menu(
    label: 'Akuntansi',
    icon: Icons.account_balance_outlined,
    page: AkuntansiPage(),
  ),
];

// ─── MainShell ────────────────────────────────────────────────────────────────

class MainShell extends StatefulWidget {
  final UserEntity user;
  const MainShell({super.key, required this.user});

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> {
  int _activeIndex = 0;

  void _onMenuTap(int index) => setState(() => _activeIndex = index);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: Column(
        children: [
          _TopNavbar(user: widget.user),
          _MenuNavbar(activeIndex: _activeIndex, onTap: _onMenuTap),
          const Divider(height: 1, thickness: 1, color: Color(0xFFE5E7EB)),
          Expanded(child: _menus[_activeIndex].page),
        ],
      ),
    );
  }
}

// ─── Navbar 1: Brand + Breadcrumb + Actions ───────────────────────────────────

class _TopNavbar extends StatelessWidget {
  final UserEntity user;
  const _TopNavbar({required this.user});

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 56,
      padding: const EdgeInsets.symmetric(horizontal: 20),
      decoration: const BoxDecoration(
        color: Colors.white,
        border: Border(bottom: BorderSide(color: Color(0xFFE5E7EB), width: 1)),
      ),
      child: Row(
        children: [
          // Brand
          Row(
            children: [
              Container(
                width: 28,
                height: 28,
                decoration: BoxDecoration(
                  color: AppColors.primary,
                  borderRadius: BorderRadius.circular(7),
                ),
                child: const Icon(Icons.account_balance_rounded,
                    color: Colors.white, size: 15),
              ),
              const SizedBox(width: 10),
              const Text(
                'BMT',
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w700,
                  color: AppColors.textPrimary,
                  letterSpacing: 0.5,
                ),
              ),
              const SizedBox(width: 6),
              const Text(
                '/',
                style: TextStyle(color: Color(0xFFD1D5DB), fontSize: 14),
              ),
              const SizedBox(width: 6),
              const Text(
                'Management',
                style: TextStyle(
                  fontSize: 13,
                  color: AppColors.textSecondary,
                  fontWeight: FontWeight.w400,
                ),
              ),
            ],
          ),

          const Spacer(),

          // Notification
          IconButton(
            onPressed: () {},
            icon: const Icon(Icons.notifications_none_rounded,
                color: AppColors.textSecondary, size: 20),
            tooltip: 'Notifikasi',
            style: IconButton.styleFrom(
              minimumSize: const Size(36, 36),
              padding: EdgeInsets.zero,
            ),
          ),

          const SizedBox(width: 4),

          // User avatar + dropdown
          _UserMenu(user: user),
        ],
      ),
    );
  }
}

class _UserMenu extends StatelessWidget {
  final UserEntity user;
  const _UserMenu({required this.user});

  @override
  Widget build(BuildContext context) {
    final initials = user.nama.isNotEmpty
        ? user.nama.trim().split(' ').map((w) => w[0]).take(2).join().toUpperCase()
        : 'U';

    return PopupMenuButton<String>(
      offset: const Offset(0, 40),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(10),
        side: const BorderSide(color: Color(0xFFE5E7EB)),
      ),
      elevation: 4,
      shadowColor: Colors.black12,
      onSelected: (value) {
        if (value == 'logout') {
          context.read<AuthBloc>().add(const LogoutRequested());
        }
      },
      itemBuilder: (_) => [
        PopupMenuItem(
          enabled: false,
          padding: const EdgeInsets.fromLTRB(16, 12, 16, 8),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                user.nama,
                style: const TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              const SizedBox(height: 2),
              Text(
                user.username,
                style: const TextStyle(
                    fontSize: 11, color: AppColors.textSecondary),
              ),
            ],
          ),
        ),
        const PopupMenuDivider(height: 1),
        const PopupMenuItem(
          value: 'profile',
          padding: EdgeInsets.symmetric(horizontal: 16, vertical: 10),
          child: Row(
            children: [
              Icon(Icons.person_outline_rounded,
                  size: 16, color: AppColors.textSecondary),
              SizedBox(width: 10),
              Text('My Profile',
                  style: TextStyle(fontSize: 13, color: AppColors.textPrimary)),
            ],
          ),
        ),
        const PopupMenuItem(
          value: 'password',
          padding: EdgeInsets.symmetric(horizontal: 16, vertical: 10),
          child: Row(
            children: [
              Icon(Icons.lock_outline_rounded,
                  size: 16, color: AppColors.textSecondary),
              SizedBox(width: 10),
              Text('Change Password',
                  style: TextStyle(fontSize: 13, color: AppColors.textPrimary)),
            ],
          ),
        ),
        const PopupMenuDivider(height: 1),
        const PopupMenuItem(
          value: 'logout',
          padding: EdgeInsets.symmetric(horizontal: 16, vertical: 10),
          child: Row(
            children: [
              Icon(Icons.logout_rounded, size: 16, color: Colors.red),
              SizedBox(width: 10),
              Text('Logout',
                  style: TextStyle(fontSize: 13, color: Colors.red)),
            ],
          ),
        ),
      ],
      child: Row(
        children: [
          Container(
            width: 30,
            height: 30,
            decoration: BoxDecoration(
              color: AppColors.primaryPale,
              borderRadius: BorderRadius.circular(8),
              border: Border.all(color: AppColors.border),
            ),
            child: Center(
              child: Text(
                initials,
                style: const TextStyle(
                  fontSize: 11,
                  fontWeight: FontWeight.w700,
                  color: AppColors.primary,
                ),
              ),
            ),
          ),
          const SizedBox(width: 6),
          const Icon(Icons.keyboard_arrow_down_rounded,
              size: 16, color: AppColors.textSecondary),
        ],
      ),
    );
  }
}

// ─── Navbar 2: Menu tabs ──────────────────────────────────────────────────────

class _MenuNavbar extends StatelessWidget {
  final int activeIndex;
  final void Function(int) onTap;

  const _MenuNavbar({required this.activeIndex, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 44,
      color: Colors.white,
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Row(
        children: List.generate(_menus.length, (i) {
          final menu = _menus[i];
          final isActive = activeIndex == i;
          return _MenuTab(
            label: menu.label,
            icon: menu.icon,
            isActive: isActive,
            onTap: () => onTap(i),
          );
        }),
      ),
    );
  }
}

class _MenuTab extends StatelessWidget {
  final String label;
  final IconData icon;
  final bool isActive;
  final VoidCallback onTap;

  const _MenuTab({
    required this.label,
    required this.icon,
    required this.isActive,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        margin: const EdgeInsets.only(right: 4),
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
        decoration: BoxDecoration(
          color: isActive ? AppColors.primaryPale : Colors.transparent,
          borderRadius: BorderRadius.circular(6),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              icon,
              size: 14,
              color: isActive ? AppColors.primary : AppColors.textSecondary,
            ),
            const SizedBox(width: 6),
            Text(
              label,
              style: TextStyle(
                fontSize: 13,
                fontWeight:
                    isActive ? FontWeight.w600 : FontWeight.w400,
                color: isActive ? AppColors.primary : AppColors.textSecondary,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
