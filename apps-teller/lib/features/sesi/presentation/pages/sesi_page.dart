import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';
import '../../domain/entities/sesi_entity.dart';
import '../bloc/sesi_bloc.dart';

class SesiPage extends StatefulWidget {
  const SesiPage({super.key});

  @override
  State<SesiPage> createState() => _SesiPageState();
}

class _SesiPageState extends State<SesiPage> {
  @override
  void initState() {
    super.initState();
    context.read<SesiBloc>().add(const LoadSesiAktif());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('Sesi Teller'),
      ),
      body: BlocConsumer<SesiBloc, SesiState>(
        listener: (context, state) {
          if (state is SesiBukaSuccess) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: const Text('Sesi berhasil dibuka'),
                backgroundColor: AppColors.success,
                behavior: SnackBarBehavior.floating,
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10)),
              ),
            );
            context.read<SesiBloc>().add(const LoadSesiAktif());
          }
          if (state is SesiTutupSuccess) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: const Text('Sesi berhasil ditutup'),
                backgroundColor: AppColors.success,
                behavior: SnackBarBehavior.floating,
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10)),
              ),
            );
            context.read<SesiBloc>().add(const LoadSesiAktif());
          }
          if (state is SesiError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: AppColors.error,
                behavior: SnackBarBehavior.floating,
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10)),
              ),
            );
          }
        },
        builder: (context, state) {
          if (state is SesiLoading) {
            return const Center(
              child: CircularProgressIndicator(color: AppColors.primary),
            );
          }

          if (state is SesiAktifLoaded) {
            if (state.sesi == null) return _BukaSesiView();
            return _SesiAktifView(sesi: state.sesi!);
          }

          if (state is PecahanLoaded) {
            return _PecahanInputView(state: state);
          }

          return const Center(child: Text(AppStrings.noSesiAktif));
        },
      ),
    );
  }
}

// ── Buka Sesi View ───────────────────────────────────────────────────────────

class _BukaSesiView extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSizes.pagePadding),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              width: 96,
              height: 96,
              decoration: BoxDecoration(
                gradient: const LinearGradient(
                  colors: [Color(0xFF0C4A2B), Color(0xFF1A7A4A)],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
                shape: BoxShape.circle,
                boxShadow: [
                  BoxShadow(
                    color: AppColors.primary.withOpacity(0.25),
                    blurRadius: 24,
                    offset: const Offset(0, 8),
                  ),
                ],
              ),
              child: const Icon(Icons.lock_open_rounded,
                  size: 44, color: Colors.white),
            ),
            const SizedBox(height: AppSizes.xl),
            const Text(
              AppStrings.noSesiAktif,
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
                color: AppColors.textPrimary,
              ),
            ),
            const SizedBox(height: 8),
            const Text(
              'Buka sesi untuk mulai bertransaksi.\nSemua tombol transaksi dinonaktifkan tanpa sesi.',
              textAlign: TextAlign.center,
              style: TextStyle(
                color: AppColors.textSecondary,
                height: 1.5,
              ),
            ),
            const SizedBox(height: AppSizes.xl),
            _EmeraldButton(
              icon: Icons.lock_open_rounded,
              label: AppStrings.bukaSesi,
              onPressed: () =>
                  context.read<SesiBloc>().add(const LoadPecahanAktif()),
            ),
          ],
        ),
      ),
    );
  }
}

// ── Sesi Aktif View ──────────────────────────────────────────────────────────

class _SesiAktifView extends StatelessWidget {
  final SesiEntity sesi;
  const _SesiAktifView({required this.sesi});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(AppSizes.pagePadding),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // ── Sesi card ──────────────────────────────────────────────────────
          Container(
            padding: const EdgeInsets.all(24),
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [Color(0xFF0C4A2B), Color(0xFF1A7A4A)],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: BorderRadius.circular(20),
              boxShadow: [
                BoxShadow(
                  color: AppColors.primary.withOpacity(0.25),
                  blurRadius: 24,
                  offset: const Offset(0, 8),
                ),
              ],
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 12, vertical: 5),
                      decoration: BoxDecoration(
                        color: Colors.white.withOpacity(0.15),
                        borderRadius: BorderRadius.circular(20),
                        border: Border.all(
                            color: Colors.white.withOpacity(0.25), width: 1),
                      ),
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Container(
                            width: 7,
                            height: 7,
                            decoration: const BoxDecoration(
                              color: Color(0xFF4ADE80),
                              shape: BoxShape.circle,
                            ),
                          ),
                          const SizedBox(width: 6),
                          const Text(
                            'SESI AKTIF',
                            style: TextStyle(
                              color: Colors.white,
                              fontWeight: FontWeight.bold,
                              fontSize: 11,
                              letterSpacing: 0.5,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 20),
                Text(
                  'Modal Awal',
                  style: TextStyle(
                      color: Colors.white.withOpacity(0.7), fontSize: 13),
                ),
                const SizedBox(height: 4),
                Text(
                  formatRupiah(sesi.modalAwal),
                  style: const TextStyle(
                    fontSize: 32,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                    letterSpacing: -0.5,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  'Dibuka: ${formatTanggalWaktu(sesi.dibukaPada)}',
                  style: TextStyle(
                      color: Colors.white.withOpacity(0.6), fontSize: 13),
                ),
              ],
            ),
          ),
          const SizedBox(height: AppSizes.xl),

          const Text(
            'Redenominasi Awal',
            style: TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 15,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: AppSizes.sm),

          Expanded(
            child: ListView.separated(
              itemCount: sesi.pecahan.length,
              separatorBuilder: (_, __) =>
                  const Divider(height: 1, color: AppColors.divider),
              itemBuilder: (context, index) {
                final p = sesi.pecahan[index];
                return Container(
                  padding: const EdgeInsets.symmetric(
                      vertical: 12, horizontal: 4),
                  child: Row(
                    children: [
                      Container(
                        width: 40,
                        height: 40,
                        decoration: BoxDecoration(
                          color: AppColors.primaryPale,
                          borderRadius: BorderRadius.circular(10),
                        ),
                        child: const Icon(Icons.payments_rounded,
                            color: AppColors.primary, size: 20),
                      ),
                      const SizedBox(width: 14),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(p.label,
                                style: const TextStyle(
                                    fontWeight: FontWeight.w500,
                                    color: AppColors.textPrimary)),
                            Text('${p.jumlah} lembar/koin',
                                style: const TextStyle(
                                    fontSize: 12,
                                    color: AppColors.textSecondary)),
                          ],
                        ),
                      ),
                      Text(
                        formatRupiah(p.subtotal),
                        style: const TextStyle(
                          fontWeight: FontWeight.bold,
                          color: AppColors.primary,
                        ),
                      ),
                    ],
                  ),
                );
              },
            ),
          ),
          const SizedBox(height: AppSizes.md),

          // ── Tutup Sesi Button ──────────────────────────────────────────────
          SizedBox(
            width: double.infinity,
            height: 50,
            child: OutlinedButton.icon(
              icon: const Icon(Icons.lock_rounded, color: AppColors.error),
              label: const Text(
                AppStrings.tutupSesi,
                style: TextStyle(color: AppColors.error),
              ),
              style: OutlinedButton.styleFrom(
                side: const BorderSide(color: AppColors.error, width: 1.5),
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12)),
              ),
              onPressed: () => _showTutupSesiDialog(context, sesi.id),
            ),
          ),
        ],
      ),
    );
  }

  void _showTutupSesiDialog(BuildContext context, String sesiId) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        title: const Text(AppStrings.tutupSesi),
        content: const Text(
          'Apakah Anda yakin ingin menutup sesi?\nPastikan kas sudah dihitung.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text(AppStrings.cancel),
          ),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.error,
              minimumSize: Size.zero,
              padding: const EdgeInsets.symmetric(
                  horizontal: AppSizes.md, vertical: AppSizes.sm),
            ),
            onPressed: () {
              Navigator.pop(ctx);
              context.read<SesiBloc>().add(
                    TutupSesiRequested(sesiId: sesiId, pecahanAkhir: const []),
                  );
            },
            child: const Text(AppStrings.tutupSesi),
          ),
        ],
      ),
    );
  }
}

// ── Pecahan Input View ───────────────────────────────────────────────────────

class _PecahanInputView extends StatefulWidget {
  final PecahanLoaded state;
  const _PecahanInputView({required this.state});

  @override
  State<_PecahanInputView> createState() => _PecahanInputViewState();
}

class _PecahanInputViewState extends State<_PecahanInputView> {
  late Map<String, int> _jumlahMap;

  @override
  void initState() {
    super.initState();
    _jumlahMap = Map.from(widget.state.jumlahMap);
  }

  int get _totalModal {
    return widget.state.pecahan.fold(
      0,
      (sum, p) => sum + p.nominal * (_jumlahMap[p.id] ?? 0),
    );
  }

  void _buka() {
    final pecahanList = widget.state.pecahan
        .where((p) => (_jumlahMap[p.id] ?? 0) > 0)
        .map((p) => {'pecahan_id': p.id, 'jumlah': _jumlahMap[p.id]})
        .toList();
    context.read<SesiBloc>().add(BukaSesiRequested(pecahan: pecahanList));
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        // ── Total banner ────────────────────────────────────────────────────
        Container(
          padding: const EdgeInsets.all(AppSizes.md),
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              colors: [Color(0xFF0C4A2B), Color(0xFF1A7A4A)],
              begin: Alignment.centerLeft,
              end: Alignment.centerRight,
            ),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text('Total Modal Awal',
                  style: TextStyle(color: Colors.white70, fontSize: 13)),
              Text(
                formatRupiah(_totalModal),
                style: const TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                  fontSize: 20,
                ),
              ),
            ],
          ),
        ),

        // ── Pecahan list ────────────────────────────────────────────────────
        Expanded(
          child: ListView.separated(
            padding: const EdgeInsets.all(AppSizes.md),
            itemCount: widget.state.pecahan.length,
            separatorBuilder: (_, __) =>
                const Divider(height: 1, color: AppColors.divider),
            itemBuilder: (context, index) {
              final pecahan = widget.state.pecahan[index];
              final jumlah = _jumlahMap[pecahan.id] ?? 0;

              return Padding(
                padding: const EdgeInsets.symmetric(vertical: 8),
                child: Row(
                  children: [
                    Container(
                      width: 40,
                      height: 40,
                      decoration: BoxDecoration(
                        color: AppColors.primaryPale,
                        borderRadius: BorderRadius.circular(10),
                      ),
                      child: const Icon(Icons.payments_rounded,
                          color: AppColors.primary, size: 20),
                    ),
                    const SizedBox(width: 14),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(pecahan.label,
                              style: const TextStyle(
                                  fontWeight: FontWeight.w500,
                                  color: AppColors.textPrimary)),
                          Text(
                            'Subtotal: ${formatRupiah(pecahan.nominal * jumlah)}',
                            style: const TextStyle(
                                color: AppColors.textSecondary, fontSize: 12),
                          ),
                        ],
                      ),
                    ),
                    SizedBox(
                      width: 80,
                      child: TextFormField(
                        initialValue: jumlah.toString(),
                        keyboardType: TextInputType.number,
                        inputFormatters: [
                          FilteringTextInputFormatter.digitsOnly
                        ],
                        textAlign: TextAlign.center,
                        decoration: const InputDecoration(
                          contentPadding: EdgeInsets.symmetric(
                              vertical: 8, horizontal: 8),
                        ),
                        onChanged: (val) => setState(() {
                          _jumlahMap[pecahan.id] = int.tryParse(val) ?? 0;
                        }),
                      ),
                    ),
                  ],
                ),
              );
            },
          ),
        ),

        // ── Buka Sesi button ────────────────────────────────────────────────
        Padding(
          padding: const EdgeInsets.all(AppSizes.pagePadding),
          child: SizedBox(
            width: double.infinity,
            height: 52,
            child: ElevatedButton.icon(
              icon: const Icon(Icons.lock_open_rounded),
              label: const Text(AppStrings.bukaSesi),
              onPressed: _totalModal > 0 ? _buka : null,
            ),
          ),
        ),
      ],
    );
  }
}

class _EmeraldButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onPressed;

  const _EmeraldButton({
    required this.icon,
    required this.label,
    required this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 52,
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF0C4A2B), Color(0xFF1A7A4A)],
          begin: Alignment.centerLeft,
          end: Alignment.centerRight,
        ),
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: AppColors.primary.withOpacity(0.3),
            blurRadius: 12,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onPressed,
          borderRadius: BorderRadius.circular(12),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(icon, color: Colors.white, size: 20),
              const SizedBox(width: 8),
              Text(
                label,
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 15,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
