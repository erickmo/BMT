import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../../domain/entities/pecahan_entity.dart';
import '../../domain/entities/sesi_entity.dart';
import '../../domain/repositories/sesi_repository.dart';
import '../datasources/sesi_remote_ds.dart';

class SesiRepositoryImpl implements SesiRepository {
  final SesiRemoteDataSource remoteDataSource;

  SesiRepositoryImpl({required this.remoteDataSource});

  @override
  Future<Either<Failure, SesiEntity?>> getSesiAktif() async {
    try {
      final model = await remoteDataSource.getSesiAktif();
      return Right(model?.toEntity());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }

  @override
  Future<Either<Failure, List<PecahanEntity>>> getPecahanAktif() async {
    try {
      final models = await remoteDataSource.getPecahanAktif();
      return Right(models.map((m) => m.toEntity()).toList());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }

  @override
  Future<Either<Failure, SesiEntity>> bukaSesi({
    required List<Map<String, dynamic>> pecahan,
  }) async {
    try {
      final model = await remoteDataSource.bukaSesi(pecahan: pecahan);
      return Right(model.toEntity());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }

  @override
  Future<Either<Failure, SesiEntity>> tutupSesi({
    required String sesiId,
    required List<Map<String, dynamic>> pecahanAkhir,
  }) async {
    try {
      final model = await remoteDataSource.tutupSesi(
        sesiId: sesiId,
        pecahanAkhir: pecahanAkhir,
      );
      return Right(model.toEntity());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }
}
