import '../../domain/entities/pecahan_entity.dart';
import '../../domain/entities/sesi_entity.dart';

class SesiModel {
  final String id;
  final String tellerId;
  final String cabangId;
  final int modalAwal;
  final int modalAkhir;
  final String status;
  final String dibukaPada;
  final String? ditutupPada;
  final List<PecahanSesiModel> pecahan;

  const SesiModel({
    required this.id,
    required this.tellerId,
    required this.cabangId,
    required this.modalAwal,
    required this.modalAkhir,
    required this.status,
    required this.dibukaPada,
    this.ditutupPada,
    required this.pecahan,
  });

  factory SesiModel.fromJson(Map<String, dynamic> json) {
    final rawPecahan = json['pecahan'] as List<dynamic>? ?? [];
    return SesiModel(
      id: json['id'] as String,
      tellerId: json['teller_id'] as String,
      cabangId: json['cabang_id'] as String,
      modalAwal: (json['modal_awal'] as num).toInt(),
      modalAkhir: (json['modal_akhir'] as num?)?.toInt() ?? 0,
      status: json['status'] as String,
      dibukaPada: json['dibuka_pada'] as String,
      ditutupPada: json['ditutup_pada'] as String?,
      pecahan: rawPecahan
          .map((e) => PecahanSesiModel.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }

  SesiEntity toEntity() {
    return SesiEntity(
      id: id,
      tellerId: tellerId,
      cabangId: cabangId,
      modalAwal: modalAwal,
      modalAkhir: modalAkhir,
      status: status,
      dibukaPada: DateTime.parse(dibukaPada).toLocal(),
      ditutupPada:
          ditutupPada != null ? DateTime.parse(ditutupPada!).toLocal() : null,
      pecahan: pecahan.map((p) => p.toEntity()).toList(),
    );
  }
}

class PecahanSesiModel {
  final String pecahanId;
  final int nominal;
  final String label;
  final String jenis;
  final int jumlah;

  const PecahanSesiModel({
    required this.pecahanId,
    required this.nominal,
    required this.label,
    required this.jenis,
    required this.jumlah,
  });

  factory PecahanSesiModel.fromJson(Map<String, dynamic> json) {
    return PecahanSesiModel(
      pecahanId: json['pecahan_id'] as String,
      nominal: (json['nominal'] as num).toInt(),
      label: json['label'] as String? ?? '',
      jenis: json['jenis'] as String? ?? 'KERTAS',
      jumlah: (json['jumlah'] as num).toInt(),
    );
  }

  PecahanSesiEntity toEntity() {
    return PecahanSesiEntity(
      pecahanId: pecahanId,
      nominal: nominal,
      label: label,
      jenis: jenis,
      jumlah: jumlah,
    );
  }
}

class PecahanModel {
  final String id;
  final int nominal;
  final String jenis;
  final String label;
  final bool isAktif;
  final int urutan;

  const PecahanModel({
    required this.id,
    required this.nominal,
    required this.jenis,
    required this.label,
    required this.isAktif,
    required this.urutan,
  });

  factory PecahanModel.fromJson(Map<String, dynamic> json) {
    return PecahanModel(
      id: json['id'] as String,
      nominal: (json['nominal'] as num).toInt(),
      jenis: json['jenis'] as String,
      label: json['label'] as String,
      isAktif: json['is_aktif'] as bool? ?? true,
      urutan: (json['urutan'] as num?)?.toInt() ?? 0,
    );
  }

  PecahanEntity toEntity() {
    return PecahanEntity(
      id: id,
      nominal: nominal,
      jenis: jenis,
      label: label,
      isAktif: isAktif,
      urutan: urutan,
    );
  }
}
