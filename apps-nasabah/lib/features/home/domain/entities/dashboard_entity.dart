import 'package:equatable/equatable.dart';
import '../../../rekening/domain/entities/rekening_entity.dart';

class DashboardEntity extends Equatable {
  final String namaNasabah;
  final String nomorNasabah;
  final int totalSaldo;
  final List<RekeningEntity> rekening;

  const DashboardEntity({
    required this.namaNasabah,
    required this.nomorNasabah,
    required this.totalSaldo,
    required this.rekening,
  });

  @override
  List<Object?> get props => [namaNasabah, nomorNasabah, totalSaldo, rekening];
}
