import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/rekening_entity.dart';
import '../../domain/entities/transaksi_entity.dart';
import '../../domain/usecases/get_daftar_rekening_usecase.dart';
import '../../domain/usecases/get_riwayat_transaksi_usecase.dart';
import '../../domain/repositories/rekening_repository.dart';

part 'rekening_event.dart';
part 'rekening_state.dart';

class RekeningBloc extends Bloc<RekeningEvent, RekeningState> {
  final GetDaftarRekeningUseCase getDaftarRekeningUseCase;
  final GetRiwayatTransaksiUseCase getRiwayatTransaksiUseCase;
  final RekeningRepository rekeningRepository;

  RekeningBloc({
    required this.getDaftarRekeningUseCase,
    required this.getRiwayatTransaksiUseCase,
    required this.rekeningRepository,
  }) : super(const RekeningInitial()) {
    on<LoadDaftarRekening>(_onLoadDaftar);
    on<LoadDetailRekening>(_onLoadDetail);
    on<LoadRiwayatTransaksi>(_onLoadRiwayat);
  }

  Future<void> _onLoadDaftar(
    LoadDaftarRekening event,
    Emitter<RekeningState> emit,
  ) async {
    emit(const RekeningLoading());
    final result = await getDaftarRekeningUseCase();
    result.fold(
      (failure) => emit(RekeningError(message: failure.message)),
      (rekening) => emit(DaftarRekeningLoaded(rekening: rekening)),
    );
  }

  Future<void> _onLoadDetail(
    LoadDetailRekening event,
    Emitter<RekeningState> emit,
  ) async {
    emit(const RekeningLoading());
    final result = await rekeningRepository.getDetailRekening(event.rekeningId);
    result.fold(
      (failure) => emit(RekeningError(message: failure.message)),
      (rekening) => emit(DetailRekeningLoaded(rekening: rekening)),
    );
  }

  Future<void> _onLoadRiwayat(
    LoadRiwayatTransaksi event,
    Emitter<RekeningState> emit,
  ) async {
    if (state is DetailRekeningLoaded) {
      final current = state as DetailRekeningLoaded;
      emit(current.copyWith(isLoadingMore: true));
    }
    final result = await getRiwayatTransaksiUseCase(
      RiwayatTransaksiParams(
        rekeningId: event.rekeningId,
        page: event.page,
      ),
    );

    result.fold(
      (failure) => emit(RekeningError(message: failure.message)),
      (transaksi) {
        if (state is DetailRekeningLoaded) {
          final current = state as DetailRekeningLoaded;
          final all = event.page == 1
              ? transaksi
              : [...current.transaksi, ...transaksi];
          emit(current.copyWith(transaksi: all, isLoadingMore: false));
        }
      },
    );
  }
}
