import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/storage/local_storage.dart';
import '../../../../core/storage/secure_storage.dart';
import '../../domain/entities/staf_entity.dart';
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
  Future<Either<Failure, StafEntity>> login({
    required String username,
    required String password,
    required String deviceId,
  }) async {
    try {
      final model = await remoteDataSource.login(
        username: username,
        password: password,
        deviceId: deviceId,
      );

      await secureStorage.saveTokens(
        accessToken: model.accessToken,
        refreshToken: model.refreshToken,
      );
      await secureStorage.saveUserId(model.staf.id);
      await localStorage.setLoggedIn(true);
      await localStorage.setNamaTeller(model.staf.nama);
      await localStorage.setRole(model.staf.role);
      await localStorage.setCabangId(model.staf.cabangId);

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
    } catch (_) {}
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
