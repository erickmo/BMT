import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import '../../../../core/errors/failures.dart';
import '../entities/transaksi_entity.dart';
import '../repositories/rekening_repository.dart';

class GetRiwayatTransaksiUseCase {
  final RekeningRepository repository;

  GetRiwayatTransaksiUseCase(this.repository);

  Future<Either<Failure, List<TransaksiEntity>>> call(
    RiwayatTransaksiParams params,
  ) {
    return repository.getRiwayatTransaksi(
      params.rekeningId,
      page: params.page,
      limit: params.limit,
    );
  }
}

class RiwayatTransaksiParams extends Equatable {
  final String rekeningId;
  final int page;
  final int limit;

  const RiwayatTransaksiParams({
    required this.rekeningId,
    this.page = 1,
    this.limit = 20,
  });

  @override
  List<Object> get props => [rekeningId, page, limit];
}
