import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/utils/formatters.dart';
import '../../domain/entities/saldo_info_entity.dart';
import '../bloc/kiosk_bloc.dart';

class KioskPage extends StatefulWidget {
  const KioskPage({super.key});

  @override
  State<KioskPage> createState() => _KioskPageState();
}

class _KioskPageState extends State<KioskPage> {
  Timer? _autoResetTimer;
  int _countdown = 10;
  Timer? _countdownTimer;

  void _startAutoReset() {
    _autoResetTimer?.cancel();
    _countdownTimer?.cancel();
    _countdown = 10;
    _countdownTimer = Timer.periodic(const Duration(seconds: 1), (t) {
      if (mounted) {
        setState(() => _countdown--);
        if (_countdown <= 0) {
          t.cancel();
        }
      }
    });
    _autoResetTimer = Timer(const Duration(seconds: 10), () {
      if (mounted) {
        context.read<KioskBloc>().add(const KioskReset());
      }
    });
  }

  void _cancelAutoReset() {
    _autoResetTimer?.cancel();
    _countdownTimer?.cancel();
    _countdown = 10;
  }

  @override
  void dispose() {
    _autoResetTimer?.cancel();
    _countdownTimer?.cancel();
    super.dispose();
  }

  // Simulate NFC tap for demo (real app uses nfc_manager package)
  void _simulateNfcTap() {
    context
        .read<KioskBloc>()
        .add(const NfcTagDetected('04:AB:CD:EF:12:34:56'));
  }

  @override
  Widget build(BuildContext context) {
    return BlocConsumer<KioskBloc, KioskState>(
      listener: (context, state) {
        if (state is KioskShowSaldo || state is KioskError) {
          _startAutoReset();
        } else if (state is KioskIdle) {
          _cancelAutoReset();
        }
      },
      builder: (context, state) {
        return Scaffold(
          backgroundColor: AppColors.background,
          body: Stack(
            children: [
              if (state is KioskIdle) _IdleScreen(onSimulateTap: _simulateNfcTap),
              if (state is KioskLoading) const _LoadingScreen(),
              if (state is KioskShowSaldo)
                _SaldoScreen(
                  saldoInfo: state.saldoInfo,
                  countdown: _countdown,
                  onReset: () => context.read<KioskBloc>().add(const KioskReset()),
                ),
              if (state is KioskError)
                _ErrorScreen(
                  message: state.message,
                  countdown: _countdown,
                  onReset: () => context.read<KioskBloc>().add(const KioskReset()),
                ),
            ],
          ),
        );
      },
    );
  }
}

class _IdleScreen extends StatefulWidget {
  final VoidCallback onSimulateTap;
  const _IdleScreen({required this.onSimulateTap});

  @override
  State<_IdleScreen> createState() => _IdleScreenState();
}

class _IdleScreenState extends State<_IdleScreen>
    with SingleTickerProviderStateMixin {
  late AnimationController _animCtrl;
  late Animation<double> _pulseAnim;

  @override
  void initState() {
    super.initState();
    _animCtrl = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 2),
    )..repeat(reverse: true);
    _pulseAnim = Tween<double>(begin: 0.85, end: 1.0).animate(
      CurvedAnimation(parent: _animCtrl, curve: Curves.easeInOut),
    );
  }

  @override
  void dispose() {
    _animCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: widget.onSimulateTap,
      child: Container(
        width: double.infinity,
        height: double.infinity,
        color: AppColors.background,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            // BMT Logo area
            const Text(
              'BMT',
              style: TextStyle(
                color: AppColors.textSecondary,
                fontSize: 18,
                letterSpacing: 4,
              ),
            ),
            const SizedBox(height: 48),
            // NFC animation
            AnimatedBuilder(
              animation: _pulseAnim,
              builder: (_, __) {
                return Stack(
                  alignment: Alignment.center,
                  children: [
                    // Outer ring
                    Container(
                      width: 220 * _pulseAnim.value,
                      height: 220 * _pulseAnim.value,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        border: Border.all(
                          color: AppColors.nfcRing.withOpacity(
                              0.2 * (2 - _pulseAnim.value)),
                          width: 2,
                        ),
                      ),
                    ),
                    // Middle ring
                    Container(
                      width: 160 * _pulseAnim.value,
                      height: 160 * _pulseAnim.value,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        border: Border.all(
                          color: AppColors.nfcRing.withOpacity(
                              0.4 * (2 - _pulseAnim.value)),
                          width: 2,
                        ),
                      ),
                    ),
                    // Icon
                    const Icon(Icons.nfc, size: 100, color: AppColors.nfcIcon),
                  ],
                );
              },
            ),
            const SizedBox(height: 48),
            const Text(
              'Tempelkan Kartu NFC',
              style: TextStyle(
                color: AppColors.textPrimary,
                fontSize: 28,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            const Text(
              'untuk melihat saldo rekening Anda',
              style: TextStyle(
                color: AppColors.textSecondary,
                fontSize: 16,
              ),
            ),
            const SizedBox(height: 60),
            const Text(
              'Tap layar untuk simulasi (demo)',
              style: TextStyle(color: AppColors.textHint, fontSize: 12),
            ),
          ],
        ),
      ),
    );
  }
}

class _LoadingScreen extends StatelessWidget {
  const _LoadingScreen();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      height: double.infinity,
      color: AppColors.background,
      child: const Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          CircularProgressIndicator(color: AppColors.nfcRing, strokeWidth: 3),
          SizedBox(height: 24),
          Text(
            'Membaca kartu...',
            style: TextStyle(color: AppColors.textSecondary, fontSize: 18),
          ),
        ],
      ),
    );
  }
}

class _SaldoScreen extends StatelessWidget {
  final SaldoInfoEntity saldoInfo;
  final int countdown;
  final VoidCallback onReset;

  const _SaldoScreen({
    required this.saldoInfo,
    required this.countdown,
    required this.onReset,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      height: double.infinity,
      color: AppColors.background,
      child: Row(
        children: [
          // Left: Saldo info
          Expanded(
            flex: 2,
            child: Padding(
              padding: const EdgeInsets.all(48),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Icon(Icons.person, size: 48, color: AppColors.nfcRing),
                  const SizedBox(height: 16),
                  Text(
                    saldoInfo.namaNasabah,
                    style: const TextStyle(
                      color: AppColors.textPrimary,
                      fontSize: 28,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    saldoInfo.nomorRekening,
                    style: const TextStyle(
                      color: AppColors.textSecondary,
                      fontSize: 13,
                      fontFamily: 'monospace',
                    ),
                  ),
                  const SizedBox(height: 32),
                  const Text(
                    'SALDO',
                    style: TextStyle(
                      color: AppColors.textSecondary,
                      fontSize: 14,
                      letterSpacing: 2,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    formatRupiah(saldoInfo.saldo),
                    style: const TextStyle(
                      color: AppColors.nfcIcon,
                      fontSize: 48,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 48),
                  // Auto reset countdown
                  GestureDetector(
                    onTap: onReset,
                    child: Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 24, vertical: 12),
                      decoration: BoxDecoration(
                        border: Border.all(color: AppColors.textHint),
                        borderRadius: BorderRadius.circular(30),
                      ),
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          const Icon(Icons.replay,
                              color: AppColors.textSecondary, size: 16),
                          const SizedBox(width: 8),
                          Text(
                            'Reset otomatis dalam ${countdown}s',
                            style: const TextStyle(
                                color: AppColors.textSecondary, fontSize: 13),
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
          // Right: Transaction history
          Expanded(
            flex: 3,
            child: Container(
              margin: const EdgeInsets.all(32),
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: AppColors.surface,
                borderRadius: BorderRadius.circular(16),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    '5 Transaksi Terakhir',
                    style: TextStyle(
                      color: AppColors.textPrimary,
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 16),
                  Expanded(
                    child: ListView.separated(
                      itemCount: saldoInfo.transaksiTerakhir.length,
                      separatorBuilder: (_, __) => Divider(
                          color: AppColors.textHint.withOpacity(0.2),
                          height: 1),
                      itemBuilder: (_, i) {
                        final tx = saldoInfo.transaksiTerakhir[i];
                        return Padding(
                          padding: const EdgeInsets.symmetric(vertical: 14),
                          child: Row(
                            children: [
                              Icon(
                                tx.isKredit
                                    ? Icons.arrow_downward
                                    : Icons.arrow_upward,
                                color: tx.isKredit
                                    ? AppColors.positif
                                    : AppColors.negatif,
                                size: 20,
                              ),
                              const SizedBox(width: 12),
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(tx.keterangan,
                                        style: const TextStyle(
                                            color: AppColors.textPrimary,
                                            fontWeight: FontWeight.w500)),
                                    Text(
                                      formatTanggalWaktu(tx.tanggal),
                                      style: const TextStyle(
                                          color: AppColors.textSecondary,
                                          fontSize: 12),
                                    ),
                                  ],
                                ),
                              ),
                              Text(
                                '${tx.isKredit ? '+' : '-'} ${formatRupiah(tx.nominal)}',
                                style: TextStyle(
                                  color: tx.isKredit
                                      ? AppColors.positif
                                      : AppColors.negatif,
                                  fontWeight: FontWeight.bold,
                                  fontSize: 15,
                                ),
                              ),
                            ],
                          ),
                        );
                      },
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _ErrorScreen extends StatelessWidget {
  final String message;
  final int countdown;
  final VoidCallback onReset;

  const _ErrorScreen({
    required this.message,
    required this.countdown,
    required this.onReset,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onReset,
      child: Container(
        width: double.infinity,
        height: double.infinity,
        color: AppColors.background,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error_outline,
                size: 72, color: AppColors.error),
            const SizedBox(height: 24),
            Text(
              message,
              style: const TextStyle(
                  color: AppColors.textPrimary, fontSize: 20),
            ),
            const SizedBox(height: 32),
            Text(
              'Tap untuk kembali ($countdown)',
              style: const TextStyle(
                  color: AppColors.textSecondary, fontSize: 14),
            ),
          ],
        ),
      ),
    );
  }
}
