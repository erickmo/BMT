import 'package:get_it/get_it.dart';
import 'core/network/dio_client.dart';

final sl = GetIt.instance;

Future<void> initDependencies() async {
  sl.registerLazySingleton<DioClient>(() => DioClient());
}
