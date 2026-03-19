import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class DioClient {
  static const String _baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );

  final FlutterSecureStorage _storage;
  late final Dio _dio;

  static const _keyToken = 'developer_token';

  DioClient(this._storage) {
    _dio = Dio(
      BaseOptions(
        baseUrl: _baseUrl,
        connectTimeout: const Duration(seconds: 30),
        receiveTimeout: const Duration(seconds: 30),
        headers: {'Content-Type': 'application/json', 'Accept': 'application/json'},
      ),
    );
    _dio.interceptors.add(_DevTokenInterceptor(_storage));
  }

  Dio get dio => _dio;

  Future<void> saveToken(String token) => _storage.write(key: _keyToken, value: token);
  Future<String?> getToken() => _storage.read(key: _keyToken);
  Future<void> clearToken() => _storage.delete(key: _keyToken);
  Future<bool> hasToken() async {
    final t = await getToken();
    return t != null && t.isNotEmpty;
  }
}

class _DevTokenInterceptor extends Interceptor {
  final FlutterSecureStorage _storage;
  static const _keyToken = 'developer_token';

  _DevTokenInterceptor(this._storage);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    final token = await _storage.read(key: _keyToken);
    if (token != null && token.isNotEmpty) {
      options.headers['Developer-Token'] = token;
    }
    handler.next(options);
  }
}
