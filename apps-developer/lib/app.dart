import 'package:flutter/material.dart';
import 'core/theme/app_theme.dart';
import 'features/auth/presentation/pages/login_page.dart';
import 'features/shared/widgets/main_shell.dart';
import 'injection_container.dart';
import 'core/network/dio_client.dart';

class AppDeveloper extends StatefulWidget {
  const AppDeveloper({super.key});

  @override
  State<AppDeveloper> createState() => _AppDeveloperState();
}

class _AppDeveloperState extends State<AppDeveloper> {
  bool _isLoggedIn = false;

  @override
  void initState() {
    super.initState();
    _checkLogin();
  }

  Future<void> _checkLogin() async {
    final hasToken = await sl<DioClient>().hasToken();
    if (mounted) setState(() => _isLoggedIn = hasToken);
  }

  Future<void> _onLogin(String token) async {
    await sl<DioClient>().saveToken(token);
    if (mounted) setState(() => _isLoggedIn = true);
  }

  Future<void> _onLogout() async {
    await sl<DioClient>().clearToken();
    if (mounted) setState(() => _isLoggedIn = false);
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'BMT Developer',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      home: _isLoggedIn
          ? MainShell(onLogout: _onLogout)
          : LoginPage(onLogin: _onLogin),
    );
  }
}
