import 'package:equatable/equatable.dart';

class NasabahEntity extends Equatable {
  final String id;
  final String bmtId;
  final String cabangId;
  final String nomorNasabah;
  final String namaLengkap;
  final String nik;
  final String telepon;
  final String? email;
  final String alamat;
  final String status;
  final DateTime createdAt;

  const NasabahEntity({
    required this.id,
    required this.bmtId,
    required this.cabangId,
    required this.nomorNasabah,
    required this.namaLengkap,
    required this.nik,
    required this.telepon,
    this.email,
    required this.alamat,
    required this.status,
    required this.createdAt,
  });

  bool get isAktif => status == 'AKTIF';

  @override
  List<Object?> get props => [id, nomorNasabah];
}
