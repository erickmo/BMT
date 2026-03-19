import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import '../../domain/entities/saldo_info_entity.dart';

part 'kiosk_event.dart';
part 'kiosk_state.dart';

class KioskBloc extends Bloc<KioskEvent, KioskState> {
  KioskBloc() : super(KioskIdle()) {
    on<NfcTagDetected>(_onNfcTagDetected);
    on<KioskReset>(_onReset);
  }

  Future<void> _onNfcTagDetected(
    NfcTagDetected event,
    Emitter<KioskState> emit,
  ) async {
    emit(KioskLoading());
    try {
      // Simulate API call — replace with actual DioClient call
      await Future.delayed(const Duration(milliseconds: 800));

      // Mock response data
      final saldoInfo = SaldoInfoEntity(
        namaNasabah: 'Muhammad Faqih',
        nomorRekening: 'ANNUR-KDR-SU-00000001',
        saldo: 250000,
        transaksiTerakhir: [
          TransaksiSingkat(
            keterangan: 'Pembelian Kantin',
            nominal: 15000,
            isKredit: false,
            tanggal: DateTime.now().subtract(const Duration(hours: 2)),
          ),
          TransaksiSingkat(
            keterangan: 'Top Up Wali',
            nominal: 100000,
            isKredit: true,
            tanggal: DateTime.now().subtract(const Duration(days: 1)),
          ),
          TransaksiSingkat(
            keterangan: 'Pembelian Buku',
            nominal: 35000,
            isKredit: false,
            tanggal: DateTime.now().subtract(const Duration(days: 1, hours: 3)),
          ),
          TransaksiSingkat(
            keterangan: 'Pembelian Laundry',
            nominal: 20000,
            isKredit: false,
            tanggal: DateTime.now().subtract(const Duration(days: 2)),
          ),
          TransaksiSingkat(
            keterangan: 'Top Up Wali',
            nominal: 200000,
            isKredit: true,
            tanggal: DateTime.now().subtract(const Duration(days: 3)),
          ),
        ],
      );
      emit(KioskShowSaldo(saldoInfo));
    } catch (e) {
      emit(KioskError('Kartu tidak dikenali atau terjadi kesalahan'));
    }
  }

  Future<void> _onReset(
    KioskReset event,
    Emitter<KioskState> emit,
  ) async {
    emit(KioskIdle());
  }
}
