import 'package:equatable/equatable.dart';

class TransaksiEntity extends Equatable {
  final String id;
  final String rekeningId;
  final String jenis; // DEBIT | KREDIT
  final String tipe;  // e.g. SETOR_TUNAI, TARIK_TUNAI, etc.
  final int nominal;
  final int saldoAkhir;
  final String keterangan;
  final DateTime tanggal;
  final String? referensi;

  const TransaksiEntity({
    required this.id,
    required this.rekeningId,
    required this.jenis,
    required this.tipe,
    required this.nominal,
    required this.saldoAkhir,
    required this.keterangan,
    required this.tanggal,
    this.referensi,
  });

  bool get isDebit => jenis == 'DEBIT';
  bool get isKredit => jenis == 'KREDIT';

  @override
  List<Object?> get props => [
        id,
        rekeningId,
        jenis,
        tipe,
        nominal,
        saldoAkhir,
        keterangan,
        tanggal,
        referensi,
      ];
}
