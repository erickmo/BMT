import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../entities/rekening_entity.dart';
import '../repositories/rekening_repository.dart';

class GetDaftarRekeningUseCase {
  final RekeningRepository repository;

  GetDaftarRekeningUseCase(this.repository);

  Future<Either<Failure, List<RekeningEntity>>> call() {
    return repository.getDaftarRekening();
  }
}
