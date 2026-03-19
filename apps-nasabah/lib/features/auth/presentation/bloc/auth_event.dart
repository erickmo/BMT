part of 'auth_bloc.dart';

sealed class AuthEvent extends Equatable {
  const AuthEvent();

  @override
  List<Object> get props => [];
}

class AuthCheckRequested extends AuthEvent {
  const AuthCheckRequested();
}

class LoginRequested extends AuthEvent {
  final String nomorNasabah;
  final String pin;
  final String deviceId;

  const LoginRequested({
    required this.nomorNasabah,
    required this.pin,
    required this.deviceId,
  });

  @override
  List<Object> get props => [nomorNasabah, pin, deviceId];
}

class LogoutRequested extends AuthEvent {
  const LogoutRequested();
}
