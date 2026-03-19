part of 'kiosk_bloc.dart';

abstract class KioskEvent extends Equatable {
  const KioskEvent();
  @override
  List<Object?> get props => [];
}

class NfcTagDetected extends KioskEvent {
  final String uid;
  const NfcTagDetected(this.uid);
  @override
  List<Object?> get props => [uid];
}

class KioskReset extends KioskEvent {
  const KioskReset();
}
