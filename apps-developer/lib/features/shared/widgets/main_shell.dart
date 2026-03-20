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
    _NavItem(
        icon: Icons.corporate_fare_rounded,
        label: AppStrings.navBmt),
    _NavItem(
        icon: Icons.payments_outlined,
        label: AppStrings.navPecahan),
    _NavItem(
        icon: Icons.tune_rounded,
        label: AppStrings.navSettings),
    _NavItem(
        icon: Icons.analytics_outlined,
        label: AppStrings.navUsageLog),
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
      width: 240,
      decoration: const BoxDecoration(
        gradient: LinearGradient(
          colors: [Color(0xFF040F07), Color(0xFF071A0D)],
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
        ),
        border: Border(
          right: BorderSide(color: Color(0x1A34D399), width: 1),
        ),
      ),
      child: Column(
        children: [
          // ── Header ────────────────────────────────────────────────────────
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 28, 20, 20),
            child: Row(
              children: [
                Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: const Color(0xFF1A7A4A).withOpacity(0.3),
                    borderRadius: BorderRadius.circular(10),
                    border: Border.all(
                        color: const Color(0xFF34D399).withOpacity(0.3),
                        width: 1),
                  ),
                  child: const Icon(Icons.terminal_rounded,
                      color: Color(0xFF34D399), size: 20),
                ),
                const SizedBox(width: 12),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'BMT',
                      style: TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                        fontSize: 15,
                      ),
                    ),
                    Text(
                      'Developer',
                      style: TextStyle(
                        color: const Color(0xFF34D399).withOpacity(0.8),
                        fontSize: 11,
                        letterSpacing: 0.5,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          Container(
            height: 1,
            margin: const EdgeInsets.symmetric(horizontal: 20),
            color: Colors.white.withOpacity(0.06),
          ),
          const SizedBox(height: 8),
          // ── Nav ───────────────────────────────────────────────────────────
          Expanded(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10),
              child: ListView.separated(
                itemCount: items.length,
                separatorBuilder: (_, __) => const SizedBox(height: 2),
                itemBuilder: (_, i) {
                  final selected = selectedIndex == i;
                  return _DevNavItem(
                    icon: items[i].icon,
                    label: items[i].label,
                    selected: selected,
                    onTap: () => onItemTap(i),
                  );
                },
              ),
            ),
          ),
          Container(
            height: 1,
            margin: const EdgeInsets.symmetric(horizontal: 20),
            color: Colors.white.withOpacity(0.06),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(10, 8, 10, 20),
            child: _DevNavItem(
              icon: Icons.logout_rounded,
              label: 'Logout',
              selected: false,
              onTap: onLogout,
              isDestructive: true,
            ),
          ),
        ],
      ),
    );
  }
}

class _DevNavItem extends StatelessWidget {
  final IconData icon;
  final String label;
  final bool selected;
  final VoidCallback onTap;
  final bool isDestructive;

  const _DevNavItem({
    required this.icon,
    required this.label,
    required this.selected,
    required this.onTap,
    this.isDestructive = false,
  });

  @override
  Widget build(BuildContext context) {
    final Color iconColor;
    final Color textColor;

    if (isDestructive) {
      iconColor = const Color(0xFFFCA5A5);
      textColor = const Color(0xFFFCA5A5);
    } else if (selected) {
      iconColor = const Color(0xFF34D399);
      textColor = Colors.white;
    } else {
      iconColor = const Color(0xFF34D399).withOpacity(0.5);
      textColor = const Color(0xFF34D399).withOpacity(0.5);
    }

    return Container(
      decoration: BoxDecoration(
        color: selected ? const Color(0xFF1A7A4A).withOpacity(0.2) : null,
        borderRadius: BorderRadius.circular(10),
        border: selected
            ? Border.all(
                color: const Color(0xFF34D399).withOpacity(0.15), width: 1)
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
                AnimatedContainer(
                  duration: const Duration(milliseconds: 200),
                  width: 3,
                  height: selected ? 20 : 0,
                  decoration: BoxDecoration(
                    color: const Color(0xFF34D399),
                    borderRadius: BorderRadius.circular(2),
                  ),
                  margin: EdgeInsets.only(right: selected ? 10 : 0),
                ),
                if (!selected) const SizedBox(width: 13),
                Icon(icon, color: iconColor, size: 18),
                const SizedBox(width: 10),
                Text(
                  label,
                  style: TextStyle(
                    color: textColor,
                    fontSize: 13,
                    fontWeight: selected ? FontWeight.w600 : FontWeight.normal,
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
