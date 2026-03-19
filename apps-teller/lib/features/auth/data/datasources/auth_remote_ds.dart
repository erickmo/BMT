import 'package:dio/dio.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/network/dio_client.dart';
import '../models/auth_response_model.dart';

abstract interface class AuthRemoteDataSource {
  Future<AuthResponseModel> login({
    required String username,
    required String password,
    required String deviceId,
  });

  Future<void> logout();
}

class AuthRemoteDataSourceImpl implements AuthRemoteDataSource {
  final DioClient client;

  AuthRemoteDataSourceImpl(this.client);

  @override
  Future<AuthResponseModel> login({
    required String username,
    required String password,
    required String deviceId,
  }) async {
    try {
      final response = await client.post('/auth/staf/login', data: {
        'username': username,
        'password': password,
        'device_id': deviceId,
      });

      final data = response.data as Map<String, dynamic>;
      if (data['success'] != true) {
        throw ServerFailure(
          message: data['message'] as String? ?? 'Login gagal',
          statusCode: response.statusCode,
        );
      }

      return AuthResponseModel.fromJson(
        data['data'] as Map<String, dynamic>,
      );
    } on DioException catch (e) {
      throw _mapError(e);
    }
  }

  @override
  Future<void> logout() async {
    try {
      await client.post('/auth/staf/logout');
    } on DioException catch (e) {
      if (e.response?.statusCode == 401) return;
      throw _mapError(e);
    }
  }

  Failure _mapError(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.sendTimeout:
        return const TimeoutFailure();
      case DioExceptionType.connectionError:
        return const NetworkFailure();
      case DioExceptionType.badResponse:
        final code = e.response?.statusCode;
        if (code == 401) return const UnauthorizedFailure();
        final msg = e.response?.data?['message'] as String?;
        return ServerFailure(message: msg ?? 'Terjadi kesalahan server', statusCode: code);
      default:
        return const UnknownFailure();
    }
  }
}
