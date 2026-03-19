import 'package:equatable/equatable.dart';

class ProfilEntity extends Equatable {
  final String id;
  final String nomorNasabah;
  final String nama;
  final String? email;
  final String? telepon;
  final String? alamat;
  final String? fotoUrl;
  final String? nik;
  final DateTime? tanggalLahir;

  const ProfilEntity({
    required this.id,
    required this.nomorNasabah,
    required this.nama,
    this.email,
    this.telepon,
    this.alamat,
    this.fotoUrl,
    this.nik,
    this.tanggalLahir,
  });

  @override
  List<Object?> get props => [
        id,
        nomorNasabah,
        nama,
        email,
        telepon,
        alamat,
        fotoUrl,
        nik,
        tanggalLahir,
      ];
}
