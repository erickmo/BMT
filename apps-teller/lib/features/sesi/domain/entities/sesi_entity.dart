import 'package:equatable/equatable.dart';

class SesiEntity extends Equatable {
  final String id;
  final String tellerId;
  final String cabangId;
  final int modalAwal;
  final int modalAkhir;
  final String status; // AKTIF | TUTUP
  final DateTime dibukaPada;
  final DateTime? ditutupPada;
  final List<PecahanSesiEntity> pecahan;

  const SesiEntity({
    required this.id,
    required this.tellerId,
    required this.cabangId,
    required this.modalAwal,
    required this.modalAkhir,
    required this.status,
    required this.dibukaPada,
    this.ditutupPada,
    required this.pecahan,
  });

  bool get isAktif => status == 'AKTIF';

  @override
  List<Object?> get props => [
        id,
        tellerId,
        cabangId,
        modalAwal,
        modalAkhir,
        status,
        dibukaPada,
        ditutupPada,
        pecahan,
      ];
}

class PecahanSesiEntity extends Equatable {
  final String pecahanId;
  final int nominal;
  final String label;
  final String jenis; // LOGAM | KERTAS
  final int jumlah;

  const PecahanSesiEntity({
    required this.pecahanId,
    required this.nominal,
    required this.label,
    required this.jenis,
    required this.jumlah,
  });

  int get subtotal => nominal * jumlah;

  @override
  List<Object> get props => [pecahanId, nominal, label, jenis, jumlah];
}
