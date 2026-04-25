package model

import "go.mongodb.org/mongo-driver/v2/bson"

type RegulasiMongo struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	Nomor         int           `bson:"nomor"`
	Tahun         int           `bson:"tahun"`
	JenisRegulasi string        `bson:"jenis_regulasi"`
	Title         string        `bson:"title"`
	Status        string        `bson:"status"`
	CreatedAt     string        `bson:"created_at"`
	UpdatedAt     string        `bson:"updated_at"`
	Babs          []Bab         `bson:"babs"`
}

type Bab struct {
	Number int     `bson:"number"`
	Title  string  `bson:"title"`
	Pasals []Pasal `bson:"pasals"`
}

type Pasal struct {
	Number      int      `bson:"number"`
	Title       string   `bson:"title"`
	Description string   `bson:"description"`
	Ayats       []string `bson:"ayats"`
	Part        *Part    `bson:"part,omitempty"`
}

type Part struct {
	Name  string `bson:"name"`
	Title string `bson:"title"`
}
