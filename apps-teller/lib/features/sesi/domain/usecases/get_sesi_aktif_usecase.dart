import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../entities/sesi_entity.dart';
import '../repositories/sesi_repository.dart';

class GetSesiAktifUseCase {
  final SesiRepository repository;

  GetSesiAktifUseCase(this.repository);

  Future<Either<Failure, SesiEntity?>> call() {
    return repository.getSesiAktif();
  }
}
