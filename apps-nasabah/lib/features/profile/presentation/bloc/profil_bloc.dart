import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/profil_entity.dart';
import '../../domain/usecases/get_profil_usecase.dart';

part 'profil_event.dart';
part 'profil_state.dart';

class ProfilBloc extends Bloc<ProfilEvent, ProfilState> {
  final GetProfilUseCase getProfilUseCase;

  ProfilBloc({required this.getProfilUseCase}) : super(const ProfilInitial()) {
    on<LoadProfil>(_onLoadProfil);
  }

  Future<void> _onLoadProfil(
    LoadProfil event,
    Emitter<ProfilState> emit,
  ) async {
    emit(const ProfilLoading());
    final result = await getProfilUseCase();
    result.fold(
      (failure) => emit(ProfilError(message: failure.message)),
      (profil) => emit(ProfilLoaded(profil: profil)),
    );
  }
}
