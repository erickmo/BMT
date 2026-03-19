import 'package:dartz/dartz.dart';
import '../../../../core/errors/failures.dart';
import '../entities/sesi_entity.dart';
import '../repositories/sesi_repository.dart';

class BukaSesiUseCase {
  final SesiRepository repository;

  BukaSesiUseCase(this.repository);

  Future<Either<Failure, SesiEntity>> call({
    required List<Map<String, dynamic>> pecahan,
  }) {
    return repository.bukaSesi(pecahan: pecahan);
  }
}
