import 'package:intl/intl.dart';

String formatRupiah(int amount) {
  final formatter = NumberFormat.currency(
    locale: 'id_ID',
    symbol: 'Rp ',
    decimalDigits: 0,
  );
  return formatter.format(amount);
}

String formatRupiahCompact(int amount) {
  if (amount >= 1000000000) {
    return 'Rp ${(amount / 1000000000).toStringAsFixed(1)} M';
  } else if (amount >= 1000000) {
    return 'Rp ${(amount / 1000000).toStringAsFixed(1)} Jt';
  } else if (amount >= 1000) {
    return 'Rp ${(amount / 1000).toStringAsFixed(0)} Rb';
  }
  return formatRupiah(amount);
}

String formatTanggal(DateTime date) =>
    DateFormat('dd MMMM yyyy', 'id_ID').format(date);

String formatTanggalWaktu(DateTime date) =>
    DateFormat('dd MMM yyyy, HH:mm', 'id_ID').format(date);

String formatWaktu(DateTime date) => DateFormat('HH:mm').format(date);

DateTime? parseTanggal(String? raw) {
  if (raw == null || raw.isEmpty) return null;
  try {
    return DateTime.parse(raw).toLocal();
  } catch (_) {
    return null;
  }
}

/// Parse nominal string (may contain Rp, commas, dots) to int
int parseNominal(String raw) {
  final cleaned = raw.replaceAll(RegExp(r'[^0-9]'), '');
  return int.tryParse(cleaned) ?? 0;
}
