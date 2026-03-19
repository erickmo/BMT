import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../entities/staf_entity.dart';

abstract interface class AuthRepository {
  Future<Either<Failure, StafEntity>> login({
    required String username,
    required String password,
    required String deviceId,
  });

  Future<Either<Failure, void>> logout();

  Future<Either<Failure, bool>> isLoggedIn();
}
