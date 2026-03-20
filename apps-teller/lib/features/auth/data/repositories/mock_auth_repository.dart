import 'package:dartz/dartz.dart';

import '../../../../core/errors/failures.dart';
import '../../../../core/storage/local_storage.dart';
import '../../../../core/storage/secure_storage.dart';
import '../../domain/entities/staf_entity.dart';
import '../../domain/repositories/auth_repository.dart';

/// MockAuthRepository dipakai saat `--dart-define=MOCK_LOGIN=true`.
/// Bypass semua request API — langsung return StafEntity mock.
class MockAuthRepository implements AuthRepository {
  final LocalStorage localStorage;
  final SecureStorage secureStorage;

  static const _mockStaf = StafEntity(
    id: '00000000-0000-0000-0000-000000000001',
    nama: 'Mo Intinusa',
    email: 'mo@intinusa.id',
    role: 'TELLER',
    cabangId: '00000000-0000-0000-0000-000000000001',
    cabangNama: 'Cabang Utama',
    bmtId: '00000000-0000-0000-0000-000000000001',
    accessToken: 'mock_access_token',
    refreshToken: 'mock_refresh_token',
  );

  static const _validEmail = 'mo@intinusa.id';
  static const _validPassword = '123123';

  MockAuthRepository({
    required this.localStorage,
    required this.secureStorage,
  });

  Future<void> _seedLocalStorage() async {
    await localStorage.setLoggedIn(true);
    await localStorage.setNamaTeller(_mockStaf.nama);
    await localStorage.setRole(_mockStaf.role);
    await localStorage.setCabangId(_mockStaf.cabangId);
    try {
      await secureStorage.saveTokens(
        accessToken: _mockStaf.accessToken,
        refreshToken: _mockStaf.refreshToken,
      );
      await secureStorage.saveUserId(_mockStaf.id);
    } catch (_) {
      // Keychain mungkin tidak tersedia di mode debug sandbox
    }
  }

  @override
  Future<Either<Failure, StafEntity>> login({
    required String username,
    required String password,
    required String deviceId,
  }) async {
    if (username != _validEmail || password != _validPassword) {
      return const Left(ServerFailure(message: 'Email atau password salah'));
    }
    await _seedLocalStorage();
    return const Right(_mockStaf);
  }

  @override
  Future<Either<Failure, bool>> isLoggedIn() async {
    await _seedLocalStorage();
    return const Right(true);
  }

  @override
  Future<Either<Failure, void>> logout() async {
    try {
      await secureStorage.clearAll();
    } catch (_) {
      // Keychain mungkin tidak tersedia di mode debug sandbox
    }
    await localStorage.clear();
    return const Right(null);
  }
}
