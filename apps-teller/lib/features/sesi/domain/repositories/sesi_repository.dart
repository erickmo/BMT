import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../entities/pecahan_entity.dart';
import '../entities/sesi_entity.dart';

abstract interface class SesiRepository {
  Future<Either<Failure, SesiEntity?>> getSesiAktif();

  Future<Either<Failure, List<PecahanEntity>>> getPecahanAktif();

  Future<Either<Failure, SesiEntity>> bukaSesi({
    required List<Map<String, dynamic>> pecahan,
  });

  Future<Either<Failure, SesiEntity>> tutupSesi({
    required String sesiId,
    required List<Map<String, dynamic>> pecahanAkhir,
  });
}
