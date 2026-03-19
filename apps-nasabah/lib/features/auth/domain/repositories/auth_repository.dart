import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../entities/user_entity.dart';

abstract interface class AuthRepository {
  Future<Either<Failure, UserEntity>> login({
    required String nomorNasabah,
    required String pin,
    required String deviceId,
  });

  Future<Either<Failure, UserEntity>> refreshToken({
    required String refreshToken,
  });

  Future<Either<Failure, void>> logout();

  Future<Either<Failure, bool>> isLoggedIn();
}
