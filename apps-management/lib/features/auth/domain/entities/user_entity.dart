import 'package:equatable/equatable.dart';

class UserEntity extends Equatable {
  final String id;
  final String nama;
  final String username;
  final String role;
  final String bmtId;
  final String? cabangId;
  final String accessToken;
  final String refreshToken;

  const UserEntity({
    required this.id,
    required this.nama,
    required this.username,
    required this.role,
    required this.bmtId,
    this.cabangId,
    required this.accessToken,
    required this.refreshToken,
  });

  bool get isManajerBMT => role == 'MANAJER_BMT';
  bool get isAdminBMT => role == 'ADMIN_BMT';
  bool get isManajerCabang => role == 'MANAJER_CABANG';
  bool get isKomite => role == 'KOMITE';
  bool get isAO => role == 'ACCOUNT_OFFICER';
  bool get isFinance => role == 'FINANCE';
  bool get isAuditor => role == 'AUDITOR_BMT';

  bool get canApproveForm =>
      isManajerBMT || isAdminBMT || isManajerCabang || isKomite;

  @override
  List<Object?> get props => [id, nama, username, role, bmtId, cabangId];
}
