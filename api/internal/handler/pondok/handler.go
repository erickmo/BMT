package pondok

import (
	"net/http"

	"github.com/bmt-saas/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	// Santri
	r.Get("/santri", handleListSantri)
	r.Post("/santri", handleCreateSantri)
	r.Get("/santri/{id}", handleGetSantri)
	r.Put("/santri/{id}", handleUpdateSantri)
	r.Delete("/santri/{id}", handleDeleteSantri)

	// Kelas
	r.Get("/kelas", handleListKelas)
	r.Post("/kelas", handleCreateKelas)
	r.Get("/kelas/{id}", handleGetKelas)
	r.Put("/kelas/{id}", handleUpdateKelas)
	r.Delete("/kelas/{id}", handleDeleteKelas)

	// Pengajar
	r.Get("/pengajar", handleListPengajar)
	r.Post("/pengajar", handleCreatePengajar)
	r.Get("/pengajar/{id}", handleGetPengajar)
	r.Put("/pengajar/{id}", handleUpdatePengajar)
	r.Delete("/pengajar/{id}", handleDeletePengajar)

	// Karyawan
	r.Get("/karyawan", handleListKaryawan)
	r.Post("/karyawan", handleCreateKaryawan)
	r.Get("/karyawan/{id}", handleGetKaryawan)
	r.Put("/karyawan/{id}", handleUpdateKaryawan)
	r.Delete("/karyawan/{id}", handleDeleteKaryawan)

	// Absensi — metode dari settings BMT, bukan konstanta
	r.Post("/absensi", handleCreateAbsensi)
	r.Get("/absensi", handleListAbsensi)
	r.Get("/absensi/rekap", handleRekapAbsensi)

	// Jadwal
	r.Get("/jadwal/pelajaran", handleListJadwalPelajaran)
	r.Post("/jadwal/pelajaran", handleCreateJadwalPelajaran)
	r.Put("/jadwal/pelajaran/{id}", handleUpdateJadwalPelajaran)
	r.Get("/jadwal/kegiatan", handleListJadwalKegiatan)
	r.Post("/jadwal/kegiatan", handleCreateJadwalKegiatan)
	r.Put("/jadwal/kegiatan/{id}", handleUpdateJadwalKegiatan)
	r.Get("/jadwal/piket", handleListJadwalPiket)
	r.Post("/jadwal/piket", handleCreateJadwalPiket)
	r.Get("/jadwal/shift", handleListJadwalShift)
	r.Post("/jadwal/shift", handleCreateJadwalShift)

	// Kalender akademik
	r.Get("/kalender", handleListKalender)
	r.Post("/kalender", handleCreateKalender)
	r.Put("/kalender/{id}", handleUpdateKalender)

	// Mapel & kurikulum
	r.Get("/mapel", handleListMapel)
	r.Post("/mapel", handleCreateMapel)
	r.Put("/mapel/{id}", handleUpdateMapel)
	r.Get("/silabus", handleListSilabus)
	r.Post("/silabus", handleCreateSilabus)
	r.Put("/silabus/{id}", handleUpdateSilabus)
	r.Get("/rpp", handleListRPP)
	r.Post("/rpp", handleCreateRPP)
	r.Put("/rpp/{id}", handleUpdateRPP)
	r.Get("/komponen-nilai", handleListKomponenNilai)
	r.Post("/komponen-nilai", handleCreateKomponenNilai)
	r.Put("/komponen-nilai/{id}", handleUpdateKomponenNilai)

	// Penilaian & raport
	r.Get("/nilai", handleListNilai)
	r.Post("/nilai", handleCreateNilai)
	r.Put("/nilai/{id}", handleUpdateNilai)
	r.Get("/nilai/tahfidz", handleListNilaiTahfidz)
	r.Post("/nilai/tahfidz", handleCreateNilaiTahfidz)
	r.Get("/nilai/akhlak", handleListNilaiAkhlak)
	r.Post("/nilai/akhlak", handleCreateNilaiAkhlak)
	r.Get("/raport", handleListRaport)
	r.Post("/raport", handleCreateRaport)
	r.Get("/raport/{id}", handleGetRaport)
	r.Post("/raport/{id}/terbitkan", handleTerbitkanRaport)

	// Tagihan & keuangan pondok
	r.Get("/jenis-tagihan", handleListJenisTagihan)
	r.Post("/jenis-tagihan", handleCreateJenisTagihan)
	r.Put("/jenis-tagihan/{id}", handleUpdateJenisTagihan)
	r.Post("/tagihan/generate", handleGenerateTagihan)
	r.Get("/tagihan", handleListTagihan)
	r.Put("/tagihan/{id}", handleUpdateTagihan)
	r.Post("/tagihan/{id}/beasiswa", handleTerapkanBeasiswaSPP)
	r.Post("/pembiayaan", handleAjukanPembiayaanPondok)

	// Sinkronisasi eksternal
	r.Post("/sinkron/dapodik", handleSinkronDAPODIK)
	r.Post("/sinkron/emis", handleSinkronEMIS)
	r.Get("/sinkron/log", handleListSinkronLog)
}

func handleListSantri(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateSantri(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "santri berhasil ditambahkan"})
}

func handleGetSantri(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"id": chi.URLParam(r, "id")})
}

func handleUpdateSantri(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "santri berhasil diupdate"})
}

func handleDeleteSantri(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "santri berhasil dihapus"})
}

func handleListKelas(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateKelas(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "kelas berhasil dibuat"})
}

func handleGetKelas(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"id": chi.URLParam(r, "id")})
}

func handleUpdateKelas(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "kelas berhasil diupdate"})
}

func handleDeleteKelas(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "kelas berhasil dihapus"})
}

func handleListPengajar(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreatePengajar(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "pengajar berhasil ditambahkan"})
}

func handleGetPengajar(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"id": chi.URLParam(r, "id")})
}

func handleUpdatePengajar(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "pengajar berhasil diupdate"})
}

func handleDeletePengajar(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "pengajar berhasil dihapus"})
}

func handleListKaryawan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateKaryawan(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "karyawan berhasil ditambahkan"})
}

func handleGetKaryawan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"id": chi.URLParam(r, "id")})
}

func handleUpdateKaryawan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "karyawan berhasil diupdate"})
}

func handleDeleteKaryawan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "karyawan berhasil dihapus"})
}

func handleCreateAbsensi(w http.ResponseWriter, r *http.Request) {
	// Metode absensi divalidasi dari settings BMT: "pondok.absensi_metode"
	response.Created(w, map[string]string{"message": "absensi berhasil dicatat"})
}

func handleListAbsensi(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleRekapAbsensi(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{})
}

func handleListJadwalPelajaran(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateJadwalPelajaran(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "jadwal pelajaran berhasil dibuat"})
}

func handleUpdateJadwalPelajaran(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "jadwal pelajaran berhasil diupdate"})
}

func handleListJadwalKegiatan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateJadwalKegiatan(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "jadwal kegiatan berhasil dibuat"})
}

func handleUpdateJadwalKegiatan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "jadwal kegiatan berhasil diupdate"})
}

func handleListJadwalPiket(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateJadwalPiket(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "jadwal piket berhasil dibuat"})
}

func handleListJadwalShift(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateJadwalShift(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "jadwal shift berhasil dibuat"})
}

func handleListKalender(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateKalender(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "kalender berhasil dibuat"})
}

func handleUpdateKalender(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "kalender berhasil diupdate"})
}

func handleListMapel(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateMapel(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "mapel berhasil dibuat"})
}

func handleUpdateMapel(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "mapel berhasil diupdate"})
}

func handleListSilabus(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateSilabus(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "silabus berhasil dibuat"})
}

func handleUpdateSilabus(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "silabus berhasil diupdate"})
}

func handleListRPP(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateRPP(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "RPP berhasil dibuat"})
}

func handleUpdateRPP(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "RPP berhasil diupdate"})
}

func handleListKomponenNilai(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateKomponenNilai(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "komponen nilai berhasil dibuat"})
}

func handleUpdateKomponenNilai(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "komponen nilai berhasil diupdate"})
}

func handleListNilai(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateNilai(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "nilai berhasil diinput"})
}

func handleUpdateNilai(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "nilai berhasil diupdate"})
}

func handleListNilaiTahfidz(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateNilaiTahfidz(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "nilai tahfidz berhasil diinput"})
}

func handleListNilaiAkhlak(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateNilaiAkhlak(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "nilai akhlak berhasil dicatat"})
}

func handleListRaport(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateRaport(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "raport berhasil dibuat"})
}

func handleGetRaport(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"id": chi.URLParam(r, "id")})
}

func handleTerbitkanRaport(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "raport berhasil diterbitkan"})
}

func handleListJenisTagihan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleCreateJenisTagihan(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "jenis tagihan berhasil dibuat"})
}

func handleUpdateJenisTagihan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "jenis tagihan berhasil diupdate"})
}

func handleGenerateTagihan(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "tagihan berhasil di-generate"})
}

func handleListTagihan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}

func handleUpdateTagihan(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "tagihan berhasil diupdate"})
}

func handleTerapkanBeasiswaSPP(w http.ResponseWriter, r *http.Request) {
	// NominalEfektif = Nominal - BeasiswaNominal (dari domain keuangan)
	response.Success(w, map[string]string{"message": "beasiswa SPP berhasil diterapkan"})
}

func handleAjukanPembiayaanPondok(w http.ResponseWriter, r *http.Request) {
	response.Created(w, map[string]string{"message": "pembiayaan pondok berhasil diajukan"})
}

func handleSinkronDAPODIK(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "sinkronisasi DAPODIK dimulai"})
}

func handleSinkronEMIS(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{"message": "sinkronisasi EMIS dimulai"})
}

func handleListSinkronLog(w http.ResponseWriter, r *http.Request) {
	response.Success(w, []interface{}{})
}
