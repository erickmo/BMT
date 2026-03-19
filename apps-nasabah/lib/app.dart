import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'core/theme/app_theme.dart';
import 'features/auth/presentation/bloc/auth_bloc.dart';
import 'features/auth/presentation/pages/login_page.dart';
import 'features/home/presentation/bloc/home_bloc.dart';
import 'features/home/presentation/pages/home_page.dart';
import 'features/profile/presentation/bloc/profil_bloc.dart';
import 'features/profile/presentation/pages/profil_page.dart';
import 'features/rekening/presentation/bloc/rekening_bloc.dart';
import 'features/rekening/presentation/pages/rekening_detail_page.dart';
import 'features/rekening/presentation/pages/rekening_list_page.dart';
import 'injection_container.dart';

class AppNasabahRoot extends StatefulWidget {
  const AppNasabahRoot({super.key});

  @override
  State<AppNasabahRoot> createState() => _AppNasabahRootState();
}

class _AppNasabahRootState extends State<AppNasabahRoot> {
  late final GoRouter _router;

  @override
  void initState() {
    super.initState();
    _router = GoRouter(
      initialLocation: '/login',
      routes: [
        GoRoute(
          path: '/login',
          builder: (context, state) => BlocProvider(
            create: (_) => sl<AuthBloc>(),
            child: const LoginPage(),
          ),
        ),
        GoRoute(
          path: '/home',
          builder: (context, state) => MultiBlocProvider(
            providers: [
              BlocProvider(create: (_) => sl<AuthBloc>()),
              BlocProvider(create: (_) => sl<HomeBloc>()),
            ],
            child: const HomePage(),
          ),
        ),
        GoRoute(
          path: '/rekening',
          builder: (context, state) => BlocProvider(
            create: (_) => sl<RekeningBloc>(),
            child: const RekeningListPage(),
          ),
        ),
        GoRoute(
          path: '/rekening/:id',
          builder: (context, state) {
            final id = state.pathParameters['id']!;
            return BlocProvider(
              create: (_) => sl<RekeningBloc>(),
              child: RekeningDetailPage(rekeningId: id),
            );
          },
        ),
        GoRoute(
          path: '/profil',
          builder: (context, state) => MultiBlocProvider(
            providers: [
              BlocProvider(create: (_) => sl<AuthBloc>()),
              BlocProvider(create: (_) => sl<ProfilBloc>()),
            ],
            child: const ProfilPage(),
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'BMT Nasabah',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light,
      routerConfig: _router,
    );
  }
}
