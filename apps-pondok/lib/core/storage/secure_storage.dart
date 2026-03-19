import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class SecureStorage {
  static const _keyAccessToken = 'access_token';
  static const _keyRefreshToken = 'refresh_token';
  static const _keyUserRole = 'user_role';
  static const _keyUserId = 'user_id';
  static const _keyBmtId = 'bmt_id';
  static const _keyCabangId = 'cabang_id';

  final FlutterSecureStorage _storage;

  SecureStorage()
      : _storage = const FlutterSecureStorage(
          aOptions: AndroidOptions(encryptedSharedPreferences: true),
        );

  Future<void> saveAccessToken(String token) =>
      _storage.write(key: _keyAccessToken, value: token);
  Future<String?> getAccessToken() => _storage.read(key: _keyAccessToken);

  Future<void> saveRefreshToken(String token) =>
      _storage.write(key: _keyRefreshToken, value: token);
  Future<String?> getRefreshToken() => _storage.read(key: _keyRefreshToken);

  Future<void> saveUserRole(String role) =>
      _storage.write(key: _keyUserRole, value: role);
  Future<String?> getUserRole() => _storage.read(key: _keyUserRole);

  Future<void> saveUserId(String id) =>
      _storage.write(key: _keyUserId, value: id);
  Future<String?> getUserId() => _storage.read(key: _keyUserId);

  Future<void> saveBmtId(String id) =>
      _storage.write(key: _keyBmtId, value: id);
  Future<String?> getBmtId() => _storage.read(key: _keyBmtId);

  Future<void> saveCabangId(String id) =>
      _storage.write(key: _keyCabangId, value: id);
  Future<String?> getCabangId() => _storage.read(key: _keyCabangId);

  Future<void> saveAuthData({
    required String accessToken,
    required String refreshToken,
    required String userId,
    required String role,
    required String bmtId,
    String? cabangId,
  }) async {
    await Future.wait([
      saveAccessToken(accessToken),
      saveRefreshToken(refreshToken),
      saveUserId(userId),
      saveUserRole(role),
      saveBmtId(bmtId),
      if (cabangId != null) saveCabangId(cabangId),
    ]);
  }

  Future<void> clearAll() => _storage.deleteAll();

  Future<bool> isLoggedIn() async {
    final token = await getAccessToken();
    return token != null && token.isNotEmpty;
  }
}
