import 'package:intl/intl.dart';

String formatRupiah(int amount) {
  return NumberFormat.currency(locale: 'id_ID', symbol: 'Rp ', decimalDigits: 0).format(amount);
}

String formatTanggal(DateTime dt) {
  return DateFormat('dd MMM yyyy', 'id_ID').format(dt);
}

String formatTanggalWaktu(DateTime dt) {
  return DateFormat('dd MMM yyyy HH:mm', 'id_ID').format(dt);
}

String formatWaktu(DateTime dt) {
  return DateFormat('HH:mm').format(dt);
}

String formatPeriode(String periode) {
  try {
    final parts = periode.split('-');
    final dt = DateTime(int.parse(parts[0]), int.parse(parts[1]));
    return DateFormat('MMMM yyyy', 'id_ID').format(dt);
  } catch (_) {
    return periode;
  }
}

String formatNilai(double nilai) {
  return nilai.toStringAsFixed(1);
}
