import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'core/theme/app_theme.dart';
import 'features/auth/presentation/bloc/auth_bloc.dart';
import 'features/auth/presentation/pages/login_page.dart';
import 'features/dashboard/presentation/pages/dashboard_page.dart';
import 'injection_container.dart';

const _mockLogin = bool.fromEnvironment('MOCK_LOGIN');

class AppTellerRoot extends StatefulWidget {
  const AppTellerRoot({super.key});

  @override
  State<AppTellerRoot> createState() => _AppTellerRootState();
}

class _AppTellerRootState extends State<AppTellerRoot> {
  late final GoRouter _router;

  @override
  void initState() {
    super.initState();
    _router = GoRouter(
      initialLocation: _mockLogin ? '/dashboard' : '/login',
      routes: [
        GoRoute(
          path: '/login',
          builder: (context, state) => BlocProvider(
            create: (_) => sl<AuthBloc>(),
            child: const LoginPage(),
          ),
        ),
        GoRoute(
          path: '/dashboard',
          builder: (context, state) => BlocProvider(
            create: (_) => sl<AuthBloc>(),
            child: const DashboardPage(),
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'BMT Teller',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light,
      routerConfig: _router,
    );
  }
}
