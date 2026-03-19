import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../entities/transaksi_entity.dart';

abstract interface class TransaksiRepository {
  Future<Either<Failure, List<NasabahSearchResult>>> cariNasabah(String query);

  Future<Either<Failure, TransaksiResultEntity>> setor({
    required String rekeningId,
    required int nominal,
    required String keterangan,
  });

  Future<Either<Failure, TransaksiResultEntity>> tarik({
    required String rekeningId,
    required int nominal,
    required String keterangan,
  });
}
