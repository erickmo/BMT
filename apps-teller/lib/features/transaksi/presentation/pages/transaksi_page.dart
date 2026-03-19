import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_sizes.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';
import '../../../../core/utils/validators.dart';
import '../../domain/entities/transaksi_entity.dart';
import '../bloc/transaksi_bloc.dart';

class TransaksiPage extends StatefulWidget {
  const TransaksiPage({super.key});

  @override
  State<TransaksiPage> createState() => _TransaksiPageState();
}

class _TransaksiPageState extends State<TransaksiPage>
    with SingleTickerProviderStateMixin {
  late final TabController _tabController;
  final _searchCtrl = TextEditingController();

  NasabahSearchResult? _selectedNasabah;
  RekeningSearchResult? _selectedRekening;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _searchCtrl.dispose();
    super.dispose();
  }

  void _resetPilihan() {
    setState(() {
      _selectedNasabah = null;
      _selectedRekening = null;
    });
    _searchCtrl.clear();
    context.read<TransaksiBloc>().add(const ResetTransaksi());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Transaksi Tunai'),
        bottom: TabBar(
          controller: _tabController,
          labelColor: Colors.white,
          unselectedLabelColor: Colors.white60,
          indicatorColor: Colors.white,
          tabs: const [
            Tab(icon: Icon(Icons.arrow_downward), text: AppStrings.setor),
            Tab(icon: Icon(Icons.arrow_upward), text: AppStrings.tarik),
          ],
        ),
      ),
      body: BlocConsumer<TransaksiBloc, TransaksiState>(
        listener: (context, state) {
          if (state is TransaksiSuccess) {
            _showSuksesDialog(context, state.result);
          }
          if (state is TransaksiError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: AppColors.error,
                behavior: SnackBarBehavior.floating,
              ),
            );
          }
        },
        builder: (context, state) {
          return TabBarView(
            controller: _tabController,
            children: [
              _TransaksiForm(
                jenis: 'SETOR',
                selectedNasabah: _selectedNasabah,
                selectedRekening: _selectedRekening,
                searchCtrl: _searchCtrl,
                onNasabahSelected: (n) =>
                    setState(() => _selectedNasabah = n),
                onRekeningSelected: (r) =>
                    setState(() => _selectedRekening = r),
                onReset: _resetPilihan,
                isLoading: state is TransaksiLoading,
                searchResults: state is NasabahSearchLoaded
                    ? state.results
                    : null,
              ),
              _TransaksiForm(
                jenis: 'TARIK',
                selectedNasabah: _selectedNasabah,
                selectedRekening: _selectedRekening,
                searchCtrl: _searchCtrl,
                onNasabahSelected: (n) =>
                    setState(() => _selectedNasabah = n),
                onRekeningSelected: (r) =>
                    setState(() => _selectedRekening = r),
                onReset: _resetPilihan,
                isLoading: state is TransaksiLoading,
                searchResults: state is NasabahSearchLoaded
                    ? state.results
                    : null,
              ),
            ],
          );
        },
      ),
    );
  }

  void _showSuksesDialog(
    BuildContext context,
    TransaksiResultEntity result,
  ) {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (ctx) => AlertDialog(
        title: Row(
          children: [
            const Icon(Icons.check_circle, color: AppColors.success),
            const SizedBox(width: AppSizes.sm),
            const Text(AppStrings.transaksiSukses),
          ],
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _SlipRow(label: 'Nasabah', value: result.namaNasabah),
            _SlipRow(label: 'Rekening', value: result.nomorRekening),
            _SlipRow(
              label: 'Jenis',
              value: result.jenis == 'KREDIT' ? 'Setor Tunai' : 'Tarik Tunai',
            ),
            _SlipRow(label: 'Nominal', value: formatRupiah(result.nominal)),
            _SlipRow(
                label: 'Saldo Sebelum',
                value: formatRupiah(result.saldoSebelum)),
            _SlipRow(
                label: 'Saldo Akhir',
                value: formatRupiah(result.saldoAkhir)),
            _SlipRow(
                label: 'Waktu',
                value: formatTanggalWaktu(result.tanggal)),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.pop(ctx);
              _resetPilihan();
            },
            child: const Text(AppStrings.close),
          ),
        ],
      ),
    );
  }
}

class _TransaksiForm extends StatefulWidget {
  final String jenis; // SETOR | TARIK
  final NasabahSearchResult? selectedNasabah;
  final RekeningSearchResult? selectedRekening;
  final TextEditingController searchCtrl;
  final void Function(NasabahSearchResult) onNasabahSelected;
  final void Function(RekeningSearchResult) onRekeningSelected;
  final VoidCallback onReset;
  final bool isLoading;
  final List<NasabahSearchResult>? searchResults;

  const _TransaksiForm({
    required this.jenis,
    required this.selectedNasabah,
    required this.selectedRekening,
    required this.searchCtrl,
    required this.onNasabahSelected,
    required this.onRekeningSelected,
    required this.onReset,
    required this.isLoading,
    this.searchResults,
  });

  @override
  State<_TransaksiForm> createState() => _TransaksiFormState();
}

class _TransaksiFormState extends State<_TransaksiForm> {
  final _formKey = GlobalKey<FormState>();
  final _nominalCtrl = TextEditingController();
  final _keteranganCtrl = TextEditingController();

  @override
  void dispose() {
    _nominalCtrl.dispose();
    _keteranganCtrl.dispose();
    super.dispose();
  }

  void _proses() {
    if (!(_formKey.currentState?.validate() ?? false)) return;
    if (widget.selectedRekening == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Pilih rekening terlebih dahulu')),
      );
      return;
    }

    final nominal = parseNominal(_nominalCtrl.text);
    final rek = widget.selectedRekening!;

    if (widget.jenis == 'SETOR') {
      context.read<TransaksiBloc>().add(
            SetorRequested(
              rekeningId: rek.id,
              nominal: nominal,
              keterangan: _keteranganCtrl.text.trim(),
            ),
          );
    } else {
      context.read<TransaksiBloc>().add(
            TarikRequested(
              rekeningId: rek.id,
              nominal: nominal,
              keterangan: _keteranganCtrl.text.trim(),
            ),
          );
    }
  }

  @override
  Widget build(BuildContext context) {
    final color =
        widget.jenis == 'SETOR' ? AppColors.setor : AppColors.tarik;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSizes.pagePadding),
      child: Form(
        key: _formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Cari nasabah
            if (widget.selectedNasabah == null) ...[
              TextFormField(
                controller: widget.searchCtrl,
                decoration: const InputDecoration(
                  labelText: AppStrings.cariNasabah,
                  hintText: AppStrings.nomorNasabahOrNama,
                  prefixIcon: Icon(Icons.search),
                ),
                onChanged: (v) {
                  if (v.length >= 3) {
                    context.read<TransaksiBloc>().add(CariNasabah(v));
                  }
                },
              ),
              if (widget.searchResults != null) ...[
                const SizedBox(height: AppSizes.sm),
                Card(
                  child: Column(
                    children: widget.searchResults!.map((n) {
                      return ListTile(
                        title: Text(n.nama),
                        subtitle: Text(n.nomorNasabah),
                        onTap: () => widget.onNasabahSelected(n),
                      );
                    }).toList(),
                  ),
                ),
              ],
            ] else ...[
              // Nasabah terpilih
              Card(
                child: ListTile(
                  leading: const CircleAvatar(
                    backgroundColor: AppColors.primary,
                    child: Icon(Icons.person, color: Colors.white),
                  ),
                  title: Text(widget.selectedNasabah!.nama),
                  subtitle: Text(widget.selectedNasabah!.nomorNasabah),
                  trailing: IconButton(
                    icon: const Icon(Icons.close),
                    onPressed: widget.onReset,
                  ),
                ),
              ),
              const SizedBox(height: AppSizes.md),

              // Pilih rekening
              const Text(
                AppStrings.pilihRekening,
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: AppSizes.sm),
              ...widget.selectedNasabah!.rekening.map(
                (rek) {
                  final isSelected =
                      widget.selectedRekening?.id == rek.id;
                  return Card(
                    child: InkWell(
                      borderRadius:
                          BorderRadius.circular(AppSizes.radiusMd),
                      onTap:
                          rek.isAktif ? () => widget.onRekeningSelected(rek) : null,
                      child: Padding(
                        padding: const EdgeInsets.all(AppSizes.sm),
                        child: Row(
                          children: [
                            Icon(
                              isSelected
                                  ? Icons.radio_button_checked
                                  : Icons.radio_button_unchecked,
                              color: isSelected
                                  ? AppColors.primary
                                  : AppColors.textHint,
                            ),
                            const SizedBox(width: AppSizes.sm),
                            Expanded(
                              child: Column(
                                crossAxisAlignment:
                                    CrossAxisAlignment.start,
                                children: [
                                  Text(rek.jenisNama,
                                      style: const TextStyle(
                                          fontWeight: FontWeight.w500)),
                                  Text(rek.nomorRekening,
                                      style: const TextStyle(
                                          fontSize: 13,
                                          color: AppColors.textSecondary)),
                                  Text(
                                    'Saldo: ${formatRupiah(rek.saldo)}',
                                    style: const TextStyle(
                                      color: AppColors.primary,
                                      fontWeight: FontWeight.w500,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            if (!rek.isAktif)
                              const Chip(label: Text('Blokir')),
                          ],
                        ),
                      ),
                    ),
                  );
                },
              ),

              if (widget.selectedRekening != null) ...[
                const SizedBox(height: AppSizes.md),
                TextFormField(
                  controller: _nominalCtrl,
                  keyboardType: TextInputType.number,
                  inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                  decoration: InputDecoration(
                    labelText: AppStrings.nominal,
                    prefixText: 'Rp ',
                    hintText: '0',
                    focusedBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(AppSizes.radiusMd),
                      borderSide: BorderSide(color: color, width: 2),
                    ),
                  ),
                  validator: Validators.nominal,
                ),
                const SizedBox(height: AppSizes.md),
                TextFormField(
                  controller: _keteranganCtrl,
                  decoration: const InputDecoration(
                    labelText: AppStrings.keterangan,
                    hintText: 'Opsional',
                  ),
                  maxLines: 2,
                ),
                const SizedBox(height: AppSizes.xl),
                ElevatedButton(
                  style: ElevatedButton.styleFrom(
                    backgroundColor: color,
                  ),
                  onPressed: widget.isLoading ? null : _proses,
                  child: widget.isLoading
                      ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            valueColor:
                                AlwaysStoppedAnimation(Colors.white),
                          ),
                        )
                      : Text(
                          widget.jenis == 'SETOR'
                              ? AppStrings.setor
                              : AppStrings.tarik,
                        ),
                ),
              ],
            ],
          ],
        ),
      ),
    );
  }
}

class _SlipRow extends StatelessWidget {
  final String label;
  final String value;

  const _SlipRow({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 3),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label,
              style: const TextStyle(color: AppColors.textSecondary)),
          Text(value, style: const TextStyle(fontWeight: FontWeight.w500)),
        ],
      ),
    );
  }
}
