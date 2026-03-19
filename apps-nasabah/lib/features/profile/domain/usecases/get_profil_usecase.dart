import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../entities/profil_entity.dart';
import '../repositories/profil_repository.dart';

class GetProfilUseCase {
  final ProfilRepository repository;

  GetProfilUseCase(this.repository);

  Future<Either<Failure, ProfilEntity>> call() {
    return repository.getProfil();
  }
}
