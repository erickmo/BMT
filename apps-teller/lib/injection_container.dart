import 'package:get_it/get_it.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'core/network/dio_client.dart';
import 'core/storage/local_storage.dart';
import 'core/storage/secure_storage.dart';

// Auth
import 'features/auth/data/datasources/auth_remote_ds.dart';
import 'features/auth/data/repositories/auth_repository_impl.dart';
import 'features/auth/domain/repositories/auth_repository.dart';
import 'features/auth/domain/usecases/login_usecase.dart';
import 'features/auth/domain/usecases/logout_usecase.dart';
import 'features/auth/presentation/bloc/auth_bloc.dart';

// Sesi
import 'features/sesi/data/datasources/sesi_remote_ds.dart';
import 'features/sesi/data/repositories/sesi_repository_impl.dart';
import 'features/sesi/domain/repositories/sesi_repository.dart';
import 'features/sesi/domain/usecases/buka_sesi_usecase.dart';
import 'features/sesi/domain/usecases/get_sesi_aktif_usecase.dart';
import 'features/sesi/domain/usecases/tutup_sesi_usecase.dart';
import 'features/sesi/presentation/bloc/sesi_bloc.dart';

// Transaksi
import 'features/transaksi/data/datasources/transaksi_remote_ds.dart';
import 'features/transaksi/data/repositories/transaksi_repository_impl.dart';
import 'features/transaksi/domain/repositories/transaksi_repository.dart';
import 'features/transaksi/presentation/bloc/transaksi_bloc.dart';

final sl = GetIt.instance;

Future<void> initDependencies() async {
  // ── External ──────────────────────────────────────────────────────────────
  final prefs = await SharedPreferences.getInstance();
  sl.registerSingleton<SharedPreferences>(prefs);

  // ── Core ──────────────────────────────────────────────────────────────────
  sl.registerLazySingleton<SecureStorage>(() => SecureStorage());

  sl.registerLazySingleton<LocalStorage>(
    () => LocalStorage(sl<SharedPreferences>()),
  );

  const baseUrl = String.fromEnvironment(
    'API_URL',
    defaultValue: 'http://localhost:8080',
  );

  sl.registerLazySingleton<DioClient>(
    () => DioClient(baseUrl: baseUrl, storage: sl<SecureStorage>()),
  );

  // ── Auth ──────────────────────────────────────────────────────────────────
  sl.registerLazySingleton<AuthRemoteDataSource>(
    () => AuthRemoteDataSourceImpl(sl<DioClient>()),
  );

  sl.registerLazySingleton<AuthRepository>(
    () => AuthRepositoryImpl(
      remoteDataSource: sl<AuthRemoteDataSource>(),
      secureStorage: sl<SecureStorage>(),
      localStorage: sl<LocalStorage>(),
    ),
  );

  sl.registerLazySingleton(() => LoginUseCase(sl<AuthRepository>()));
  sl.registerLazySingleton(() => LogoutUseCase(sl<AuthRepository>()));

  sl.registerFactory(
    () => AuthBloc(
      loginUseCase: sl<LoginUseCase>(),
      logoutUseCase: sl<LogoutUseCase>(),
      authRepository: sl<AuthRepository>(),
    ),
  );

  // ── Sesi ──────────────────────────────────────────────────────────────────
  sl.registerLazySingleton<SesiRemoteDataSource>(
    () => SesiRemoteDataSourceImpl(sl<DioClient>()),
  );

  sl.registerLazySingleton<SesiRepository>(
    () => SesiRepositoryImpl(remoteDataSource: sl<SesiRemoteDataSource>()),
  );

  sl.registerLazySingleton(() => GetSesiAktifUseCase(sl<SesiRepository>()));
  sl.registerLazySingleton(() => BukaSesiUseCase(sl<SesiRepository>()));
  sl.registerLazySingleton(() => TutupSesiUseCase(sl<SesiRepository>()));

  sl.registerFactory(
    () => SesiBloc(
      getSesiAktifUseCase: sl<GetSesiAktifUseCase>(),
      bukaSesiUseCase: sl<BukaSesiUseCase>(),
      tutupSesiUseCase: sl<TutupSesiUseCase>(),
      sesiRepository: sl<SesiRepository>(),
    ),
  );

  // ── Transaksi ─────────────────────────────────────────────────────────────
  sl.registerLazySingleton<TransaksiRemoteDataSource>(
    () => TransaksiRemoteDataSourceImpl(sl<DioClient>()),
  );

  sl.registerLazySingleton<TransaksiRepository>(
    () => TransaksiRepositoryImpl(
        remoteDataSource: sl<TransaksiRemoteDataSource>()),
  );

  sl.registerFactory(
    () => TransaksiBloc(repository: sl<TransaksiRepository>()),
  );
}
