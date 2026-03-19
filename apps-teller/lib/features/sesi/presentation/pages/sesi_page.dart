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
      appBar: AppBar(title: const Text('Sesi Teller')),
      body: BlocConsumer<SesiBloc, SesiState>(
        listener: (context, state) {
          if (state is SesiBukaSuccess) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Text('Sesi berhasil dibuka'),
                backgroundColor: AppColors.success,
              ),
            );
            context.read<SesiBloc>().add(const LoadSesiAktif());
          }
          if (state is SesiTutupSuccess) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Text('Sesi berhasil ditutup'),
                backgroundColor: AppColors.success,
              ),
            );
            context.read<SesiBloc>().add(const LoadSesiAktif());
          }
          if (state is SesiError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: AppColors.error,
              ),
            );
          }
        },
        builder: (context, state) {
          if (state is SesiLoading) {
            return const Center(child: CircularProgressIndicator());
          }

          if (state is SesiAktifLoaded) {
            if (state.sesi == null) {
              return _BukaSesiView();
            }
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

class _BukaSesiView extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSizes.pagePadding),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(
              Icons.lock_open,
              size: 72,
              color: AppColors.primary,
            ),
            const SizedBox(height: AppSizes.lg),
            const Text(
              AppStrings.noSesiAktif,
              style: TextStyle(
                fontSize: 18,
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: AppSizes.md),
            const Text(
              'Buka sesi untuk mulai bertransaksi.\nSemua tombol transaksi dinonaktifkan tanpa sesi.',
              textAlign: TextAlign.center,
              style: TextStyle(color: AppColors.textSecondary),
            ),
            const SizedBox(height: AppSizes.xl),
            ElevatedButton.icon(
              icon: const Icon(Icons.lock_open),
              label: const Text(AppStrings.bukaSesi),
              onPressed: () {
                context.read<SesiBloc>().add(const LoadPecahanAktif());
              },
            ),
          ],
        ),
      ),
    );
  }
}

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
          Card(
            child: Padding(
              padding: const EdgeInsets.all(AppSizes.cardPadding),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: AppSizes.sm,
                          vertical: AppSizes.xs,
                        ),
                        decoration: BoxDecoration(
                          color: AppColors.success.withValues(alpha: 0.15),
                          borderRadius:
                              BorderRadius.circular(AppSizes.radiusFull),
                        ),
                        child: const Text(
                          'SESI AKTIF',
                          style: TextStyle(
                            color: AppColors.success,
                            fontWeight: FontWeight.bold,
                            fontSize: 12,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSizes.md),
                  const Text(
                    'Modal Awal',
                    style: TextStyle(
                      color: AppColors.textSecondary,
                      fontSize: 13,
                    ),
                  ),
                  Text(
                    formatRupiah(sesi.modalAwal),
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                      color: AppColors.primary,
                    ),
                  ),
                  const SizedBox(height: AppSizes.sm),
                  Text(
                    'Dibuka: ${formatTanggalWaktu(sesi.dibukaPada)}',
                    style: const TextStyle(
                      color: AppColors.textSecondary,
                      fontSize: 13,
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: AppSizes.lg),
          const Text(
            'Redenominasi Awal',
            style: TextStyle(fontWeight: FontWeight.bold, fontSize: 15),
          ),
          const SizedBox(height: AppSizes.sm),
          Expanded(
            child: ListView.builder(
              itemCount: sesi.pecahan.length,
              itemBuilder: (context, index) {
                final p = sesi.pecahan[index];
                return ListTile(
                  leading: const Icon(Icons.money, color: AppColors.primary),
                  title: Text(p.label),
                  subtitle: Text('${p.jumlah} lembar/koin'),
                  trailing: Text(
                    formatRupiah(p.subtotal),
                    style: const TextStyle(fontWeight: FontWeight.bold),
                  ),
                );
              },
            ),
          ),
          const SizedBox(height: AppSizes.md),
          OutlinedButton.icon(
            icon: const Icon(Icons.lock, color: AppColors.error),
            label: const Text(
              AppStrings.tutupSesi,
              style: TextStyle(color: AppColors.error),
            ),
            style: OutlinedButton.styleFrom(
              side: const BorderSide(color: AppColors.error),
            ),
            onPressed: () {
              _showTutupSesiDialog(context, sesi.id);
            },
          ),
        ],
      ),
    );
  }

  void _showTutupSesiDialog(BuildContext context, String sesiId) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
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
                horizontal: AppSizes.md,
                vertical: AppSizes.sm,
              ),
            ),
            onPressed: () {
              Navigator.pop(ctx);
              context.read<SesiBloc>().add(
                    TutupSesiRequested(
                      sesiId: sesiId,
                      pecahanAkhir: const [],
                    ),
                  );
            },
            child: const Text(AppStrings.tutupSesi),
          ),
        ],
      ),
    );
  }
}

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
        .map((p) => {
              'pecahan_id': p.id,
              'jumlah': _jumlahMap[p.id],
            })
        .toList();

    context.read<SesiBloc>().add(BukaSesiRequested(pecahan: pecahanList));
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        // Total
        Container(
          padding: const EdgeInsets.all(AppSizes.md),
          color: AppColors.primary,
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                'Total Modal Awal',
                style: TextStyle(color: Colors.white70),
              ),
              Text(
                formatRupiah(_totalModal),
                style: const TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                  fontSize: 18,
                ),
              ),
            ],
          ),
        ),

        Expanded(
          child: ListView.separated(
            padding: const EdgeInsets.all(AppSizes.md),
            itemCount: widget.state.pecahan.length,
            separatorBuilder: (_, __) => const Divider(),
            itemBuilder: (context, index) {
              final pecahan = widget.state.pecahan[index];
              final jumlah = _jumlahMap[pecahan.id] ?? 0;

              return Row(
                children: [
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          pecahan.label,
                          style: const TextStyle(fontWeight: FontWeight.w500),
                        ),
                        Text(
                          'Subtotal: ${formatRupiah(pecahan.nominal * jumlah)}',
                          style: const TextStyle(
                            color: AppColors.textSecondary,
                            fontSize: 12,
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(
                    width: 80,
                    child: TextFormField(
                      initialValue: jumlah.toString(),
                      keyboardType: TextInputType.number,
                      inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                      textAlign: TextAlign.center,
                      decoration: const InputDecoration(
                        contentPadding: EdgeInsets.symmetric(
                          vertical: AppSizes.sm,
                          horizontal: AppSizes.xs,
                        ),
                      ),
                      onChanged: (val) {
                        setState(() {
                          _jumlahMap[pecahan.id] = int.tryParse(val) ?? 0;
                        });
                      },
                    ),
                  ),
                ],
              );
            },
          ),
        ),

        Padding(
          padding: const EdgeInsets.all(AppSizes.pagePadding),
          child: ElevatedButton.icon(
            icon: const Icon(Icons.lock_open),
            label: const Text(AppStrings.bukaSesi),
            onPressed: _totalModal > 0 ? _buka : null,
          ),
        ),
      ],
    );
  }
}
