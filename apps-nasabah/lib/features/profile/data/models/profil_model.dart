import '../../domain/entities/profil_entity.dart';

class ProfilModel {
  final String id;
  final String nomorNasabah;
  final String nama;
  final String? email;
  final String? telepon;
  final String? alamat;
  final String? fotoUrl;
  final String? nik;
  final String? tanggalLahir;

  const ProfilModel({
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

  factory ProfilModel.fromJson(Map<String, dynamic> json) {
    return ProfilModel(
      id: json['id'] as String,
      nomorNasabah: json['nomor_nasabah'] as String,
      nama: json['nama'] as String,
      email: json['email'] as String?,
      telepon: json['telepon'] as String?,
      alamat: json['alamat'] as String?,
      fotoUrl: json['foto_url'] as String?,
      nik: json['nik'] as String?,
      tanggalLahir: json['tanggal_lahir'] as String?,
    );
  }

  ProfilEntity toEntity() {
    return ProfilEntity(
      id: id,
      nomorNasabah: nomorNasabah,
      nama: nama,
      email: email,
      telepon: telepon,
      alamat: alamat,
      fotoUrl: fotoUrl,
      nik: nik,
      tanggalLahir:
          tanggalLahir != null ? DateTime.tryParse(tanggalLahir!) : null,
    );
  }
}
