import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../../domain/entities/transaksi_entity.dart';
import '../../domain/repositories/transaksi_repository.dart';
import '../datasources/transaksi_remote_ds.dart';

class TransaksiRepositoryImpl implements TransaksiRepository {
  final TransaksiRemoteDataSource remoteDataSource;

  TransaksiRepositoryImpl({required this.remoteDataSource});

  @override
  Future<Either<Failure, List<NasabahSearchResult>>> cariNasabah(
    String query,
  ) async {
    try {
      final models = await remoteDataSource.cariNasabah(query);
      return Right(models.map((m) => m.toEntity()).toList());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }

  @override
  Future<Either<Failure, TransaksiResultEntity>> setor({
    required String rekeningId,
    required int nominal,
    required String keterangan,
  }) async {
    try {
      final model = await remoteDataSource.setor(
        rekeningId: rekeningId,
        nominal: nominal,
        keterangan: keterangan,
      );
      return Right(model.toEntity());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }

  @override
  Future<Either<Failure, TransaksiResultEntity>> tarik({
    required String rekeningId,
    required int nominal,
    required String keterangan,
  }) async {
    try {
      final model = await remoteDataSource.tarik(
        rekeningId: rekeningId,
        nominal: nominal,
        keterangan: keterangan,
      );
      return Right(model.toEntity());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }
}
