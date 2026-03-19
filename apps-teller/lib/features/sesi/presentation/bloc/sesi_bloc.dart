import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/pecahan_entity.dart';
import '../../domain/entities/sesi_entity.dart';
import '../../domain/repositories/sesi_repository.dart';
import '../../domain/usecases/buka_sesi_usecase.dart';
import '../../domain/usecases/get_sesi_aktif_usecase.dart';
import '../../domain/usecases/tutup_sesi_usecase.dart';

part 'sesi_event.dart';
part 'sesi_state.dart';

class SesiBloc extends Bloc<SesiEvent, SesiState> {
  final GetSesiAktifUseCase getSesiAktifUseCase;
  final BukaSesiUseCase bukaSesiUseCase;
  final TutupSesiUseCase tutupSesiUseCase;
  final SesiRepository sesiRepository;

  SesiBloc({
    required this.getSesiAktifUseCase,
    required this.bukaSesiUseCase,
    required this.tutupSesiUseCase,
    required this.sesiRepository,
  }) : super(const SesiInitial()) {
    on<LoadSesiAktif>(_onLoadSesiAktif);
    on<LoadPecahanAktif>(_onLoadPecahan);
    on<BukaSesiRequested>(_onBukaSesi);
    on<TutupSesiRequested>(_onTutupSesi);
  }

  Future<void> _onLoadSesiAktif(
    LoadSesiAktif event,
    Emitter<SesiState> emit,
  ) async {
    emit(const SesiLoading());
    final result = await getSesiAktifUseCase();
    result.fold(
      (failure) => emit(SesiError(message: failure.message)),
      (sesi) => emit(SesiAktifLoaded(sesi: sesi)),
    );
  }

  Future<void> _onLoadPecahan(
    LoadPecahanAktif event,
    Emitter<SesiState> emit,
  ) async {
    emit(const SesiLoading());
    final result = await sesiRepository.getPecahanAktif();
    result.fold(
      (failure) => emit(SesiError(message: failure.message)),
      (pecahan) => emit(
        PecahanLoaded(
          pecahan: pecahan,
          jumlahMap: {for (final p in pecahan) p.id: 0},
        ),
      ),
    );
  }

  Future<void> _onBukaSesi(
    BukaSesiRequested event,
    Emitter<SesiState> emit,
  ) async {
    emit(const SesiLoading());
    final result = await bukaSesiUseCase(pecahan: event.pecahan);
    result.fold(
      (failure) => emit(SesiError(message: failure.message)),
      (sesi) => emit(SesiBukaSuccess(sesi: sesi)),
    );
  }

  Future<void> _onTutupSesi(
    TutupSesiRequested event,
    Emitter<SesiState> emit,
  ) async {
    emit(const SesiLoading());
    final result = await tutupSesiUseCase(
      sesiId: event.sesiId,
      pecahanAkhir: event.pecahanAkhir,
    );
    result.fold(
      (failure) => emit(SesiError(message: failure.message)),
      (sesi) => emit(SesiTutupSuccess(sesi: sesi)),
    );
  }
}
