import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import '../../../../core/constants/app_colors.dart';
import '../../../../core/constants/app_strings.dart';
import '../../../../core/utils/formatters.dart';

enum KasirStep { inputNominal, tapNfc, inputPin, konfirmasi, hasil }

class KasirPage extends StatefulWidget {
  const KasirPage({super.key});

  @override
  State<KasirPage> createState() => _KasirPageState();
}

class _KasirPageState extends State<KasirPage> {
  KasirStep _step = KasirStep.inputNominal;
  int _nominal = 0;
  String _nominalText = '';
  String _pin = '';
  String? _namaNasabah;
  bool _isSuccess = false;
  bool _isProcessing = false;

  void _onNominalDigit(String digit) {
    if (_nominalText.length >= 12) return;
    setState(() {
      _nominalText += digit;
      _nominal = int.tryParse(_nominalText) ?? 0;
    });
  }

  void _onNominalDelete() {
    if (_nominalText.isEmpty) return;
    setState(() {
      _nominalText = _nominalText.substring(0, _nominalText.length - 1);
      _nominal = int.tryParse(_nominalText) ?? 0;
    });
  }

  void _onNominalConfirm() {
    if (_nominal <= 0) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Masukkan nominal yang valid')),
      );
      return;
    }
    setState(() => _step = KasirStep.tapNfc);
    _simulateNfcScan();
  }

  void _simulateNfcScan() async {
    await Future.delayed(const Duration(seconds: 2));
    if (mounted) {
      setState(() {
        _namaNasabah = 'Muhammad Faqih';
        _step = KasirStep.inputPin;
      });
    }
  }

  void _onPinDigit(String digit) {
    if (_pin.length >= 6) return;
    setState(() {
      _pin += digit;
      if (_pin.length == 6) {
        _step = KasirStep.konfirmasi;
      }
    });
  }

  void _onPinDelete() {
    if (_pin.isEmpty) return;
    setState(() => _pin = _pin.substring(0, _pin.length - 1));
  }

  void _onConfirm() async {
    setState(() => _isProcessing = true);
    await Future.delayed(const Duration(seconds: 1));
    if (mounted) {
      setState(() {
        _isProcessing = false;
        _isSuccess = true;
        _step = KasirStep.hasil;
      });
    }
  }

  void _reset() {
    setState(() {
      _step = KasirStep.inputNominal;
      _nominal = 0;
      _nominalText = '';
      _pin = '';
      _namaNasabah = null;
      _isSuccess = false;
      _isProcessing = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.kasirBg,
      appBar: AppBar(
        title: const Text(AppStrings.kasir),
        backgroundColor: AppColors.kasirBg,
        foregroundColor: Colors.white,
        automaticallyImplyLeading: false,
        actions: [
          if (_step != KasirStep.inputNominal)
            IconButton(
              icon: const Icon(Icons.close),
              onPressed: _reset,
              tooltip: 'Batal',
            ),
        ],
      ),
      body: Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 420),
          child: _buildStep(),
        ),
      ),
    );
  }

  Widget _buildStep() {
    switch (_step) {
      case KasirStep.inputNominal:
        return _InputNominalView(
          nominalText: _nominalText,
          nominal: _nominal,
          onDigit: _onNominalDigit,
          onDelete: _onNominalDelete,
          onConfirm: _onNominalConfirm,
        );
      case KasirStep.tapNfc:
        return _TapNfcView();
      case KasirStep.inputPin:
        return _InputPinView(
          namaNasabah: _namaNasabah ?? '',
          nominal: _nominal,
          pin: _pin,
          onDigit: _onPinDigit,
          onDelete: _onPinDelete,
        );
      case KasirStep.konfirmasi:
        return _KonfirmasiView(
          namaNasabah: _namaNasabah ?? '',
          nominal: _nominal,
          isProcessing: _isProcessing,
          onConfirm: _onConfirm,
          onCancel: _reset,
        );
      case KasirStep.hasil:
        return _HasilView(
          isSuccess: _isSuccess,
          namaNasabah: _namaNasabah ?? '',
          nominal: _nominal,
          onReset: _reset,
        );
    }
  }
}

class _InputNominalView extends StatelessWidget {
  final String nominalText;
  final int nominal;
  final void Function(String) onDigit;
  final VoidCallback onDelete;
  final VoidCallback onConfirm;

  const _InputNominalView({
    required this.nominalText,
    required this.nominal,
    required this.onDigit,
    required this.onDelete,
    required this.onConfirm,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Text(
            AppStrings.inputNominal,
            style: TextStyle(color: Colors.white60, fontSize: 16),
          ),
          const SizedBox(height: 16),
          Text(
            nominalText.isEmpty ? 'Rp 0' : formatRupiah(nominal),
            style: const TextStyle(
              color: Colors.white,
              fontSize: 42,
              fontWeight: FontWeight.bold,
              fontFamily: 'monospace',
            ),
          ),
          const SizedBox(height: 32),
          _Numpad(onDigit: onDigit, onDelete: onDelete),
          const SizedBox(height: 24),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: nominal > 0 ? onConfirm : null,
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.primary,
                minimumSize: const Size.fromHeight(56),
              ),
              child: const Text('Lanjut', style: TextStyle(fontSize: 18)),
            ),
          ),
        ],
      ),
    );
  }
}

class _TapNfcView extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return const Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Icon(Icons.nfc, size: 120, color: AppColors.nfcActive),
        SizedBox(height: 24),
        Text(
          AppStrings.tapNfc,
          style: TextStyle(color: Colors.white, fontSize: 22, fontWeight: FontWeight.bold),
        ),
        SizedBox(height: 12),
        Text(
          'Menunggu kartu...',
          style: TextStyle(color: Colors.white60, fontSize: 14),
        ),
        SizedBox(height: 32),
        CircularProgressIndicator(color: AppColors.nfcActive),
      ],
    );
  }
}

class _InputPinView extends StatelessWidget {
  final String namaNasabah;
  final int nominal;
  final String pin;
  final void Function(String) onDigit;
  final VoidCallback onDelete;

  const _InputPinView({
    required this.namaNasabah,
    required this.nominal,
    required this.pin,
    required this.onDigit,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.person, size: 48, color: AppColors.nfcActive),
          const SizedBox(height: 8),
          Text(
            namaNasabah,
            style: const TextStyle(
                color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold),
          ),
          Text(
            formatRupiah(nominal),
            style: const TextStyle(
                color: AppColors.primaryLight, fontSize: 16),
          ),
          const SizedBox(height: 24),
          const Text(AppStrings.inputPin,
              style: TextStyle(color: Colors.white60, fontSize: 14)),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: List.generate(
              6,
              (i) => Container(
                margin: const EdgeInsets.symmetric(horizontal: 8),
                width: 16,
                height: 16,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: i < pin.length ? Colors.white : Colors.white24,
                ),
              ),
            ),
          ),
          const SizedBox(height: 32),
          _Numpad(onDigit: onDigit, onDelete: onDelete),
        ],
      ),
    );
  }
}

class _KonfirmasiView extends StatelessWidget {
  final String namaNasabah;
  final int nominal;
  final bool isProcessing;
  final VoidCallback onConfirm;
  final VoidCallback onCancel;

  const _KonfirmasiView({
    required this.namaNasabah,
    required this.nominal,
    required this.isProcessing,
    required this.onConfirm,
    required this.onCancel,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.check_circle_outline,
              size: 64, color: AppColors.nfcActive),
          const SizedBox(height: 24),
          const Text(AppStrings.konfirmasi,
              style: TextStyle(
                  color: Colors.white, fontSize: 22, fontWeight: FontWeight.bold)),
          const SizedBox(height: 24),
          Card(
            color: Colors.white10,
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                children: [
                  _ConfirmRow('Nasabah', namaNasabah),
                  const Divider(color: Colors.white24),
                  _ConfirmRow('Nominal', formatRupiah(nominal)),
                  const Divider(color: Colors.white24),
                  _ConfirmRow('Waktu', formatTanggalWaktu(DateTime.now())),
                ],
              ),
            ),
          ),
          const SizedBox(height: 32),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: isProcessing ? null : onConfirm,
              style: ElevatedButton.styleFrom(
                  backgroundColor: AppColors.success,
                  minimumSize: const Size.fromHeight(56)),
              child: isProcessing
                  ? const CircularProgressIndicator(color: Colors.white)
                  : const Text('Konfirmasi Bayar',
                      style: TextStyle(fontSize: 18)),
            ),
          ),
          const SizedBox(height: 12),
          SizedBox(
            width: double.infinity,
            child: OutlinedButton(
              onPressed: isProcessing ? null : onCancel,
              style: OutlinedButton.styleFrom(
                side: const BorderSide(color: Colors.white30),
                foregroundColor: Colors.white60,
                minimumSize: const Size.fromHeight(48),
              ),
              child: const Text(AppStrings.cancel),
            ),
          ),
        ],
      ),
    );
  }
}

class _ConfirmRow extends StatelessWidget {
  final String label;
  final String value;
  const _ConfirmRow(this.label, this.value);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: Colors.white60)),
          Text(value,
              style: const TextStyle(
                  color: Colors.white, fontWeight: FontWeight.w600)),
        ],
      ),
    );
  }
}

class _HasilView extends StatelessWidget {
  final bool isSuccess;
  final String namaNasabah;
  final int nominal;
  final VoidCallback onReset;

  const _HasilView({
    required this.isSuccess,
    required this.namaNasabah,
    required this.nominal,
    required this.onReset,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            isSuccess ? Icons.check_circle : Icons.error,
            size: 80,
            color: isSuccess ? AppColors.success : AppColors.error,
          ),
          const SizedBox(height: 16),
          Text(
            isSuccess ? AppStrings.berhasil : AppStrings.gagal,
            style: TextStyle(
              color: isSuccess ? AppColors.success : AppColors.error,
              fontSize: 26,
              fontWeight: FontWeight.bold,
            ),
          ),
          if (isSuccess) ...[
            const SizedBox(height: 8),
            Text(formatRupiah(nominal),
                style: const TextStyle(color: Colors.white, fontSize: 22)),
            const SizedBox(height: 4),
            Text(namaNasabah,
                style: const TextStyle(color: Colors.white60, fontSize: 14)),
          ],
          const SizedBox(height: 40),
          if (isSuccess)
            SizedBox(
              width: double.infinity,
              child: OutlinedButton.icon(
                onPressed: () {},
                icon: const Icon(Icons.receipt, color: Colors.white60),
                label: const Text(AppStrings.cetak,
                    style: TextStyle(color: Colors.white60)),
                style: OutlinedButton.styleFrom(
                  side: const BorderSide(color: Colors.white24),
                  minimumSize: const Size.fromHeight(48),
                ),
              ),
            ),
          const SizedBox(height: 12),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: onReset,
              style: ElevatedButton.styleFrom(
                  backgroundColor: AppColors.primary,
                  minimumSize: const Size.fromHeight(52)),
              child: const Text(AppStrings.transaksiLagi),
            ),
          ),
        ],
      ),
    );
  }
}

class _Numpad extends StatelessWidget {
  final void Function(String) onDigit;
  final VoidCallback onDelete;

  const _Numpad({required this.onDigit, required this.onDelete});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        for (final row in [
          ['1', '2', '3'],
          ['4', '5', '6'],
          ['7', '8', '9'],
          ['000', '0', '⌫'],
        ])
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: row.map((key) {
              return GestureDetector(
                onTap: () {
                  if (key == '⌫') {
                    onDelete();
                  } else {
                    onDigit(key);
                  }
                },
                child: Container(
                  width: 88,
                  height: 60,
                  margin: const EdgeInsets.all(6),
                  decoration: BoxDecoration(
                    color: key == '⌫'
                        ? Colors.red.withOpacity(0.2)
                        : Colors.white10,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  alignment: Alignment.center,
                  child: Text(
                    key,
                    style: TextStyle(
                      color: key == '⌫' ? Colors.red.shade300 : Colors.white,
                      fontSize: 20,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              );
            }).toList(),
          ),
      ],
    );
  }
}
