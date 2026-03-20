import 'package:flutter/material.dart';
import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_strings.dart';
import '../../dashboard/presentation/pages/dashboard_page.dart';
import '../../nasabah/presentation/pages/nasabah_list_page.dart';
import '../../rekening/presentation/pages/rekening_list_page.dart';
import '../../form/presentation/pages/form_list_page.dart';
import '../../settings/presentation/pages/settings_page.dart';

class MainShell extends StatefulWidget {
  const MainShell({super.key});

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> {
  int _selectedIndex = 0;

  final List<_NavItem> _navItems = const [
    _NavItem(
      icon: Icons.dashboard_outlined,
      activeIcon: Icons.dashboard_rounded,
      label: AppStrings.navDashboard,
    ),
    _NavItem(
      icon: Icons.people_outline_rounded,
      activeIcon: Icons.people_rounded,
      label: AppStrings.navNasabah,
    ),
    _NavItem(
      icon: Icons.account_balance_wallet_outlined,
      activeIcon: Icons.account_balance_wallet_rounded,
      label: AppStrings.navRekening,
    ),
    _NavItem(
      icon: Icons.description_outlined,
      activeIcon: Icons.description_rounded,
      label: AppStrings.navForm,
    ),
    _NavItem(
      icon: Icons.settings_outlined,
      activeIcon: Icons.settings_rounded,
      label: AppStrings.navSettings,
    ),
  ];

  final List<Widget> _pages = const [
    DashboardPage(),
    NasabahListPage(),
    RekeningListPage(),
    FormListPage(),
    SettingsPage(),
  ];

  @override
  Widget build(BuildContext context) {
    final isWide = MediaQuery.of(context).size.width > 700;

    return Scaffold(
      body: Row(
        children: [
          if (isWide)
            _EmeraldSidebar(
              items: _navItems,
              selectedIndex: _selectedIndex,
              onItemTap: (i) => setState(() => _selectedIndex = i),
              appTitle: 'BMT',
              appSubtitle: 'Management',
              headerIcon: Icons.business_rounded,
            ),
          Expanded(child: _pages[_selectedIndex]),
        ],
      ),
      bottomNavigationBar: isWide
          ? null
          : _buildBottomNav(),
    );
  }

  Widget _buildBottomNav() {
    return Container(
      decoration: const BoxDecoration(
        color: AppColors.surface,
        border: Border(top: BorderSide(color: AppColors.border, width: 1)),
      ),
      child: NavigationBar(
        backgroundColor: AppColors.surface,
        elevation: 0,
        selectedIndex: _selectedIndex,
        onDestinationSelected: (i) => setState(() => _selectedIndex = i),
        destinations: _navItems
            .map((item) => NavigationDestination(
                  icon: Icon(item.icon),
                  selectedIcon: Icon(item.activeIcon),
                  label: item.label,
                ))
            .toList(),
      ),
    );
  }
}

class _NavItem {
  final IconData icon;
  final IconData activeIcon;
  final String label;
  const _NavItem({
    required this.icon,
    required this.activeIcon,
    required this.label,
  });
}

/// Reusable elegant sidebar for management & pondok apps
class _EmeraldSidebar extends StatelessWidget {
  final List<_NavItem> items;
  final int selectedIndex;
  final ValueChanged<int> onItemTap;
  final String appTitle;
  final String appSubtitle;
  final IconData headerIcon;
  final VoidCallback? onLogout;

  const _EmeraldSidebar({
    required this.items,
    required this.selectedIndex,
    required this.onItemTap,
    required this.appTitle,
    required this.appSubtitle,
    required this.headerIcon,
    this.onLogout,
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
        border: Border(
          right: BorderSide(color: Color(0x1A4ADE80), width: 1),
        ),
      ),
      child: Column(
        children: [
          // ── Logo / Brand ──────────────────────────────────────────────────
          _SidebarHeader(
            title: appTitle,
            subtitle: appSubtitle,
            icon: headerIcon,
          ),
          // ── Divider ───────────────────────────────────────────────────────
          _SidebarDivider(),
          const SizedBox(height: 8),
          // ── Nav Items ─────────────────────────────────────────────────────
          Expanded(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10),
              child: ListView.separated(
                itemCount: items.length,
                separatorBuilder: (_, __) => const SizedBox(height: 2),
                itemBuilder: (_, i) {
                  final item = items[i];
                  final selected = selectedIndex == i;
                  return _SidebarItem(
                    icon: selected ? item.activeIcon : item.icon,
                    label: item.label,
                    selected: selected,
                    onTap: () => onItemTap(i),
                  );
                },
              ),
            ),
          ),
          // ── Bottom ────────────────────────────────────────────────────────
          _SidebarDivider(),
          Padding(
            padding: const EdgeInsets.fromLTRB(10, 8, 10, 20),
            child: _SidebarItem(
              icon: Icons.logout_rounded,
              label: AppStrings.logout,
              selected: false,
              onTap: onLogout ?? () {},
              isDestructive: true,
            ),
          ),
        ],
      ),
    );
  }
}

class _SidebarHeader extends StatelessWidget {
  final String title;
  final String subtitle;
  final IconData icon;

  const _SidebarHeader({
    required this.title,
    required this.subtitle,
    required this.icon,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 28, 20, 20),
      child: Row(
        children: [
          Container(
            width: 42,
            height: 42,
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [Color(0xFF1A7A4A), Color(0xFF2EA878)],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: BorderRadius.circular(11),
              boxShadow: [
                BoxShadow(
                  color: const Color(0xFF1A7A4A).withOpacity(0.4),
                  blurRadius: 8,
                  offset: const Offset(0, 3),
                ),
              ],
            ),
            child: Icon(icon, color: Colors.white, size: 22),
          ),
          const SizedBox(width: 12),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                title,
                style: const TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                  fontSize: 15,
                ),
              ),
              Text(
                subtitle,
                style: const TextStyle(
                  color: Color(0xFF6EE7B7),
                  fontSize: 11,
                  letterSpacing: 0.3,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _SidebarDivider extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      height: 1,
      margin: const EdgeInsets.symmetric(horizontal: 20),
      color: Colors.white.withOpacity(0.07),
    );
  }
}

class _SidebarItem extends StatelessWidget {
  final IconData icon;
  final String label;
  final bool selected;
  final VoidCallback onTap;
  final bool isDestructive;

  const _SidebarItem({
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
      iconColor = const Color(0xFF6EE7B7).withOpacity(0.7);
      textColor = const Color(0xFF6EE7B7).withOpacity(0.7);
    }

    return Container(
      decoration: BoxDecoration(
        color: selected
            ? const Color(0xFF1A7A4A).withOpacity(0.22)
            : null,
        borderRadius: BorderRadius.circular(10),
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
                // Active indicator bar
                AnimatedContainer(
                  duration: const Duration(milliseconds: 200),
                  width: 3,
                  height: selected ? 20 : 0,
                  decoration: BoxDecoration(
                    color: const Color(0xFF4ADE80),
                    borderRadius: BorderRadius.circular(2),
                  ),
                  margin: EdgeInsets.only(right: selected ? 10 : 0),
                ),
                if (!selected) const SizedBox(width: 13),
                Icon(icon, color: iconColor, size: 18),
                const SizedBox(width: 10),
                Expanded(
                  child: Text(
                    label,
                    style: TextStyle(
                      color: textColor,
                      fontSize: 13,
                      fontWeight:
                          selected ? FontWeight.w600 : FontWeight.normal,
                    ),
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
