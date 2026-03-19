import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../entities/rekening_entity.dart';
import '../entities/transaksi_entity.dart';

abstract interface class RekeningRepository {
  Future<Either<Failure, List<RekeningEntity>>> getDaftarRekening();

  Future<Either<Failure, RekeningEntity>> getDetailRekening(String rekeningId);

  Future<Either<Failure, List<TransaksiEntity>>> getRiwayatTransaksi(
    String rekeningId, {
    int page = 1,
    int limit = 20,
  });
}
