part of 'sesi_bloc.dart';

sealed class SesiState extends Equatable {
  const SesiState();

  @override
  List<Object?> get props => [];
}

class SesiInitial extends SesiState {
  const SesiInitial();
}

class SesiLoading extends SesiState {
  const SesiLoading();
}

class SesiAktifLoaded extends SesiState {
  final SesiEntity? sesi; // null = tidak ada sesi aktif

  const SesiAktifLoaded({this.sesi});

  @override
  List<Object?> get props => [sesi];
}

class PecahanLoaded extends SesiState {
  final List<PecahanEntity> pecahan;
  final Map<String, int> jumlahMap; // pecahan_id → jumlah

  const PecahanLoaded({required this.pecahan, required this.jumlahMap});

  int get totalModal => pecahan.fold(
        0,
        (sum, p) => sum + p.nominal * (jumlahMap[p.id] ?? 0),
      );

  @override
  List<Object?> get props => [pecahan, jumlahMap];
}

class SesiBukaSuccess extends SesiState {
  final SesiEntity sesi;

  const SesiBukaSuccess({required this.sesi});

  @override
  List<Object?> get props => [sesi];
}

class SesiTutupSuccess extends SesiState {
  final SesiEntity sesi;

  const SesiTutupSuccess({required this.sesi});

  @override
  List<Object?> get props => [sesi];
}

class SesiError extends SesiState {
  final String message;

  const SesiError({required this.message});

  @override
  List<Object?> get props => [message];
}
