package pengiriman

import (
	"encoding/json"
	"errors"
)

var (
	ErrAlamatTidakLengkap = errors.New("alamat pengiriman tidak lengkap")
)

// Alamat represents a shipping address embedded in pesanan.alamat_kirim (JSON)
type Alamat struct {
	NamaPenerima string `json:"nama_penerima"`
	Telepon      string `json:"telepon"`
	Provinsi     string `json:"provinsi"`
	Kabupaten    string `json:"kabupaten"`
	Kecamatan    string `json:"kecamatan"`
	Kelurahan    string `json:"kelurahan"`
	JalanDetail  string `json:"jalan_detail"`
	KodePos      string `json:"kode_pos,omitempty"`
	Catatan      string `json:"catatan,omitempty"`
}

// InfoPengiriman holds delivery tracking information
type InfoPengiriman struct {
	Kurir      string `json:"kurir"`
	NomorResi  string `json:"nomor_resi"`
	StatusInfo string `json:"status_info,omitempty"`
}

func (a *Alamat) Validasi() error {
	if a.NamaPenerima == "" {
		return errors.New("nama penerima wajib diisi")
	}
	if a.Telepon == "" {
		return errors.New("telepon penerima wajib diisi")
	}
	if a.Provinsi == "" || a.Kabupaten == "" || a.JalanDetail == "" {
		return ErrAlamatTidakLengkap
	}
	return nil
}

func (a *Alamat) ToJSON() (json.RawMessage, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func AlamatFromJSON(raw json.RawMessage) (*Alamat, error) {
	var a Alamat
	if err := json.Unmarshal(raw, &a); err != nil {
		return nil, errors.New("gagal membaca data alamat pengiriman")
	}
	return &a, nil
}
