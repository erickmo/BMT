import 'package:intl/intl.dart';

String formatRupiah(int amount) {
  return NumberFormat.currency(
    locale: 'id_ID',
    symbol: 'Rp ',
    decimalDigits: 0,
  ).format(amount);
}

String formatRupiahCompact(int amount) {
  if (amount >= 1000000000) {
    return 'Rp ${(amount / 1000000000).toStringAsFixed(1)}M';
  } else if (amount >= 1000000) {
    return 'Rp ${(amount / 1000000).toStringAsFixed(1)}jt';
  } else if (amount >= 1000) {
    return 'Rp ${(amount / 1000).toStringAsFixed(0)}rb';
  }
  return formatRupiah(amount);
}

String formatTanggal(DateTime dt) {
  return DateFormat('dd MMM yyyy', 'id_ID').format(dt);
}

String formatTanggalWaktu(DateTime dt) {
  return DateFormat('dd MMM yyyy HH:mm', 'id_ID').format(dt);
}

String formatTanggalPendek(DateTime dt) {
  return DateFormat('dd/MM/yyyy').format(dt);
}

String formatWaktu(DateTime dt) {
  return DateFormat('HH:mm').format(dt);
}

String formatPeriode(String periode) {
  // "2025-01" → "Januari 2025"
  try {
    final parts = periode.split('-');
    final dt = DateTime(int.parse(parts[0]), int.parse(parts[1]));
    return DateFormat('MMMM yyyy', 'id_ID').format(dt);
  } catch (_) {
    return periode;
  }
}

String formatNomorRekening(String nomor) {
  // Ensure consistent display: KODE-CAB-JENIS-SEQ
  return nomor;
}

String formatPersen(double value) {
  return '${value.toStringAsFixed(2)}%';
}

String formatNominalShort(int amount) {
  if (amount >= 1000000000) {
    return '${(amount / 1000000000).toStringAsFixed(2)}M';
  } else if (amount >= 1000000) {
    return '${(amount / 1000000).toStringAsFixed(2)}jt';
  }
  return formatRupiah(amount);
}
