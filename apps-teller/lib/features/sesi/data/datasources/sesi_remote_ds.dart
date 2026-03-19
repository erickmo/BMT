import 'package:dio/dio.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/network/dio_client.dart';
import '../models/sesi_model.dart';

abstract interface class SesiRemoteDataSource {
  Future<SesiModel?> getSesiAktif();
  Future<List<PecahanModel>> getPecahanAktif();
  Future<SesiModel> bukaSesi({required List<Map<String, dynamic>> pecahan});
  Future<SesiModel> tutupSesi({
    required String sesiId,
    required List<Map<String, dynamic>> pecahanAkhir,
  });
}

class SesiRemoteDataSourceImpl implements SesiRemoteDataSource {
  final DioClient client;

  SesiRemoteDataSourceImpl(this.client);

  @override
  Future<SesiModel?> getSesiAktif() async {
    try {
      final response = await client.get('/teller/sesi/aktif');
      final data = response.data as Map<String, dynamic>;
      if (data['data'] == null) return null;
      return SesiModel.fromJson(data['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      if (e.response?.statusCode == 404) return null;
      throw _mapError(e);
    }
  }

  @override
  Future<List<PecahanModel>> getPecahanAktif() async {
    try {
      final response = await client.get('/dev/pecahan-uang');
      final data = response.data as Map<String, dynamic>;
      final items = data['data'] as List<dynamic>? ?? [];
      return items
          .map((e) => PecahanModel.fromJson(e as Map<String, dynamic>))
          .toList();
    } on DioException catch (e) {
      throw _mapError(e);
    }
  }

  @override
  Future<SesiModel> bukaSesi({
    required List<Map<String, dynamic>> pecahan,
  }) async {
    try {
      final response = await client.post(
        '/teller/sesi/buka',
        data: {'pecahan': pecahan},
      );
      final data = response.data as Map<String, dynamic>;
      return SesiModel.fromJson(data['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw _mapError(e);
    }
  }

  @override
  Future<SesiModel> tutupSesi({
    required String sesiId,
    required List<Map<String, dynamic>> pecahanAkhir,
  }) async {
    try {
      final response = await client.post(
        '/teller/sesi/tutup',
        data: {'sesi_id': sesiId, 'pecahan_akhir': pecahanAkhir},
      );
      final data = response.data as Map<String, dynamic>;
      return SesiModel.fromJson(data['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw _mapError(e);
    }
  }

  Failure _mapError(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionError:
        return const NetworkFailure();
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.receiveTimeout:
        return const TimeoutFailure();
      default:
        if (e.response?.statusCode == 401) return const UnauthorizedFailure();
        return ServerFailure(
          message: e.response?.data?['message'] as String? ?? 'Terjadi kesalahan',
          statusCode: e.response?.statusCode,
        );
    }
  }
}
