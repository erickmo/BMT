import '../../domain/entities/transaksi_entity.dart';

class TransaksiResultModel {
  final String id;
  final String rekeningId;
  final String nomorRekening;
  final String namaNasabah;
  final String jenis;
  final String tipe;
  final int nominal;
  final int saldoSebelum;
  final int saldoAkhir;
  final String keterangan;
  final String tanggal;

  const TransaksiResultModel({
    required this.id,
    required this.rekeningId,
    required this.nomorRekening,
    required this.namaNasabah,
    required this.jenis,
    required this.tipe,
    required this.nominal,
    required this.saldoSebelum,
    required this.saldoAkhir,
    required this.keterangan,
    required this.tanggal,
  });

  factory TransaksiResultModel.fromJson(Map<String, dynamic> json) {
    return TransaksiResultModel(
      id: json['id'] as String,
      rekeningId: json['rekening_id'] as String,
      nomorRekening: json['nomor_rekening'] as String? ?? '',
      namaNasabah: json['nama_nasabah'] as String? ?? '',
      jenis: json['jenis'] as String,
      tipe: json['tipe'] as String? ?? '',
      nominal: (json['nominal'] as num).toInt(),
      saldoSebelum: (json['saldo_sebelum'] as num?)?.toInt() ?? 0,
      saldoAkhir: (json['saldo_akhir'] as num).toInt(),
      keterangan: json['keterangan'] as String? ?? '',
      tanggal: json['tanggal'] as String,
    );
  }

  TransaksiResultEntity toEntity() {
    return TransaksiResultEntity(
      id: id,
      rekeningId: rekeningId,
      nomorRekening: nomorRekening,
      namaNasabah: namaNasabah,
      jenis: jenis,
      tipe: tipe,
      nominal: nominal,
      saldoSebelum: saldoSebelum,
      saldoAkhir: saldoAkhir,
      keterangan: keterangan,
      tanggal: DateTime.parse(tanggal).toLocal(),
    );
  }
}

class NasabahSearchModel {
  final String id;
  final String nomorNasabah;
  final String nama;
  final List<RekeningSearchModel> rekening;

  const NasabahSearchModel({
    required this.id,
    required this.nomorNasabah,
    required this.nama,
    required this.rekening,
  });

  factory NasabahSearchModel.fromJson(Map<String, dynamic> json) {
    final rawRek = json['rekening'] as List<dynamic>? ?? [];
    return NasabahSearchModel(
      id: json['id'] as String,
      nomorNasabah: json['nomor_nasabah'] as String,
      nama: json['nama'] as String,
      rekening: rawRek
          .map((e) => RekeningSearchModel.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }

  NasabahSearchResult toEntity() {
    return NasabahSearchResult(
      id: id,
      nomorNasabah: nomorNasabah,
      nama: nama,
      rekening: rekening.map((r) => r.toEntity()).toList(),
    );
  }
}

class RekeningSearchModel {
  final String id;
  final String nomorRekening;
  final String jenisNama;
  final int saldo;
  final String status;

  const RekeningSearchModel({
    required this.id,
    required this.nomorRekening,
    required this.jenisNama,
    required this.saldo,
    required this.status,
  });

  factory RekeningSearchModel.fromJson(Map<String, dynamic> json) {
    return RekeningSearchModel(
      id: json['id'] as String,
      nomorRekening: json['nomor_rekening'] as String,
      jenisNama: json['jenis_rekening_nama'] as String? ?? '',
      saldo: (json['saldo'] as num).toInt(),
      status: json['status'] as String,
    );
  }

  RekeningSearchResult toEntity() {
    return RekeningSearchResult(
      id: id,
      nomorRekening: nomorRekening,
      jenisNama: jenisNama,
      saldo: saldo,
      status: status,
    );
  }
}
