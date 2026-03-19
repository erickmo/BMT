import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../../domain/entities/profil_entity.dart';
import '../../domain/repositories/profil_repository.dart';
import '../datasources/profil_remote_ds.dart';

class ProfilRepositoryImpl implements ProfilRepository {
  final ProfilRemoteDataSource remoteDataSource;

  ProfilRepositoryImpl({required this.remoteDataSource});

  @override
  Future<Either<Failure, ProfilEntity>> getProfil() async {
    try {
      final model = await remoteDataSource.getProfil();
      return Right(model.toEntity());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }
}
