import 'package:flutter/material.dart';
import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_strings.dart';
import '../../bmt/presentation/pages/bmt_page.dart';
import '../../pecahan_uang/presentation/pages/pecahan_uang_page.dart';
import '../../platform_settings/presentation/pages/platform_settings_page.dart';
import '../../usage_log/presentation/pages/usage_log_page.dart';

class MainShell extends StatefulWidget {
  final VoidCallback onLogout;
  const MainShell({super.key, required this.onLogout});

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> {
  int _selectedIndex = 0;

  final List<_NavItem> _navItems = const [
    _NavItem(icon: Icons.corporate_fare, label: AppStrings.navBmt),
    _NavItem(icon: Icons.payments_outlined, label: AppStrings.navPecahan),
    _NavItem(icon: Icons.tune, label: AppStrings.navSettings),
    _NavItem(icon: Icons.analytics_outlined, label: AppStrings.navUsageLog),
  ];

  final List<Widget> _pages = const [
    BmtPage(),
    PecahanUangPage(),
    PlatformSettingsPage(),
    UsageLogPage(),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Row(
        children: [
          _DevSidebar(
            items: _navItems,
            selectedIndex: _selectedIndex,
            onItemTap: (i) => setState(() => _selectedIndex = i),
            onLogout: widget.onLogout,
          ),
          Expanded(child: _pages[_selectedIndex]),
        ],
      ),
    );
  }
}

class _NavItem {
  final IconData icon;
  final String label;
  const _NavItem({required this.icon, required this.label});
}

class _DevSidebar extends StatelessWidget {
  final List<_NavItem> items;
  final int selectedIndex;
  final ValueChanged<int> onItemTap;
  final VoidCallback onLogout;

  const _DevSidebar({
    required this.items,
    required this.selectedIndex,
    required this.onItemTap,
    required this.onLogout,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 220,
      color: AppColors.sidebarBg,
      child: Column(
        children: [
          Container(
            height: 72,
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: Row(
              children: [
                const Icon(Icons.terminal, color: AppColors.accent, size: 28),
                const SizedBox(width: 12),
                Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('BMT', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
                    Text('Developer', style: TextStyle(color: Colors.white.withOpacity(0.5), fontSize: 11)),
                  ],
                ),
              ],
            ),
          ),
          const Divider(color: Colors.white12, height: 1),
          const SizedBox(height: 8),
          Expanded(
            child: ListView.builder(
              itemCount: items.length,
              itemBuilder: (_, i) {
                final selected = selectedIndex == i;
                return Container(
                  margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 2),
                  decoration: BoxDecoration(
                    color: selected ? AppColors.sidebarActive.withOpacity(0.4) : null,
                    borderRadius: BorderRadius.circular(6),
                  ),
                  child: ListTile(
                    dense: true,
                    leading: Icon(
                      items[i].icon,
                      color: selected ? AppColors.accent : AppColors.sidebarText,
                      size: 20,
                    ),
                    title: Text(
                      items[i].label,
                      style: TextStyle(
                        color: selected ? AppColors.accent : AppColors.sidebarText,
                        fontSize: 13,
                        fontWeight: selected ? FontWeight.w600 : FontWeight.normal,
                      ),
                    ),
                    onTap: () => onItemTap(i),
                  ),
                );
              },
            ),
          ),
          const Divider(color: Colors.white12, height: 1),
          ListTile(
            dense: true,
            leading: const Icon(Icons.logout, color: Colors.red, size: 20),
            title: const Text('Logout', style: TextStyle(color: Colors.red, fontSize: 13)),
            onTap: onLogout,
          ),
          const SizedBox(height: 8),
        ],
      ),
    );
  }
}
