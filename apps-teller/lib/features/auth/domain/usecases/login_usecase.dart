import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import '../../../../core/errors/failures.dart';
import '../entities/staf_entity.dart';
import '../repositories/auth_repository.dart';

class LoginUseCase {
  final AuthRepository repository;

  LoginUseCase(this.repository);

  Future<Either<Failure, StafEntity>> call(LoginParams params) {
    return repository.login(
      username: params.username,
      password: params.password,
      deviceId: params.deviceId,
    );
  }
}

class LoginParams extends Equatable {
  final String username;
  final String password;
  final String deviceId;

  const LoginParams({
    required this.username,
    required this.password,
    required this.deviceId,
  });

  @override
  List<Object> get props => [username, password, deviceId];
}
