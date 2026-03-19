part of 'sesi_bloc.dart';

sealed class SesiEvent extends Equatable {
  const SesiEvent();

  @override
  List<Object> get props => [];
}

class LoadSesiAktif extends SesiEvent {
  const LoadSesiAktif();
}

class LoadPecahanAktif extends SesiEvent {
  const LoadPecahanAktif();
}

class BukaSesiRequested extends SesiEvent {
  final List<Map<String, dynamic>> pecahan;

  const BukaSesiRequested({required this.pecahan});

  @override
  List<Object> get props => [pecahan];
}

class TutupSesiRequested extends SesiEvent {
  final String sesiId;
  final List<Map<String, dynamic>> pecahanAkhir;

  const TutupSesiRequested({
    required this.sesiId,
    required this.pecahanAkhir,
  });

  @override
  List<Object> get props => [sesiId, pecahanAkhir];
}
