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
        if (_countdown <= 0) t.cancel();
      }
    });
    _autoResetTimer = Timer(const Duration(seconds: 10), () {
      if (mounted) context.read<KioskBloc>().add(const KioskReset());
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

  void _simulateNfcTap() {
    context.read<KioskBloc>().add(const NfcTagDetected('04:AB:CD:EF:12:34:56'));
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
              if (state is KioskIdle)
                _IdleScreen(onSimulateTap: _simulateNfcTap),
              if (state is KioskLoading) const _LoadingScreen(),
              if (state is KioskShowSaldo)
                _SaldoScreen(
                  saldoInfo: state.saldoInfo,
                  countdown: _countdown,
                  onReset: () =>
                      context.read<KioskBloc>().add(const KioskReset()),
                ),
              if (state is KioskError)
                _ErrorScreen(
                  message: state.message,
                  countdown: _countdown,
                  onReset: () =>
                      context.read<KioskBloc>().add(const KioskReset()),
                ),
            ],
          ),
        );
      },
    );
  }
}

// ── Idle Screen ─────────────────────────────────────────────────────────────

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
        decoration: const BoxDecoration(
          gradient: RadialGradient(
            center: Alignment.center,
            radius: 1.4,
            colors: [Color(0xFF0F3D22), Color(0xFF052E16)],
          ),
        ),
        child: Stack(
          children: [
            // Background pattern circles
            _BgCircle(size: 600, opacity: 0.04, offset: Offset.zero),
            _BgCircle(size: 400, opacity: 0.05, offset: Offset.zero),
            _BgCircle(size: 200, opacity: 0.07, offset: Offset.zero),
            // Content
            Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                // Logo + name
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Container(
                      width: 36,
                      height: 36,
                      decoration: BoxDecoration(
                        color: const Color(0xFF1A7A4A).withOpacity(0.4),
                        borderRadius: BorderRadius.circular(10),
                        border: Border.all(
                            color: const Color(0xFF4ADE80).withOpacity(0.3),
                            width: 1),
                      ),
                      child: const Icon(Icons.account_balance_rounded,
                          color: Color(0xFF4ADE80), size: 18),
                    ),
                    const SizedBox(width: 12),
                    const Text(
                      'BMT',
                      style: TextStyle(
                        color: Color(0xFF6EE7B7),
                        fontSize: 20,
                        letterSpacing: 6,
                        fontWeight: FontWeight.w300,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 64),

                // NFC animation
                AnimatedBuilder(
                  animation: _pulseAnim,
                  builder: (_, __) {
                    return Stack(
                      alignment: Alignment.center,
                      children: [
                        // Outer glow ring
                        Container(
                          width: 260 * _pulseAnim.value,
                          height: 260 * _pulseAnim.value,
                          decoration: BoxDecoration(
                            shape: BoxShape.circle,
                            border: Border.all(
                              color: AppColors.nfcRing.withOpacity(
                                0.12 * (2 - _pulseAnim.value),
                              ),
                              width: 1.5,
                            ),
                          ),
                        ),
                        // Middle ring
                        Container(
                          width: 190 * _pulseAnim.value,
                          height: 190 * _pulseAnim.value,
                          decoration: BoxDecoration(
                            shape: BoxShape.circle,
                            border: Border.all(
                              color: AppColors.nfcRing.withOpacity(
                                0.22 * (2 - _pulseAnim.value),
                              ),
                              width: 1.5,
                            ),
                          ),
                        ),
                        // Inner glow
                        Container(
                          width: 120,
                          height: 120,
                          decoration: BoxDecoration(
                            shape: BoxShape.circle,
                            color: AppColors.nfcRing.withOpacity(0.06),
                            border: Border.all(
                              color: AppColors.nfcRing.withOpacity(0.3),
                              width: 1.5,
                            ),
                            boxShadow: [
                              BoxShadow(
                                color: AppColors.nfcRing.withOpacity(0.15),
                                blurRadius: 30,
                                spreadRadius: 5,
                              ),
                            ],
                          ),
                          child: const Icon(
                            Icons.nfc_rounded,
                            size: 56,
                            color: AppColors.nfcIcon,
                          ),
                        ),
                      ],
                    );
                  },
                ),
                const SizedBox(height: 60),

                // Instructions
                const Text(
                  'Tempelkan Kartu NFC',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 32,
                    fontWeight: FontWeight.bold,
                    letterSpacing: -0.5,
                  ),
                ),
                const SizedBox(height: 12),
                const Text(
                  'untuk melihat saldo rekening Anda',
                  style: TextStyle(
                    color: Color(0xFF6EE7B7),
                    fontSize: 18,
                  ),
                ),
                const SizedBox(height: 56),

                // Demo hint
                Container(
                  padding: const EdgeInsets.symmetric(
                      horizontal: 20, vertical: 10),
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.04),
                    borderRadius: BorderRadius.circular(30),
                    border: Border.all(
                        color: Colors.white.withOpacity(0.08), width: 1),
                  ),
                  child: const Text(
                    '[ Tap layar untuk simulasi ]',
                    style: TextStyle(
                      color: Color(0xFF2EA878),
                      fontSize: 13,
                      letterSpacing: 0.5,
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _BgCircle extends StatelessWidget {
  final double size;
  final double opacity;
  final Offset offset;
  const _BgCircle(
      {required this.size, required this.opacity, required this.offset});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Container(
        width: size,
        height: size,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          border: Border.all(
            color: const Color(0xFF4ADE80).withOpacity(opacity),
            width: 1,
          ),
        ),
      ),
    );
  }
}

// ── Loading Screen ───────────────────────────────────────────────────────────

class _LoadingScreen extends StatelessWidget {
  const _LoadingScreen();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      height: double.infinity,
      decoration: const BoxDecoration(
        gradient: RadialGradient(
          center: Alignment.center,
          radius: 1.4,
          colors: [Color(0xFF0F3D22), Color(0xFF052E16)],
        ),
      ),
      child: const Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          CircularProgressIndicator(
            color: AppColors.nfcRing,
            strokeWidth: 3,
          ),
          SizedBox(height: 28),
          Text(
            'Membaca kartu...',
            style: TextStyle(
              color: Color(0xFF6EE7B7),
              fontSize: 20,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}

// ── Saldo Screen ─────────────────────────────────────────────────────────────

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
    return LayoutBuilder(
      builder: (context, constraints) {
        final isMobile = constraints.maxWidth < 720;
        return Container(
          width: double.infinity,
          height: double.infinity,
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              colors: [Color(0xFF052E16), Color(0xFF0A3D20)],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
          ),
          child: isMobile
              ? _buildMobileLayout(context)
              : _buildDesktopLayout(context),
        );
      },
    );
  }

  Widget _buildMobileLayout(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.fromLTRB(20, 48, 20, 32),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildSaldoCard(mobile: true),
          const SizedBox(height: 20),
          _buildTransaksiPanel(mobile: true),
        ],
      ),
    );
  }

  Widget _buildDesktopLayout(BuildContext context) {
    return Row(
      children: [
        Expanded(
          flex: 2,
          child: Padding(
            padding: const EdgeInsets.all(56),
            child: _buildSaldoCard(mobile: false),
          ),
        ),
        Expanded(
          flex: 3,
          child: Padding(
            padding: const EdgeInsets.fromLTRB(0, 40, 40, 40),
            child: _buildTransaksiPanel(mobile: false),
          ),
        ),
      ],
    );
  }

  Widget _buildSaldoCard({required bool mobile}) {
    return Column(
      mainAxisAlignment:
          mobile ? MainAxisAlignment.start : MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Avatar + nama
        Row(
          children: [
            Container(
              width: mobile ? 48 : 64,
              height: mobile ? 48 : 64,
              decoration: BoxDecoration(
                color: const Color(0xFF1A7A4A).withOpacity(0.3),
                shape: BoxShape.circle,
                border: Border.all(
                    color: AppColors.nfcRing.withOpacity(0.3), width: 1.5),
              ),
              child: Icon(Icons.person_rounded,
                  size: mobile ? 24 : 32, color: AppColors.nfcRing),
            ),
            const SizedBox(width: 14),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    saldoInfo.namaNasabah,
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: mobile ? 18 : 30,
                      fontWeight: FontWeight.bold,
                      letterSpacing: -0.5,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 3),
                  Text(
                    saldoInfo.nomorRekening,
                    style: const TextStyle(
                      color: Color(0xFF6EE7B7),
                      fontSize: 13,
                      fontFamily: 'monospace',
                      letterSpacing: 1,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
        SizedBox(height: mobile ? 24 : 40),

        // Saldo
        Container(
          width: double.infinity,
          padding: mobile
              ? const EdgeInsets.all(20)
              : EdgeInsets.zero,
          decoration: mobile
              ? BoxDecoration(
                  color: Colors.white.withOpacity(0.05),
                  borderRadius: BorderRadius.circular(20),
                  border: Border.all(
                      color: Colors.white.withOpacity(0.08), width: 1),
                )
              : null,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'S A L D O',
                style: TextStyle(
                  color: Color(0xFF6EE7B7),
                  fontSize: 12,
                  letterSpacing: 4,
                  fontWeight: FontWeight.w300,
                ),
              ),
              const SizedBox(height: 8),
              ShaderMask(
                shaderCallback: (bounds) => const LinearGradient(
                  colors: [Color(0xFF4ADE80), Color(0xFFC9A84C)],
                  begin: Alignment.centerLeft,
                  end: Alignment.centerRight,
                ).createShader(bounds),
                child: Text(
                  formatRupiah(saldoInfo.saldo),
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: mobile ? 34 : 52,
                    fontWeight: FontWeight.bold,
                    letterSpacing: -1,
                  ),
                ),
              ),
            ],
          ),
        ),
        SizedBox(height: mobile ? 20 : 48),

        // Countdown / reset
        GestureDetector(
          onTap: onReset,
          child: Container(
            padding:
                const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.05),
              borderRadius: BorderRadius.circular(30),
              border: Border.all(
                  color: Colors.white.withOpacity(0.15), width: 1),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(Icons.replay_rounded,
                    color: Color(0xFF6EE7B7), size: 15),
                const SizedBox(width: 7),
                Text(
                  'Reset dalam ${countdown}s',
                  style: const TextStyle(
                    color: Color(0xFF6EE7B7),
                    fontSize: 13,
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildTransaksiPanel({required bool mobile}) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(20),
        border:
            Border.all(color: Colors.white.withOpacity(0.08), width: 1),
      ),
      child: Column(
        mainAxisSize: mobile ? MainAxisSize.min : MainAxisSize.max,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 18, 20, 14),
            child: Row(
              children: [
                Container(
                  width: 7,
                  height: 7,
                  decoration: BoxDecoration(
                    color: AppColors.nfcRing,
                    shape: BoxShape.circle,
                    boxShadow: [
                      BoxShadow(
                        color: AppColors.nfcRing.withOpacity(0.5),
                        blurRadius: 6,
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 10),
                const Text(
                  '5 Transaksi Terakhir',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 15,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
            ),
          ),
          Container(height: 1, color: Colors.white.withOpacity(0.06)),
          // List transaksi
          if (saldoInfo.transaksiTerakhir.isEmpty)
            Padding(
              padding: const EdgeInsets.all(24),
              child: Center(
                child: Text(
                  'Belum ada transaksi',
                  style: TextStyle(
                      color: Colors.white.withOpacity(0.4), fontSize: 14),
                ),
              ),
            )
          else
            ...List.generate(saldoInfo.transaksiTerakhir.length, (i) {
              final tx = saldoInfo.transaksiTerakhir[i];
              final isLast = i == saldoInfo.transaksiTerakhir.length - 1;
              return Column(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 20, vertical: 14),
                    child: Row(
                      children: [
                        Container(
                          width: 36,
                          height: 36,
                          decoration: BoxDecoration(
                            color: tx.isKredit
                                ? AppColors.positif.withOpacity(0.12)
                                : AppColors.negatif.withOpacity(0.12),
                            borderRadius: BorderRadius.circular(10),
                          ),
                          child: Icon(
                            tx.isKredit
                                ? Icons.arrow_downward_rounded
                                : Icons.arrow_upward_rounded,
                            color: tx.isKredit
                                ? AppColors.positif
                                : AppColors.negatif,
                            size: 17,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                tx.keterangan,
                                style: const TextStyle(
                                  color: Colors.white,
                                  fontWeight: FontWeight.w500,
                                  fontSize: 13,
                                ),
                                maxLines: 1,
                                overflow: TextOverflow.ellipsis,
                              ),
                              const SizedBox(height: 2),
                              Text(
                                formatTanggalWaktu(tx.tanggal),
                                style: const TextStyle(
                                  color: Color(0xFF6EE7B7),
                                  fontSize: 11,
                                ),
                              ),
                            ],
                          ),
                        ),
                        const SizedBox(width: 8),
                        Text(
                          '${tx.isKredit ? '+' : '-'} ${formatRupiah(tx.nominal)}',
                          style: TextStyle(
                            color: tx.isKredit
                                ? AppColors.positif
                                : AppColors.negatif,
                            fontWeight: FontWeight.bold,
                            fontSize: mobile ? 13 : 15,
                          ),
                        ),
                      ],
                    ),
                  ),
                  if (!isLast)
                    Container(
                        height: 1,
                        color: Colors.white.withOpacity(0.05)),
                ],
              );
            }),
        ],
      ),
    );
  }
}

// ── Error Screen ─────────────────────────────────────────────────────────────

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
        decoration: const BoxDecoration(
          gradient: RadialGradient(
            center: Alignment.center,
            radius: 1.4,
            colors: [Color(0xFF1A0A0A), Color(0xFF0A0505)],
          ),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              width: 100,
              height: 100,
              decoration: BoxDecoration(
                color: AppColors.error.withOpacity(0.1),
                shape: BoxShape.circle,
                border: Border.all(
                    color: AppColors.error.withOpacity(0.3), width: 1.5),
              ),
              child: const Icon(Icons.error_outline_rounded,
                  size: 50, color: AppColors.error),
            ),
            const SizedBox(height: 32),
            Text(
              message,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 22,
                fontWeight: FontWeight.w600,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 40),
            Container(
              padding: const EdgeInsets.symmetric(
                  horizontal: 24, vertical: 12),
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.05),
                borderRadius: BorderRadius.circular(30),
                border:
                    Border.all(color: Colors.white.withOpacity(0.1), width: 1),
              ),
              child: Text(
                'Tap untuk kembali ($countdown)',
                style: const TextStyle(
                  color: Color(0xFF6EE7B7),
                  fontSize: 15,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
