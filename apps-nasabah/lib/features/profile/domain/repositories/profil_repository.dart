import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../entities/profil_entity.dart';

abstract interface class ProfilRepository {
  Future<Either<Failure, ProfilEntity>> getProfil();
}
