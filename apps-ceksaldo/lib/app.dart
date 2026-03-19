import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'features/ceksaldo/presentation/bloc/kiosk_bloc.dart';
import 'features/ceksaldo/presentation/pages/kiosk_page.dart';

class AppCekSaldo extends StatelessWidget {
  const AppCekSaldo({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => KioskBloc(),
      child: MaterialApp(
        title: 'Cek Saldo',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
          useMaterial3: true,
          colorScheme: ColorScheme.fromSeed(
            seedColor: const Color(0xFF1A237E),
            brightness: Brightness.dark,
          ),
          scaffoldBackgroundColor: const Color(0xFF0A1628),
        ),
        home: const KioskPage(),
      ),
    );
  }
}
