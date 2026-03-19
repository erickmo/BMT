import 'package:dio/dio.dart';

class DioClient {
  static const String _baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );

  late final Dio _dio;

  DioClient() {
    _dio = Dio(
      BaseOptions(
        baseUrl: _baseUrl,
        connectTimeout: const Duration(seconds: 10),
        receiveTimeout: const Duration(seconds: 10),
        headers: {'Content-Type': 'application/json', 'Accept': 'application/json'},
      ),
    );
  }

  Dio get dio => _dio;

  /// GET /nfc/ceksaldo/:uid — no auth required, IP whitelisted by server
  Future<Map<String, dynamic>> cekSaldo(String uid) async {
    final response = await _dio.get('/nfc/ceksaldo/$uid');
    return response.data['data'] as Map<String, dynamic>;
  }
}
