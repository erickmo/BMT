import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../entities/sesi_entity.dart';
import '../repositories/sesi_repository.dart';

class TutupSesiUseCase {
  final SesiRepository repository;

  TutupSesiUseCase(this.repository);

  Future<Either<Failure, SesiEntity>> call({
    required String sesiId,
    required List<Map<String, dynamic>> pecahanAkhir,
  }) {
    return repository.tutupSesi(
      sesiId: sesiId,
      pecahanAkhir: pecahanAkhir,
    );
  }
}
