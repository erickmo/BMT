part of 'transaksi_bloc.dart';

sealed class TransaksiState extends Equatable {
  const TransaksiState();

  @override
  List<Object?> get props => [];
}

class TransaksiInitial extends TransaksiState {
  const TransaksiInitial();
}

class TransaksiLoading extends TransaksiState {
  const TransaksiLoading();
}

class NasabahSearchLoaded extends TransaksiState {
  final List<NasabahSearchResult> results;

  const NasabahSearchLoaded({required this.results});

  @override
  List<Object?> get props => [results];
}

class TransaksiSuccess extends TransaksiState {
  final TransaksiResultEntity result;

  const TransaksiSuccess({required this.result});

  @override
  List<Object?> get props => [result];
}

class TransaksiError extends TransaksiState {
  final String message;

  const TransaksiError({required this.message});

  @override
  List<Object?> get props => [message];
}
