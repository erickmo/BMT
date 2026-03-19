import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/dashboard_entity.dart';
import '../../../rekening/domain/usecases/get_daftar_rekening_usecase.dart';
import '../../../../core/storage/local_storage.dart';

part 'home_event.dart';
part 'home_state.dart';

class HomeBloc extends Bloc<HomeEvent, HomeState> {
  final GetDaftarRekeningUseCase getDaftarRekeningUseCase;
  final LocalStorage localStorage;

  HomeBloc({
    required this.getDaftarRekeningUseCase,
    required this.localStorage,
  }) : super(const HomeInitial()) {
    on<LoadDashboard>(_onLoadDashboard);
  }

  Future<void> _onLoadDashboard(
    LoadDashboard event,
    Emitter<HomeState> emit,
  ) async {
    emit(const HomeLoading());

    final result = await getDaftarRekeningUseCase();
    result.fold(
      (failure) => emit(HomeError(message: failure.message)),
      (rekening) {
        final totalSaldo = rekening.fold<int>(0, (sum, r) => sum + r.saldo);
        emit(
          HomeLoaded(
            dashboard: DashboardEntity(
              namaNasabah: localStorage.namaNasabah ?? 'Nasabah',
              nomorNasabah: localStorage.nomorNasabah ?? '',
              totalSaldo: totalSaldo,
              rekening: rekening,
            ),
          ),
        );
      },
    );
  }
}
