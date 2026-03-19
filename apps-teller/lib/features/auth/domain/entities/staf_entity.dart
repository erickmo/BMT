import 'package:equatable/equatable.dart';

class StafEntity extends Equatable {
  final String id;
  final String nama;
  final String email;
  final String role;
  final String cabangId;
  final String cabangNama;
  final String bmtId;
  final String accessToken;
  final String refreshToken;

  const StafEntity({
    required this.id,
    required this.nama,
    required this.email,
    required this.role,
    required this.cabangId,
    required this.cabangNama,
    required this.bmtId,
    required this.accessToken,
    required this.refreshToken,
  });

  bool get isTeller => role == 'TELLER';

  @override
  List<Object> get props => [
        id,
        nama,
        email,
        role,
        cabangId,
        cabangNama,
        bmtId,
        accessToken,
        refreshToken,
      ];
}
