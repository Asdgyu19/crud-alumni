package models

// Result struktur untuk hasil query alumni dengan gaji
type Result struct {
	ID             int     `json:"id"`
	Nama           string  `json:"nama"`
	Jurusan        string  `json:"jurusan"`
	TahunLulus     int     `json:"tahun_lulus"`
	BidangIndustri string  `json:"bidang_industri"`
	NamaPerusahaan string  `json:"nama_perusahaan"`
	PosisiJabatan  string  `json:"posisi_jabatan"`
	RangeGaji      float64 `json:"range_gaji"`
}
