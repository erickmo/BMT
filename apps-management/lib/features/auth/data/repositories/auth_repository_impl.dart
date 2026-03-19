import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/storage/secure_storage.dart';
import '../../domain/entities/user_entity.dart';
import '../../domain/repositories/auth_repository.dart';
import '../datasources/auth_remote_datasource.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthRemoteDatasource _remote;
  final SecureStorage _secureStorage;

  AuthRepositoryImpl(this._remote, this._secureStorage);

  @override
  Future<Either<Failure, UserEntity>> login({
    required String username,
    required String password,
  }) async {
    try {
      final model = await _remote.login(username: username, password: password);
      final entity = model.toEntity();
      await _secureStorage.saveAuthData(
        accessToken: entity.accessToken,
        refreshToken: entity.refreshToken,
        userId: entity.id,
        role: entity.role,
        bmtId: entity.bmtId,
        cabangId: entity.cabangId,
      );
      return Right(entity);
    } on Failure catch (f) {
      return Left(f);
    } catch (e) {
      return const Left(UnexpectedFailure());
    }
  }

  @override
  Future<Either<Failure, void>> logout() async {
    try {
      final refreshToken = await _secureStorage.getRefreshToken();
      if (refreshToken != null) {
        await _remote.logout(refreshToken);
      }
      await _secureStorage.clearAll();
      return const Right(null);
    } catch (_) {
      await _secureStorage.clearAll();
      return const Right(null);
    }
  }

  @override
  Future<Either<Failure, UserEntity?>> getCurrentUser() async {
    try {
      final isLoggedIn = await _secureStorage.isLoggedIn();
      if (!isLoggedIn) return const Right(null);

      final id = await _secureStorage.getUserId();
      final role = await _secureStorage.getUserRole();
      final bmtId = await _secureStorage.getBmtId();
      final cabangId = await _secureStorage.getCabangId();
      final accessToken = await _secureStorage.getAccessToken();
      final refreshToken = await _secureStorage.getRefreshToken();

      if (id == null || role == null || bmtId == null) return const Right(null);

      return Right(UserEntity(
        id: id,
        nama: '',
        username: '',
        role: role,
        bmtId: bmtId,
        cabangId: cabangId,
        accessToken: accessToken ?? '',
        refreshToken: refreshToken ?? '',
      ));
    } catch (_) {
      return const Left(CacheFailure());
    }
  }
}
