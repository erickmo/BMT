part of 'kiosk_bloc.dart';

abstract class KioskState extends Equatable {
  const KioskState();
  @override
  List<Object?> get props => [];
}

class KioskIdle extends KioskState {}

class KioskLoading extends KioskState {}

class KioskShowSaldo extends KioskState {
  final SaldoInfoEntity saldoInfo;
  const KioskShowSaldo(this.saldoInfo);
  @override
  List<Object?> get props => [saldoInfo];
}

class KioskError extends KioskState {
  final String message;
  const KioskError(this.message);
  @override
  List<Object?> get props => [message];
}
