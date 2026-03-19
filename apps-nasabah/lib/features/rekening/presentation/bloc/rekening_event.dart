part of 'rekening_bloc.dart';

sealed class RekeningEvent extends Equatable {
  const RekeningEvent();

  @override
  List<Object> get props => [];
}

class LoadDaftarRekening extends RekeningEvent {
  const LoadDaftarRekening();
}

class LoadDetailRekening extends RekeningEvent {
  final String rekeningId;

  const LoadDetailRekening(this.rekeningId);

  @override
  List<Object> get props => [rekeningId];
}

class LoadRiwayatTransaksi extends RekeningEvent {
  final String rekeningId;
  final int page;

  const LoadRiwayatTransaksi({required this.rekeningId, this.page = 1});

  @override
  List<Object> get props => [rekeningId, page];
}
