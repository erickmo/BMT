part of 'profil_bloc.dart';

sealed class ProfilState extends Equatable {
  const ProfilState();

  @override
  List<Object?> get props => [];
}

class ProfilInitial extends ProfilState {
  const ProfilInitial();
}

class ProfilLoading extends ProfilState {
  const ProfilLoading();
}

class ProfilLoaded extends ProfilState {
  final ProfilEntity profil;

  const ProfilLoaded({required this.profil});

  @override
  List<Object?> get props => [profil];
}

class ProfilError extends ProfilState {
  final String message;

  const ProfilError({required this.message});

  @override
  List<Object?> get props => [message];
}
