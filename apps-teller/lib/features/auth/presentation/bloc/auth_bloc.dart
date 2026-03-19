import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/staf_entity.dart';
import '../../domain/usecases/login_usecase.dart';
import '../../domain/usecases/logout_usecase.dart';
import '../../domain/repositories/auth_repository.dart';

part 'auth_event.dart';
part 'auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final LoginUseCase loginUseCase;
  final LogoutUseCase logoutUseCase;
  final AuthRepository authRepository;

  AuthBloc({
    required this.loginUseCase,
    required this.logoutUseCase,
    required this.authRepository,
  }) : super(const AuthInitial()) {
    on<AuthCheckRequested>(_onCheck);
    on<LoginRequested>(_onLogin);
    on<LogoutRequested>(_onLogout);
  }

  Future<void> _onCheck(
      AuthCheckRequested event, Emitter<AuthState> emit) async {
    emit(const AuthLoading());
    final result = await authRepository.isLoggedIn();
    result.fold(
      (_) => emit(const AuthUnauthenticated()),
      (loggedIn) => loggedIn
          ? null // stay loading — full profile fetch not wired here
          : emit(const AuthUnauthenticated()),
    );
  }

  Future<void> _onLogin(
      LoginRequested event, Emitter<AuthState> emit) async {
    emit(const AuthLoading());
    final result = await loginUseCase(
      LoginParams(
        username: event.username,
        password: event.password,
        deviceId: event.deviceId,
      ),
    );
    result.fold(
      (failure) => emit(AuthFailure(message: failure.message)),
      (staf) => emit(AuthAuthenticated(staf: staf)),
    );
  }

  Future<void> _onLogout(
      LogoutRequested event, Emitter<AuthState> emit) async {
    emit(const AuthLoading());
    await logoutUseCase();
    emit(const AuthUnauthenticated());
  }
}
