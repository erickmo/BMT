import 'package:equatable/equatable.dart';

class UserEntity extends Equatable {
  final String id;
  final String nomorNasabah;
  final String nama;
  final String? email;
  final String? telepon;
  final String? fotoUrl;
  final String accessToken;
  final String refreshToken;

  const UserEntity({
    required this.id,
    required this.nomorNasabah,
    required this.nama,
    this.email,
    this.telepon,
    this.fotoUrl,
    required this.accessToken,
    required this.refreshToken,
  });

  @override
  List<Object?> get props => [
        id,
        nomorNasabah,
        nama,
        email,
        telepon,
        fotoUrl,
        accessToken,
        refreshToken,
      ];
}
