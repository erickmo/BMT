import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/storage/local_storage.dart';
import '../../../../core/storage/secure_storage.dart';
import '../../domain/entities/user_entity.dart';
import '../../domain/repositories/auth_repository.dart';
import '../datasources/auth_remote_ds.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthRemoteDataSource remoteDataSource;
  final SecureStorage secureStorage;
  final LocalStorage localStorage;

  AuthRepositoryImpl({
    required this.remoteDataSource,
    required this.secureStorage,
    required this.localStorage,
  });

  @override
  Future<Either<Failure, UserEntity>> login({
    required String nomorNasabah,
    required String pin,
    required String deviceId,
  }) async {
    try {
      final model = await remoteDataSource.login(
        nomorNasabah: nomorNasabah,
        pin: pin,
        deviceId: deviceId,
      );

      // Persist tokens and user info locally
      await secureStorage.saveTokens(
        accessToken: model.accessToken,
        refreshToken: model.refreshToken,
      );
      await secureStorage.saveUserId(model.nasabah.id);
      await localStorage.setLoggedIn(true);
      await localStorage.setNomorNasabah(model.nasabah.nomorNasabah);
      await localStorage.setNamaNasabah(model.nasabah.nama);

      return Right(model.toEntity());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }

  @override
  Future<Either<Failure, UserEntity>> refreshToken({
    required String refreshToken,
  }) async {
    try {
      final model = await remoteDataSource.refreshToken(
        refreshToken: refreshToken,
      );
      await secureStorage.saveTokens(
        accessToken: model.accessToken,
        refreshToken: model.refreshToken,
      );
      return Right(model.toEntity());
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return Left(UnknownFailure(message: e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> logout() async {
    try {
      await remoteDataSource.logout();
    } catch (_) {
      // best-effort server logout
    }
    await secureStorage.clearAll();
    await localStorage.clear();
    return const Right(null);
  }

  @override
  Future<Either<Failure, bool>> isLoggedIn() async {
    final loggedIn = localStorage.isLoggedIn;
    final token = await secureStorage.getAccessToken();
    return Right(loggedIn && token != null);
  }
}
