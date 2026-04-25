package model

type RegulasiIndex struct {
	ID            string `json:"id"`
	MongoID       string `json:"mongo_id"`
	JudulRegulasi string `json:"judul_regulasi"`
	NomorRegulasi string `json:"nomor_regulasi"`
	Tahun         string `json:"tahun"`
	JenisRegulasi string `json:"jenis_regulasi"`

	Bab      string `json:"bab"`
	JudulBab string `json:"judul_bab"`
	Pasal    string `json:"pasal"`
	Bagian   string `json:"bagian"`
	Ayat     string `json:"ayat"`

	Isi      string `json:"isi"`
	FullText string `json:"full_text"`
}

type SearchResult struct {
	MongoID string `json:"mongo_id"`
	Pasal   string `json:"pasal"`
	Ayat    string `json:"ayat"`
}
