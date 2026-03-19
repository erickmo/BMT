import 'package:equatable/equatable.dart';

abstract class Failure extends Equatable {
  final String message;
  const Failure(this.message);
  @override
  List<Object> get props => [message];
}

class ServerFailure extends Failure {
  final int? statusCode;
  const ServerFailure(super.message, {this.statusCode});
  @override
  List<Object> get props => [message, statusCode ?? 0];
}

class NetworkFailure extends Failure {
  const NetworkFailure([super.message = 'Tidak ada koneksi internet']);
}

class UnauthorizedFailure extends Failure {
  const UnauthorizedFailure([super.message = 'Sesi habis, silakan login kembali']);
}

class NfcFailure extends Failure {
  const NfcFailure([super.message = 'Gagal membaca kartu NFC']);
}

class PinSalahFailure extends Failure {
  const PinSalahFailure([super.message = 'PIN kartu tidak sesuai']);
}

class SaldoKurangFailure extends Failure {
  const SaldoKurangFailure([super.message = 'Saldo tidak mencukupi']);
}

class UnexpectedFailure extends Failure {
  const UnexpectedFailure([super.message = 'Terjadi kesalahan yang tidak terduga']);
}
