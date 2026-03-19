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
