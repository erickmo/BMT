import 'package:dio/dio.dart';
import '../../storage/secure_storage.dart';

class AuthInterceptor extends Interceptor {
  final SecureStorage storage;
  final Dio dio;

  // Paths that do not need auth header
  static const List<String> _publicPaths = [
    '/auth/nasabah/login',
    '/auth/nasabah/refresh',
  ];

  AuthInterceptor({required this.storage, required this.dio});

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final path = options.path;
    final isPublic = _publicPaths.any((p) => path.contains(p));

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
      final requestPath = err.requestOptions.path;
      final isRefreshPath = requestPath.contains('/auth/nasabah/refresh');

      if (!isRefreshPath) {
        final refreshed = await _tryRefreshToken();
        if (refreshed) {
          final newToken = await storage.getAccessToken();
          final opts = err.requestOptions;
          opts.headers['Authorization'] = 'Bearer $newToken';
          try {
            final response = await dio.fetch(opts);
            handler.resolve(response);
            return;
          } catch (e) {
            // retry failed — fall through to error
          }
        }
        // Refresh failed — clear tokens so app navigates to login
        await storage.clearAll();
      }
    }
    handler.next(err);
  }

  Future<bool> _tryRefreshToken() async {
    try {
      final refreshToken = await storage.getRefreshToken();
      if (refreshToken == null) return false;

      final response = await dio.post(
        '/auth/nasabah/refresh',
        data: {'refresh_token': refreshToken},
      );

      final newAccessToken = response.data['data']['access_token'] as String?;
      final newRefreshToken = response.data['data']['refresh_token'] as String?;

      if (newAccessToken != null && newRefreshToken != null) {
        await storage.saveTokens(
          accessToken: newAccessToken,
          refreshToken: newRefreshToken,
        );
        return true;
      }
      return false;
    } catch (_) {
      return false;
    }
  }
}
