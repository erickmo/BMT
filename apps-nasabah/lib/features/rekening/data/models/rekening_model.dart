import '../../domain/entities/rekening_entity.dart';
import '../../domain/entities/transaksi_entity.dart';

class RekeningModel {
  final String id;
  final String nomorRekening;
  final String jenisRekeningNama;
  final String jenisRekeningKode;
  final int saldo;
  final String status;
  final String? alasanBlokir;
  final int biayaAdminBulanan;
  final String tanggalBuka;
  final String? tanggalTutup;

  const RekeningModel({
    required this.id,
    required this.nomorRekening,
    required this.jenisRekeningNama,
    required this.jenisRekeningKode,
    required this.saldo,
    required this.status,
    this.alasanBlokir,
    required this.biayaAdminBulanan,
    required this.tanggalBuka,
    this.tanggalTutup,
  });

  factory RekeningModel.fromJson(Map<String, dynamic> json) {
    return RekeningModel(
      id: json['id'] as String,
      nomorRekening: json['nomor_rekening'] as String,
      jenisRekeningNama: json['jenis_rekening_nama'] as String? ?? '',
      jenisRekeningKode: json['jenis_rekening_kode'] as String? ?? '',
      saldo: (json['saldo'] as num).toInt(),
      status: json['status'] as String,
      alasanBlokir: json['alasan_blokir'] as String?,
      biayaAdminBulanan: (json['biaya_admin_bulanan'] as num?)?.toInt() ?? 0,
      tanggalBuka: json['tanggal_buka'] as String,
      tanggalTutup: json['tanggal_tutup'] as String?,
    );
  }

  RekeningEntity toEntity() {
    return RekeningEntity(
      id: id,
      nomorRekening: nomorRekening,
      jenisRekeningNama: jenisRekeningNama,
      jenisRekeningKode: jenisRekeningKode,
      saldo: saldo,
      status: status,
      alasanBlokir: alasanBlokir,
      biayaAdminBulanan: biayaAdminBulanan,
      tanggalBuka: DateTime.parse(tanggalBuka),
      tanggalTutup:
          tanggalTutup != null ? DateTime.parse(tanggalTutup!) : null,
    );
  }
}

class TransaksiModel {
  final String id;
  final String rekeningId;
  final String jenis;
  final String tipe;
  final int nominal;
  final int saldoAkhir;
  final String keterangan;
  final String tanggal;
  final String? referensi;

  const TransaksiModel({
    required this.id,
    required this.rekeningId,
    required this.jenis,
    required this.tipe,
    required this.nominal,
    required this.saldoAkhir,
    required this.keterangan,
    required this.tanggal,
    this.referensi,
  });

  factory TransaksiModel.fromJson(Map<String, dynamic> json) {
    return TransaksiModel(
      id: json['id'] as String,
      rekeningId: json['rekening_id'] as String,
      jenis: json['jenis'] as String,
      tipe: json['tipe'] as String? ?? '',
      nominal: (json['nominal'] as num).toInt(),
      saldoAkhir: (json['saldo_akhir'] as num).toInt(),
      keterangan: json['keterangan'] as String? ?? '',
      tanggal: json['tanggal'] as String,
      referensi: json['referensi'] as String?,
    );
  }

  TransaksiEntity toEntity() {
    return TransaksiEntity(
      id: id,
      rekeningId: rekeningId,
      jenis: jenis,
      tipe: tipe,
      nominal: nominal,
      saldoAkhir: saldoAkhir,
      keterangan: keterangan,
      tanggal: DateTime.parse(tanggal).toLocal(),
      referensi: referensi,
    );
  }
}
