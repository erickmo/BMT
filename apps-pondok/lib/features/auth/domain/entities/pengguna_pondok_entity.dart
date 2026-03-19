import 'package:equatable/equatable.dart';

class PenggunaPondokEntity extends Equatable {
  final String id;
  final String nama;
  final String username;
  final String role;
  final String bmtId;
  final String cabangId;
  final String accessToken;
  final String refreshToken;

  const PenggunaPondokEntity({
    required this.id,
    required this.nama,
    required this.username,
    required this.role,
    required this.bmtId,
    required this.cabangId,
    required this.accessToken,
    required this.refreshToken,
  });

  bool get isAdmin => role == 'ADMIN_PONDOK';
  bool get isOperator => role == 'OPERATOR_PONDOK';
  bool get isBendahara => role == 'BENDAHARA_PONDOK';

  @override
  List<Object?> get props => [id, username, role];
}
