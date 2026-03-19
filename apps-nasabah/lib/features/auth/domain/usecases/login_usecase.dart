import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import '../../../../core/errors/failures.dart';
import '../entities/user_entity.dart';
import '../repositories/auth_repository.dart';

class LoginUseCase {
  final AuthRepository repository;

  LoginUseCase(this.repository);

  Future<Either<Failure, UserEntity>> call(LoginParams params) {
    return repository.login(
      nomorNasabah: params.nomorNasabah,
      pin: params.pin,
      deviceId: params.deviceId,
    );
  }
}

class LoginParams extends Equatable {
  final String nomorNasabah;
  final String pin;
  final String deviceId;

  const LoginParams({
    required this.nomorNasabah,
    required this.pin,
    required this.deviceId,
  });

  @override
  List<Object> get props => [nomorNasabah, pin, deviceId];
}
