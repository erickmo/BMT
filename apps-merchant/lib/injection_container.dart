import 'package:get_it/get_it.dart';
import 'package:dio/dio.dart';
import 'core/network/dio_client.dart';
import 'core/storage/secure_storage.dart';

final sl = GetIt.instance;

Future<void> initDependencies() async {
  sl.registerLazySingleton<SecureStorage>(() => SecureStorage());
  sl.registerLazySingleton<DioClient>(() => DioClient(sl<SecureStorage>()));
  sl.registerLazySingleton<Dio>(() => sl<DioClient>().dio);
}
