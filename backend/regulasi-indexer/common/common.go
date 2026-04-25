package common

import (
	"encoding/json"
	"net/http"
)

func UniqueStrings(ids []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, id := range ids {
		if _, exists := seen[id]; !exists {
			seen[id] = struct{}{}
			result = append(result, id)
		}
	}

	return result
}

func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func GetString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
