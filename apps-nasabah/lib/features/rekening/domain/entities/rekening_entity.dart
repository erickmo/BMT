import 'package:equatable/equatable.dart';

class RekeningEntity extends Equatable {
  final String id;
  final String nomorRekening;
  final String jenisRekeningNama;
  final String jenisRekeningKode;
  final int saldo;
  final String status;
  final String? alasanBlokir;
  final int biayaAdminBulanan;
  final DateTime tanggalBuka;
  final DateTime? tanggalTutup;

  const RekeningEntity({
    required this.id,
    required this.nomorRekening,
    required this.jenisRekeningNama,
    required this.jenisRekeningKode,
    required this.saldo,
    required this.status,
    this.alasanBlokir,
    required this.biayaAdminBulanan,
    required this.tanggalBuka,
    this.tanggalTutup,
  });

  bool get isAktif => status == 'AKTIF';
  bool get isBlokir => status == 'BLOKIR';

  @override
  List<Object?> get props => [
        id,
        nomorRekening,
        jenisRekeningNama,
        jenisRekeningKode,
        saldo,
        status,
        alasanBlokir,
        biayaAdminBulanan,
        tanggalBuka,
        tanggalTutup,
      ];
}
