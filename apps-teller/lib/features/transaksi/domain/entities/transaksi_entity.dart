import 'package:equatable/equatable.dart';

class TransaksiResultEntity extends Equatable {
  final String id;
  final String rekeningId;
  final String nomorRekening;
  final String namaNasabah;
  final String jenis; // DEBIT | KREDIT
  final String tipe;
  final int nominal;
  final int saldoSebelum;
  final int saldoAkhir;
  final String keterangan;
  final DateTime tanggal;

  const TransaksiResultEntity({
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

  @override
  List<Object> get props => [
        id,
        rekeningId,
        nomorRekening,
        namaNasabah,
        jenis,
        tipe,
        nominal,
        saldoSebelum,
        saldoAkhir,
        keterangan,
        tanggal,
      ];
}

class NasabahSearchResult extends Equatable {
  final String id;
  final String nomorNasabah;
  final String nama;
  final List<RekeningSearchResult> rekening;

  const NasabahSearchResult({
    required this.id,
    required this.nomorNasabah,
    required this.nama,
    required this.rekening,
  });

  @override
  List<Object> get props => [id, nomorNasabah, nama, rekening];
}

class RekeningSearchResult extends Equatable {
  final String id;
  final String nomorRekening;
  final String jenisNama;
  final int saldo;
  final String status;

  const RekeningSearchResult({
    required this.id,
    required this.nomorRekening,
    required this.jenisNama,
    required this.saldo,
    required this.status,
  });

  bool get isAktif => status == 'AKTIF';

  @override
  List<Object> get props => [id, nomorRekening, jenisNama, saldo, status];
}
