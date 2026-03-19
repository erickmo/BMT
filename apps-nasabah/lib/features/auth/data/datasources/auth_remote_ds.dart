import 'package:dio/dio.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/network/dio_client.dart';
import '../models/auth_response_model.dart';

abstract interface class AuthRemoteDataSource {
  Future<AuthResponseModel> login({
    required String nomorNasabah,
    required String pin,
    required String deviceId,
  });

  Future<AuthResponseModel> refreshToken({required String refreshToken});

  Future<void> logout();
}

class AuthRemoteDataSourceImpl implements AuthRemoteDataSource {
  final DioClient client;

  AuthRemoteDataSourceImpl(this.client);

  @override
  Future<AuthResponseModel> login({
    required String nomorNasabah,
    required String pin,
    required String deviceId,
  }) async {
    try {
      final response = await client.post('/auth/nasabah/login', data: {
        'nomor_nasabah': nomorNasabah,
        'pin': pin,
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
      throw _mapDioError(e);
    }
  }

  @override
  Future<AuthResponseModel> refreshToken({required String refreshToken}) async {
    try {
      final response = await client.post('/auth/nasabah/refresh', data: {
        'refresh_token': refreshToken,
      });

      final data = response.data as Map<String, dynamic>;
      return AuthResponseModel.fromJson(
        data['data'] as Map<String, dynamic>,
      );
    } on DioException catch (e) {
      throw _mapDioError(e);
    }
  }

  @override
  Future<void> logout() async {
    try {
      await client.post('/auth/nasabah/logout');
    } on DioException catch (e) {
      // If server responds 401, token already invalid — treat as success
      if (e.response?.statusCode == 401) return;
      throw _mapDioError(e);
    }
  }

  Failure _mapDioError(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.sendTimeout:
        return const TimeoutFailure();
      case DioExceptionType.connectionError:
        return const NetworkFailure();
      case DioExceptionType.badResponse:
        final statusCode = e.response?.statusCode;
        if (statusCode == 401) return const UnauthorizedFailure();
        final msg = e.response?.data?['message'] as String?;
        return ServerFailure(
          message: msg ?? 'Terjadi kesalahan server',
          statusCode: statusCode,
        );
      default:
        return const UnknownFailure();
    }
  }
}
