import '../../domain/entities/staf_entity.dart';

class AuthResponseModel {
  final String accessToken;
  final String refreshToken;
  final StafModel staf;

  const AuthResponseModel({
    required this.accessToken,
    required this.refreshToken,
    required this.staf,
  });

  factory AuthResponseModel.fromJson(Map<String, dynamic> json) {
    return AuthResponseModel(
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
      staf: StafModel.fromJson(json['staf'] as Map<String, dynamic>),
    );
  }

  StafEntity toEntity() => StafEntity(
        id: staf.id,
        nama: staf.nama,
        email: staf.email,
        role: staf.role,
        cabangId: staf.cabangId,
        cabangNama: staf.cabangNama,
        bmtId: staf.bmtId,
        accessToken: accessToken,
        refreshToken: refreshToken,
      );
}

class StafModel {
  final String id;
  final String nama;
  final String email;
  final String role;
  final String cabangId;
  final String cabangNama;
  final String bmtId;

  const StafModel({
    required this.id,
    required this.nama,
    required this.email,
    required this.role,
    required this.cabangId,
    required this.cabangNama,
    required this.bmtId,
  });

  factory StafModel.fromJson(Map<String, dynamic> json) {
    return StafModel(
      id: json['id'] as String,
      nama: json['nama'] as String,
      email: json['email'] as String? ?? '',
      role: json['role'] as String,
      cabangId: json['cabang_id'] as String,
      cabangNama: json['cabang_nama'] as String? ?? '',
      bmtId: json['bmt_id'] as String,
    );
  }
}
