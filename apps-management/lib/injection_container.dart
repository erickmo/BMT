import 'package:get_it/get_it.dart';
import 'package:dio/dio.dart';
import 'core/network/dio_client.dart';
import 'core/storage/secure_storage.dart';
import 'features/auth/data/datasources/auth_remote_datasource.dart';
import 'features/auth/data/repositories/auth_repository_impl.dart';
import 'features/auth/data/repositories/mock_auth_repository.dart';
import 'features/auth/domain/repositories/auth_repository.dart';
import 'features/auth/domain/usecases/login_usecase.dart';
import 'features/auth/presentation/bloc/auth_bloc.dart';

final sl = GetIt.instance;

const _mockLogin = bool.fromEnvironment('MOCK_LOGIN');

Future<void> initDependencies() async {
  // Core
  sl.registerLazySingleton<SecureStorage>(() => SecureStorage());
  sl.registerLazySingleton<DioClient>(() => DioClient(sl<SecureStorage>()));
  sl.registerLazySingleton<Dio>(() => sl<DioClient>().dio);

  // Auth
  if (!_mockLogin) {
    sl.registerLazySingleton<AuthRemoteDatasource>(
      () => AuthRemoteDatasourceImpl(sl<Dio>()),
    );
  }
  sl.registerLazySingleton<AuthRepository>(
    () => _mockLogin
        ? MockAuthRepository(secureStorage: sl<SecureStorage>())
        : AuthRepositoryImpl(
            sl<AuthRemoteDatasource>(),
            sl<SecureStorage>(),
          ),
  );
  sl.registerLazySingleton(() => LoginUsecase(sl<AuthRepository>()));
  sl.registerFactory(() => AuthBloc(sl<LoginUsecase>()));
}
