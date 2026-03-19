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

// Rekening
import 'features/rekening/data/datasources/rekening_remote_ds.dart';
import 'features/rekening/data/repositories/rekening_repository_impl.dart';
import 'features/rekening/domain/repositories/rekening_repository.dart';
import 'features/rekening/domain/usecases/get_daftar_rekening_usecase.dart';
import 'features/rekening/domain/usecases/get_riwayat_transaksi_usecase.dart';
import 'features/rekening/presentation/bloc/rekening_bloc.dart';

// Home
import 'features/home/presentation/bloc/home_bloc.dart';

// Profile
import 'features/profile/data/datasources/profil_remote_ds.dart';
import 'features/profile/data/repositories/profil_repository_impl.dart';
import 'features/profile/domain/repositories/profil_repository.dart';
import 'features/profile/domain/usecases/get_profil_usecase.dart';
import 'features/profile/presentation/bloc/profil_bloc.dart';

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

  // ── Rekening ──────────────────────────────────────────────────────────────
  sl.registerLazySingleton<RekeningRemoteDataSource>(
    () => RekeningRemoteDataSourceImpl(sl<DioClient>()),
  );

  sl.registerLazySingleton<RekeningRepository>(
    () => RekeningRepositoryImpl(remoteDataSource: sl<RekeningRemoteDataSource>()),
  );

  sl.registerLazySingleton(
    () => GetDaftarRekeningUseCase(sl<RekeningRepository>()),
  );
  sl.registerLazySingleton(
    () => GetRiwayatTransaksiUseCase(sl<RekeningRepository>()),
  );

  sl.registerFactory(
    () => RekeningBloc(
      getDaftarRekeningUseCase: sl<GetDaftarRekeningUseCase>(),
      getRiwayatTransaksiUseCase: sl<GetRiwayatTransaksiUseCase>(),
      rekeningRepository: sl<RekeningRepository>(),
    ),
  );

  // ── Home ──────────────────────────────────────────────────────────────────
  sl.registerFactory(
    () => HomeBloc(
      getDaftarRekeningUseCase: sl<GetDaftarRekeningUseCase>(),
      localStorage: sl<LocalStorage>(),
    ),
  );

  // ── Profile ───────────────────────────────────────────────────────────────
  sl.registerLazySingleton<ProfilRemoteDataSource>(
    () => ProfilRemoteDataSourceImpl(sl<DioClient>()),
  );

  sl.registerLazySingleton<ProfilRepository>(
    () => ProfilRepositoryImpl(remoteDataSource: sl<ProfilRemoteDataSource>()),
  );

  sl.registerLazySingleton(() => GetProfilUseCase(sl<ProfilRepository>()));

  sl.registerFactory(
    () => ProfilBloc(getProfilUseCase: sl<GetProfilUseCase>()),
  );
}
