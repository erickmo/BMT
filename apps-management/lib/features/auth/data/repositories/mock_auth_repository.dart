import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/storage/secure_storage.dart';
import '../../domain/entities/user_entity.dart';
import '../../domain/repositories/auth_repository.dart';

class MockAuthRepository implements AuthRepository {
  final SecureStorage secureStorage;

  MockAuthRepository({required this.secureStorage});

  static const _validUsername = 'mo@intinusa.id';
  static const _validPassword = '123123';

  static const _mockUser = UserEntity(
    id: '00000000-0000-0000-0000-000000000001',
    nama: 'Mo Intinusa',
    username: 'mo@intinusa.id',
    role: 'MANAJER_BMT',
    bmtId: '00000000-0000-0000-0000-000000000001',
    cabangId: null,
    accessToken: 'mock_access_token',
    refreshToken: 'mock_refresh_token',
  );

  @override
  Future<Either<Failure, UserEntity>> login({
    required String username,
    required String password,
  }) async {
    if (username != _validUsername || password != _validPassword) {
      return const Left(ServerFailure('Email atau password salah'));
    }
    try {
      await secureStorage.saveAccessToken(_mockUser.accessToken);
      await secureStorage.saveRefreshToken(_mockUser.refreshToken);
    } catch (_) {}
    return const Right(_mockUser);
  }

  @override
  Future<Either<Failure, void>> logout() async {
    try {
      await secureStorage.clearAll();
    } catch (_) {}
    return const Right(null);
  }

  @override
  Future<Either<Failure, UserEntity?>> getCurrentUser() async {
    return const Right(_mockUser);
  }
}
