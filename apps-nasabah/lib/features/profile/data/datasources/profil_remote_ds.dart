import 'package:dio/dio.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/network/dio_client.dart';
import '../models/profil_model.dart';

abstract interface class ProfilRemoteDataSource {
  Future<ProfilModel> getProfil();
}

class ProfilRemoteDataSourceImpl implements ProfilRemoteDataSource {
  final DioClient client;

  ProfilRemoteDataSourceImpl(this.client);

  @override
  Future<ProfilModel> getProfil() async {
    try {
      final response = await client.get('/nasabah/profil');
      final data = response.data as Map<String, dynamic>;
      return ProfilModel.fromJson(data['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      switch (e.type) {
        case DioExceptionType.connectionError:
          throw const NetworkFailure();
        case DioExceptionType.connectionTimeout:
        case DioExceptionType.receiveTimeout:
          throw const TimeoutFailure();
        default:
          if (e.response?.statusCode == 401) throw const UnauthorizedFailure();
          throw ServerFailure(
            message:
                e.response?.data?['message'] as String? ?? 'Terjadi kesalahan',
            statusCode: e.response?.statusCode,
          );
      }
    }
  }
}
