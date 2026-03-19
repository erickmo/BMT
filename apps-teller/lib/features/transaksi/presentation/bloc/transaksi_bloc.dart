import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/transaksi_entity.dart';
import '../../domain/repositories/transaksi_repository.dart';

part 'transaksi_event.dart';
part 'transaksi_state.dart';

class TransaksiBloc extends Bloc<TransaksiEvent, TransaksiState> {
  final TransaksiRepository repository;

  TransaksiBloc({required this.repository})
      : super(const TransaksiInitial()) {
    on<CariNasabah>(_onCariNasabah);
    on<SetorRequested>(_onSetor);
    on<TarikRequested>(_onTarik);
    on<ResetTransaksi>(_onReset);
  }

  Future<void> _onCariNasabah(
    CariNasabah event,
    Emitter<TransaksiState> emit,
  ) async {
    emit(const TransaksiLoading());
    final result = await repository.cariNasabah(event.query);
    result.fold(
      (failure) => emit(TransaksiError(message: failure.message)),
      (results) => emit(NasabahSearchLoaded(results: results)),
    );
  }

  Future<void> _onSetor(
    SetorRequested event,
    Emitter<TransaksiState> emit,
  ) async {
    emit(const TransaksiLoading());
    final result = await repository.setor(
      rekeningId: event.rekeningId,
      nominal: event.nominal,
      keterangan: event.keterangan,
    );
    result.fold(
      (failure) => emit(TransaksiError(message: failure.message)),
      (tx) => emit(TransaksiSuccess(result: tx)),
    );
  }

  Future<void> _onTarik(
    TarikRequested event,
    Emitter<TransaksiState> emit,
  ) async {
    emit(const TransaksiLoading());
    final result = await repository.tarik(
      rekeningId: event.rekeningId,
      nominal: event.nominal,
      keterangan: event.keterangan,
    );
    result.fold(
      (failure) => emit(TransaksiError(message: failure.message)),
      (tx) => emit(TransaksiSuccess(result: tx)),
    );
  }

  void _onReset(ResetTransaksi event, Emitter<TransaksiState> emit) {
    emit(const TransaksiInitial());
  }
}
