import 'package:dio/dio.dart';
import '../../../../core/errors/failures.dart';
import '../models/login_response_model.dart';

abstract class AuthRemoteDatasource {
  Future<LoginResponseModel> login({
    required String username,
    required String password,
  });

  Future<void> logout(String refreshToken);
}

class AuthRemoteDatasourceImpl implements AuthRemoteDatasource {
  final Dio _dio;
  AuthRemoteDatasourceImpl(this._dio);

  @override
  Future<LoginResponseModel> login({
    required String username,
    required String password,
  }) async {
    try {
      final response = await _dio.post(
        '/auth/staf/login',
        data: {'username': username, 'password': password},
      );
      final data = response.data['data'] as Map<String, dynamic>;
      return LoginResponseModel.fromJson(data);
    } on DioException catch (e) {
      if (e.response?.statusCode == 401) {
        throw const UnauthorizedFailure('Username atau password salah');
      }
      throw ServerFailure(
        e.response?.data?['message'] ?? 'Gagal login',
        statusCode: e.response?.statusCode,
      );
    }
  }

  @override
  Future<void> logout(String refreshToken) async {
    try {
      await _dio.post(
        '/auth/staf/logout',
        data: {'refresh_token': refreshToken},
      );
    } on DioException catch (_) {
      // Ignore logout errors — clear locally regardless
    }
  }
}
