import '../../domain/entities/user_entity.dart';

class AuthResponseModel {
  final String accessToken;
  final String refreshToken;
  final NasabahModel nasabah;

  const AuthResponseModel({
    required this.accessToken,
    required this.refreshToken,
    required this.nasabah,
  });

  factory AuthResponseModel.fromJson(Map<String, dynamic> json) {
    return AuthResponseModel(
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
      nasabah: NasabahModel.fromJson(json['nasabah'] as Map<String, dynamic>),
    );
  }

  UserEntity toEntity() {
    return UserEntity(
      id: nasabah.id,
      nomorNasabah: nasabah.nomorNasabah,
      nama: nasabah.nama,
      email: nasabah.email,
      telepon: nasabah.telepon,
      fotoUrl: nasabah.fotoUrl,
      accessToken: accessToken,
      refreshToken: refreshToken,
    );
  }
}

class NasabahModel {
  final String id;
  final String nomorNasabah;
  final String nama;
  final String? email;
  final String? telepon;
  final String? fotoUrl;

  const NasabahModel({
    required this.id,
    required this.nomorNasabah,
    required this.nama,
    this.email,
    this.telepon,
    this.fotoUrl,
  });

  factory NasabahModel.fromJson(Map<String, dynamic> json) {
    return NasabahModel(
      id: json['id'] as String,
      nomorNasabah: json['nomor_nasabah'] as String,
      nama: json['nama'] as String,
      email: json['email'] as String?,
      telepon: json['telepon'] as String?,
      fotoUrl: json['foto_url'] as String?,
    );
  }
}
