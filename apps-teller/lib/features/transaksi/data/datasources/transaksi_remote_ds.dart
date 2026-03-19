import 'package:dio/dio.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/network/dio_client.dart';
import '../models/transaksi_model.dart';

abstract interface class TransaksiRemoteDataSource {
  Future<List<NasabahSearchModel>> cariNasabah(String query);
  Future<TransaksiResultModel> setor({
    required String rekeningId,
    required int nominal,
    required String keterangan,
  });
  Future<TransaksiResultModel> tarik({
    required String rekeningId,
    required int nominal,
    required String keterangan,
  });
}

class TransaksiRemoteDataSourceImpl implements TransaksiRemoteDataSource {
  final DioClient client;

  TransaksiRemoteDataSourceImpl(this.client);

  @override
  Future<List<NasabahSearchModel>> cariNasabah(String query) async {
    try {
      final response = await client.get(
        '/teller/nasabah/cari',
        queryParameters: {'q': query},
      );
      final data = response.data as Map<String, dynamic>;
      final items = data['data'] as List<dynamic>? ?? [];
      return items
          .map((e) => NasabahSearchModel.fromJson(e as Map<String, dynamic>))
          .toList();
    } on DioException catch (e) {
      throw _mapError(e);
    }
  }

  @override
  Future<TransaksiResultModel> setor({
    required String rekeningId,
    required int nominal,
    required String keterangan,
  }) async {
    try {
      final response = await client.post(
        '/teller/rekening/$rekeningId/setor',
        data: {'nominal': nominal, 'keterangan': keterangan},
      );
      final data = response.data as Map<String, dynamic>;
      return TransaksiResultModel.fromJson(
          data['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw _mapError(e);
    }
  }

  @override
  Future<TransaksiResultModel> tarik({
    required String rekeningId,
    required int nominal,
    required String keterangan,
  }) async {
    try {
      final response = await client.post(
        '/teller/rekening/$rekeningId/tarik',
        data: {'nominal': nominal, 'keterangan': keterangan},
      );
      final data = response.data as Map<String, dynamic>;
      return TransaksiResultModel.fromJson(
          data['data'] as Map<String, dynamic>);
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
          message:
              e.response?.data?['message'] as String? ?? 'Terjadi kesalahan',
          statusCode: e.response?.statusCode,
        );
    }
  }
}
