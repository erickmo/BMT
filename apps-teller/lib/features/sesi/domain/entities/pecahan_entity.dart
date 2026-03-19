import 'package:equatable/equatable.dart';

/// Pecahan uang dari DB — bukan konstanta kode
class PecahanEntity extends Equatable {
  final String id;
  final int nominal;
  final String jenis; // LOGAM | KERTAS
  final String label;
  final bool isAktif;
  final int urutan;

  const PecahanEntity({
    required this.id,
    required this.nominal,
    required this.jenis,
    required this.label,
    required this.isAktif,
    required this.urutan,
  });

  @override
  List<Object> get props => [id, nominal, jenis, label, isAktif, urutan];
}
