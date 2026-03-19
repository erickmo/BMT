part of 'rekening_bloc.dart';

sealed class RekeningState extends Equatable {
  const RekeningState();

  @override
  List<Object?> get props => [];
}

class RekeningInitial extends RekeningState {
  const RekeningInitial();
}

class RekeningLoading extends RekeningState {
  const RekeningLoading();
}

class DaftarRekeningLoaded extends RekeningState {
  final List<RekeningEntity> rekening;

  const DaftarRekeningLoaded({required this.rekening});

  @override
  List<Object?> get props => [rekening];
}

class DetailRekeningLoaded extends RekeningState {
  final RekeningEntity rekening;
  final List<TransaksiEntity> transaksi;
  final bool isLoadingMore;

  const DetailRekeningLoaded({
    required this.rekening,
    this.transaksi = const [],
    this.isLoadingMore = false,
  });

  DetailRekeningLoaded copyWith({
    RekeningEntity? rekening,
    List<TransaksiEntity>? transaksi,
    bool? isLoadingMore,
  }) {
    return DetailRekeningLoaded(
      rekening: rekening ?? this.rekening,
      transaksi: transaksi ?? this.transaksi,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
    );
  }

  @override
  List<Object?> get props => [rekening, transaksi, isLoadingMore];
}

class RekeningError extends RekeningState {
  final String message;

  const RekeningError({required this.message});

  @override
  List<Object?> get props => [message];
}
