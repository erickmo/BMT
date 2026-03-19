part of 'profil_bloc.dart';

sealed class ProfilEvent extends Equatable {
  const ProfilEvent();

  @override
  List<Object> get props => [];
}

class LoadProfil extends ProfilEvent {
  const LoadProfil();
}
