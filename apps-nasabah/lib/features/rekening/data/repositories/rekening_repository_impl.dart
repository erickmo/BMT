import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../../domain/entities/rekening_entity.dart';
import '../../domain/entities/transaksi_entity.dart';
import '../../domain/repositories/rekening_repository.dart';
import '../datasources/rekening_remote_ds.dart';

class RekeningRepositoryImpl implements RekeningRepository {
  final RekeningRemoteDataSource remoteDataSource;

  RekeningRepositoryImpl({required this.remoteDataSource});

  @override
  Future<Either<Failure, List<RekeningEntity>>> getDaftarRekening() async {
    try {
      final models = await remoteDataSource.getDaftarRekening();
      return Right(models.map((m) => m.toEntity()).toList());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }

  @override
  Future<Either<Failure, RekeningEntity>> getDetailRekening(
    String rekeningId,
  ) async {
    try {
      final model = await remoteDataSource.getDetailRekening(rekeningId);
      return Right(model.toEntity());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }

  @override
  Future<Either<Failure, List<TransaksiEntity>>> getRiwayatTransaksi(
    String rekeningId, {
    int page = 1,
    int limit = 20,
  }) async {
    try {
      final models = await remoteDataSource.getRiwayatTransaksi(
        rekeningId,
        page: page,
        limit: limit,
      );
      return Right(models.map((m) => m.toEntity()).toList());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }
}
