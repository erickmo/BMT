part of 'transaksi_bloc.dart';

sealed class TransaksiEvent extends Equatable {
  const TransaksiEvent();

  @override
  List<Object> get props => [];
}

class CariNasabah extends TransaksiEvent {
  final String query;

  const CariNasabah(this.query);

  @override
  List<Object> get props => [query];
}

class SetorRequested extends TransaksiEvent {
  final String rekeningId;
  final int nominal;
  final String keterangan;

  const SetorRequested({
    required this.rekeningId,
    required this.nominal,
    required this.keterangan,
  });

  @override
  List<Object> get props => [rekeningId, nominal, keterangan];
}

class TarikRequested extends TransaksiEvent {
  final String rekeningId;
  final int nominal;
  final String keterangan;

  const TarikRequested({
    required this.rekeningId,
    required this.nominal,
    required this.keterangan,
  });

  @override
  List<Object> get props => [rekeningId, nominal, keterangan];
}

class ResetTransaksi extends TransaksiEvent {
  const ResetTransaksi();
}
