import 'package:dio/dio.dart';
import '../../storage/secure_storage.dart';

class AuthInterceptor extends Interceptor {
  final SecureStorage storage;
  final Dio dio;

  static const List<String> _publicPaths = [
    '/auth/staf/login',
    '/auth/staf/refresh',
  ];

  AuthInterceptor({required this.storage, required this.dio});

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final isPublic = _publicPaths.any((p) => options.path.contains(p));
    if (!isPublic) {
      final token = await storage.getAccessToken();
      if (token != null) {
        options.headers['Authorization'] = 'Bearer $token';
      }
    }
    handler.next(options);
  }

  @override
  Future<void> onError(
    DioException err,
    ErrorInterceptorHandler handler,
  ) async {
    if (err.response?.statusCode == 401) {
      final isRefresh = err.requestOptions.path.contains('/auth/staf/refresh');
      if (!isRefresh) {
        final refreshed = await _tryRefresh();
        if (refreshed) {
          final token = await storage.getAccessToken();
          final opts = err.requestOptions;
          opts.headers['Authorization'] = 'Bearer $token';
          try {
            final res = await dio.fetch(opts);
            handler.resolve(res);
            return;
          } catch (_) {}
        }
        await storage.clearAll();
      }
    }
    handler.next(err);
  }

  Future<bool> _tryRefresh() async {
    try {
      final refresh = await storage.getRefreshToken();
      if (refresh == null) return false;
      final res = await dio.post(
        '/auth/staf/refresh',
        data: {'refresh_token': refresh},
      );
      final newAccess = res.data['data']['access_token'] as String?;
      final newRefresh = res.data['data']['refresh_token'] as String?;
      if (newAccess != null && newRefresh != null) {
        await storage.saveTokens(
            accessToken: newAccess, refreshToken: newRefresh);
        return true;
      }
      return false;
    } catch (_) {
      return false;
    }
  }
}
