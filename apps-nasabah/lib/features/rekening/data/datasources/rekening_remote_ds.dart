import 'package:dio/dio.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/network/dio_client.dart';
import '../models/rekening_model.dart';

abstract interface class RekeningRemoteDataSource {
  Future<List<RekeningModel>> getDaftarRekening();
  Future<RekeningModel> getDetailRekening(String rekeningId);
  Future<List<TransaksiModel>> getRiwayatTransaksi(
    String rekeningId, {
    int page = 1,
    int limit = 20,
  });
}

class RekeningRemoteDataSourceImpl implements RekeningRemoteDataSource {
  final DioClient client;

  RekeningRemoteDataSourceImpl(this.client);

  @override
  Future<List<RekeningModel>> getDaftarRekening() async {
    try {
      final response = await client.get('/nasabah/rekening');
      final data = response.data as Map<String, dynamic>;
      final items = data['data'] as List<dynamic>? ?? [];
      return items
          .map((e) => RekeningModel.fromJson(e as Map<String, dynamic>))
          .toList();
    } on DioException catch (e) {
      throw _mapError(e);
    }
  }

  @override
  Future<RekeningModel> getDetailRekening(String rekeningId) async {
    try {
      final response = await client.get('/nasabah/rekening/$rekeningId');
      final data = response.data as Map<String, dynamic>;
      return RekeningModel.fromJson(data['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw _mapError(e);
    }
  }

  @override
  Future<List<TransaksiModel>> getRiwayatTransaksi(
    String rekeningId, {
    int page = 1,
    int limit = 20,
  }) async {
    try {
      final response = await client.get(
        '/nasabah/rekening/$rekeningId/transaksi',
        queryParameters: {'page': page, 'limit': limit},
      );
      final data = response.data as Map<String, dynamic>;
      final items =
          (data['data'] as Map<String, dynamic>?)?['items'] as List<dynamic>? ??
              [];
      return items
          .map((e) => TransaksiModel.fromJson(e as Map<String, dynamic>))
          .toList();
    } on DioException catch (e) {
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
