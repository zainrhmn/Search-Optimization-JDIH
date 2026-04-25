package handler

import (
	"context"
	"log"
	"net/http"
	"regulasi-indexer/common"
	"regulasi-indexer/model"
	"regulasi-indexer/repository"
	"strconv"

	"github.com/opensearch-project/opensearch-go/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handler struct {
	Collection *mongo.Collection
	OS         *opensearch.Client
}

type SearchRegulasiHandlerResponse struct {
	MongoID string `json:"mongo_id"`
	Title   string `json:"title"`
	Babs    []Bab  `json:"babs"`
}

type Ayat struct {
	Nomor     string `json:"nomor"`
	Isi       string `json:"isi"`
	Highlight bool   `json:"highlight"`
}

type Pasal struct {
	Nomor string `json:"nomor"`
	Ayats []Ayat `json:"ayats"`
}

type Bab struct {
	Title  string  `json:"title"`
	Number string  `json:"number"`
	Pasals []Pasal `json:"pasals"`
}

func NewHandler(c *mongo.Collection, os *opensearch.Client) *Handler {
	return &Handler{
		Collection: c,
		OS:         os,
	}
}

func (h *Handler) SearchRegulasiHandler(w http.ResponseWriter, r *http.Request) {
	_ = context.Background()

	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		http.Error(w, "keyword required", http.StatusBadRequest)
		return
	}

	// 🔥 Call your existing pipeline
	searchResults, err := repository.SearchRegulasi(h.OS, keyword)
	if err != nil {
		log.Println("got an error at SearchRegulasi: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// no result → return empty
	if len(searchResults) == 0 {
		common.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"query": keyword,
			"total": 0,
			"data":  []interface{}{},
		})
		return
	}

	ids := repository.ExtractMongoIDs(searchResults)
	mongoDocs, err := repository.GetByIDs(h.Collection, ids)
	if err != nil {
		log.Println("got an error at repository.GetByIDs: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := buildResponse(searchResults, mongoDocs)
	common.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"query": keyword,
		"total": len(response),
		"data":  response,
	})
}

func buildResponse(searchResults []model.SearchResult, mongoDocs []model.RegulasiMongo) []SearchRegulasiHandlerResponse {
	// ===== 1. Build match map =====
	matchMap := make(map[string]map[string]map[string]bool)

	for _, item := range searchResults {
		if item.MongoID == "" || item.Pasal == "" || item.Ayat == "" {
			continue
		}

		if _, ok := matchMap[item.MongoID]; !ok {
			matchMap[item.MongoID] = make(map[string]map[string]bool)
		}
		if _, ok := matchMap[item.MongoID][item.Pasal]; !ok {
			matchMap[item.MongoID][item.Pasal] = make(map[string]bool)
		}

		matchMap[item.MongoID][item.Pasal][item.Ayat] = true
	}

	// ===== 2. Map mongo docs by ID =====
	docMap := make(map[string]model.RegulasiMongo)

	for _, doc := range mongoDocs {
		docMap[doc.ID.Hex()] = doc
	}

	// ===== 3. Preserve document order from ranking =====
	var orderedDocIDs []string
	seen := make(map[string]bool)

	for _, s := range searchResults {
		if s.MongoID == "" {
			continue
		}
		if !seen[s.MongoID] {
			seen[s.MongoID] = true
			orderedDocIDs = append(orderedDocIDs, s.MongoID)
		}
	}

	// ===== 4. Build final response =====
	var results []SearchRegulasiHandlerResponse

	for _, mongoID := range orderedDocIDs {
		doc, ok := docMap[mongoID] // doc is model.RegulasiMongo
		if !ok {
			continue
		}

		docMatch := matchMap[mongoID]

		resp := SearchRegulasiHandlerResponse{
			MongoID: mongoID,
			Title:   doc.Title,
		}

		for _, b := range doc.Babs {
			babResponse := Bab{
				Title:  b.Title,
				Number: strconv.Itoa(b.Number),
			}

			for _, p := range b.Pasals {
				pasalNumber := strconv.Itoa(p.Number)

				pasal := Pasal{
					Nomor: pasalNumber,
				}

				hasMatch := false

				for i, ayatText := range p.Ayats {
					ayatNomor := strconv.Itoa(i + 1)
					highlight := false
					if docMatch != nil {
						if pasalMatch, ok := docMatch[pasalNumber]; ok {
							if pasalMatch[ayatNomor] {
								highlight = true
								hasMatch = true
							}
						}
					}

					pasal.Ayats = append(pasal.Ayats, Ayat{
						Nomor:     ayatNomor,
						Isi:       ayatText,
						Highlight: highlight,
					})
				}

				// only include matched pasal
				if hasMatch {
					babResponse.Pasals = append(babResponse.Pasals, pasal)
				}
			}

			// only include bab with results
			if len(babResponse.Pasals) > 0 {
				resp.Babs = append(resp.Babs, babResponse)
			}
		}

		results = append(results, resp)
	}

	return results
}
