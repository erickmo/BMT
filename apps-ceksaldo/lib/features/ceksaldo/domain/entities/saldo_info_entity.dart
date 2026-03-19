import 'package:equatable/equatable.dart';

class TransaksiSingkat extends Equatable {
  final String keterangan;
  final int nominal;
  final bool isKredit;
  final DateTime tanggal;

  const TransaksiSingkat({
    required this.keterangan,
    required this.nominal,
    required this.isKredit,
    required this.tanggal,
  });

  @override
  List<Object?> get props => [keterangan, nominal, isKredit, tanggal];
}

class SaldoInfoEntity extends Equatable {
  final String namaNasabah;
  final String nomorRekening;
  final int saldo;
  final List<TransaksiSingkat> transaksiTerakhir;

  const SaldoInfoEntity({
    required this.namaNasabah,
    required this.nomorRekening,
    required this.saldo,
    required this.transaksiTerakhir,
  });

  @override
  List<Object?> get props => [namaNasabah, nomorRekening, saldo];
}
