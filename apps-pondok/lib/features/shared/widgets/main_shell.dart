import 'package:flutter/material.dart';
import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_strings.dart';
import '../../santri/presentation/pages/santri_list_page.dart';
import '../../absensi/presentation/pages/absensi_page.dart';
import '../../jadwal/presentation/pages/jadwal_page.dart';
import '../../nilai/presentation/pages/nilai_page.dart';
import '../../tagihan/presentation/pages/tagihan_page.dart';

class MainShell extends StatefulWidget {
  const MainShell({super.key});

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> {
  int _selectedIndex = 0;

  final List<_NavItem> _navItems = const [
    _NavItem(icon: Icons.people_outline, label: AppStrings.navSantri),
    _NavItem(icon: Icons.how_to_reg_outlined, label: AppStrings.navAbsensi),
    _NavItem(icon: Icons.schedule_outlined, label: AppStrings.navJadwal),
    _NavItem(icon: Icons.grade_outlined, label: AppStrings.navNilai),
    _NavItem(icon: Icons.receipt_long_outlined, label: AppStrings.navTagihan),
  ];

  final List<Widget> _pages = const [
    SantriListPage(),
    AbsensiPage(),
    JadwalPage(),
    NilaiPage(),
    TagihanPage(),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _pages[_selectedIndex],
      bottomNavigationBar: NavigationBar(
        selectedIndex: _selectedIndex,
        onDestinationSelected: (i) => setState(() => _selectedIndex = i),
        destinations: _navItems
            .map((item) => NavigationDestination(
                  icon: Icon(item.icon),
                  label: item.label,
                ))
            .toList(),
      ),
    );
  }
}

class _NavItem {
  final IconData icon;
  final String label;
  const _NavItem({required this.icon, required this.label});
}
