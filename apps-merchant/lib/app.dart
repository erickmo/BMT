import 'package:flutter/material.dart';
import 'core/constants/app_colors.dart';
import 'core/constants/app_strings.dart';
import 'core/theme/app_theme.dart';
import 'features/auth/presentation/pages/login_page.dart';
import 'features/kasir/presentation/pages/kasir_page.dart';
import 'features/laporan/presentation/pages/laporan_page.dart';

class AppMerchant extends StatefulWidget {
  const AppMerchant({super.key});

  @override
  State<AppMerchant> createState() => _AppMerchantState();
}

class _AppMerchantState extends State<AppMerchant> {
  bool _isLoggedIn = false;
  int _selectedIndex = 0;

  final List<Widget> _pages = const [
    KasirPage(),
    LaporanPage(),
  ];

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: AppStrings.appName,
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      home: _isLoggedIn ? _buildMain() : LoginPage(onLogin: () {
        setState(() => _isLoggedIn = true);
      }),
    );
  }

  Widget _buildMain() {
    return Scaffold(
      body: _pages[_selectedIndex],
      bottomNavigationBar: NavigationBar(
        selectedIndex: _selectedIndex,
        onDestinationSelected: (i) => setState(() => _selectedIndex = i),
        destinations: const [
          NavigationDestination(
            icon: Icon(Icons.point_of_sale_outlined),
            selectedIcon: Icon(Icons.point_of_sale),
            label: AppStrings.navKasir,
          ),
          NavigationDestination(
            icon: Icon(Icons.bar_chart_outlined),
            selectedIcon: Icon(Icons.bar_chart),
            label: AppStrings.navLaporan,
          ),
        ],
      ),
    );
  }
}
