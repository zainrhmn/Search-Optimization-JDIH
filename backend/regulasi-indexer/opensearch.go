package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/opensearch-project/opensearch-go/v2"

	"regulasi-indexer/model"
)

var indexName = "regulasi_index"

func NewOpenSearch() (*opensearch.Client, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	})

	if err != nil {
		return nil, err
	}

	return client, nil
}

func Ping(client *opensearch.Client) error {
	res, err := client.Info()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fmt.Println("OpenSearch connected 🚀")
	return nil
}

func SearchRegulasi(client *opensearch.Client, keyword string) ([]string, error) {
	result := []string{}
	query := map[string]interface{}{
		"size": 20,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query": keyword,
							"type":  "best_fields",
							"fields": []string{
								"nomor_regulasi^6",
								"tahun^6",
								"judul_regulasi^5",
								"judul_bab^4",
								"pasal^3",
								"bagian^2",
								"isi^1",
								"full_text^0.5",
							},
							"minimum_should_match": "60%",
						},
					},
					map[string]interface{}{
						"match_phrase": map[string]interface{}{
							"judul_regulasi": map[string]interface{}{
								"query": keyword,
								"boost": 10,
							},
						},
					},
					map[string]interface{}{
						"match_phrase": map[string]interface{}{
							"full_text": map[string]interface{}{
								"query": keyword,
								"boost": 4,
							},
						},
					},
				},
			},
		},
	}
	// convert ke JSON
	body, err := json.Marshal(query)
	if err != nil {
		return result, err
	}

	// execute search
	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(indexName),
		client.Search.WithBody(bytes.NewReader(body)),
		client.Search.WithPretty(),
	)

	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	var osResult map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&osResult); err != nil {
		return result, err
	}

	hits := osResult["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, h := range hits {
		hit := h.(map[string]interface{})
		source := hit["_source"].(map[string]interface{})

		idStr, ok := source["mongo_id"].(string)
		if !ok {
			continue
		}

		result = append(result, idStr)
	}

	return result, nil
}

func ParseMongoToOpenSearch(coll *mongo.Collection, client *opensearch.Client) error {
	ctx := context.Background()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	batchSize := 500
	var bulkDocs []model.RegulasiIndex
	total := 0

	for cursor.Next(ctx) {
		var doc model.RegulasiMongo

		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Decode error: %v\n", err)
			continue
		}

		// 🔥 transform nested → flat
		transformed := transformMongoToIndex(doc)
		bulkDocs = append(bulkDocs, transformed...)

		// 🔥 bulk insert
		if len(bulkDocs) >= batchSize {
			if err := bulkIndex(client, indexName, bulkDocs); err != nil {
				log.Println("Bulk error:", err)
			} else {
				total += len(bulkDocs)
				fmt.Println("Indexed:", total)
			}
			bulkDocs = nil
		}
	}

	// 🔥 sisa data
	if len(bulkDocs) > 0 {
		if err := bulkIndex(client, indexName, bulkDocs); err != nil {
			return err
		}
		total += len(bulkDocs)
	}

	fmt.Println("✅ Total indexed:", total)
	return nil
}

func transformMongoToIndex(doc model.RegulasiMongo) []model.RegulasiIndex {
	var results []model.RegulasiIndex

	for _, bab := range doc.Babs {
		for _, pasal := range bab.Pasals {

			bagian := ""
			if pasal.Part != nil {
				bagian = pasal.Part.Title
			}

			for i, ayatText := range pasal.Ayats {

				ayatNumber := fmt.Sprintf("(%d)", i+1)

				id := fmt.Sprintf("%d_%d_pasal_%d_ayat_%d",
					doc.Nomor,
					doc.Tahun,
					pasal.Number,
					i+1,
				)

				fullText := fmt.Sprintf("%s %v %s Pasal %d %s %s",
					doc.Title,
					bab.Number,
					bab.Title,
					pasal.Number,
					ayatNumber,
					ayatText,
				)

				results = append(results, model.RegulasiIndex{
					ID:            id,
					MongoID:       doc.ID.Hex(),
					JudulRegulasi: doc.Title,
					NomorRegulasi: fmt.Sprintf("%d", doc.Nomor),
					Tahun:         fmt.Sprintf("%d", doc.Tahun),
					JenisRegulasi: doc.JenisRegulasi,

					Bab:      fmt.Sprintf("%v", bab.Number),
					JudulBab: bab.Title,
					Pasal:    fmt.Sprintf("Pasal %d", pasal.Number),
					Bagian:   bagian,
					Ayat:     ayatNumber,

					Isi:      ayatText,
					FullText: fullText,
				})
			}
		}
	}

	return results
}

func bulkIndex(client *opensearch.Client, index string, docs []model.RegulasiIndex) error {
	var buf bytes.Buffer

	for _, doc := range docs {
		meta := fmt.Sprintf(`{ "index" : { "_index" : "%s", "_id" : "%s" } }`+"\n",
			index, doc.ID)

		data, _ := json.Marshal(doc)

		buf.WriteString(meta)
		buf.Write(data)
		buf.WriteString("\n")
	}

	res, err := client.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
