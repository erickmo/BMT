import 'package:intl/intl.dart';

/// Format integer (int64 from API) as Indonesian Rupiah
String formatRupiah(int amount) {
  final formatter = NumberFormat.currency(
    locale: 'id_ID',
    symbol: 'Rp ',
    decimalDigits: 0,
  );
  return formatter.format(amount);
}

/// Format integer as compact Rupiah (e.g., Rp 1,5 Jt)
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

/// Format DateTime to Indonesian date string
String formatTanggal(DateTime date) {
  return DateFormat('dd MMMM yyyy', 'id_ID').format(date);
}

/// Format DateTime to Indonesian date-time string
String formatTanggalWaktu(DateTime date) {
  return DateFormat('dd MMM yyyy, HH:mm', 'id_ID').format(date);
}

/// Format DateTime to short date
String formatTanggalPendek(DateTime date) {
  return DateFormat('dd/MM/yyyy').format(date);
}

/// Format DateTime to time only
String formatWaktu(DateTime date) {
  return DateFormat('HH:mm').format(date);
}

/// Parse ISO8601 string safely
DateTime? parseTanggal(String? raw) {
  if (raw == null || raw.isEmpty) return null;
  try {
    return DateTime.parse(raw).toLocal();
  } catch (_) {
    return null;
  }
}

/// Mask nomor rekening: show only last 4 digits
String maskNomorRekening(String nomor) {
  if (nomor.length <= 4) return nomor;
  final visible = nomor.substring(nomor.length - 4);
  final masked = '*' * (nomor.length - 4);
  return '$masked$visible';
}

/// Format nomor rekening with dashes for display
String formatNomorRekening(String nomor) {
  return nomor; // e.g. ANNUR-KDR-SU-00000001 already formatted
}
