package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"regulasi-indexer/handler"
	"regulasi-indexer/repository"
)

func main() {
	ctx := context.Background()

	// ===== MongoDB Connection =====
	mongoClient, err := repository.ConnectMongo("mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}

	// ping
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping error:", err)
	}

	log.Println("✅ MongoDB connected")
	collection := mongoClient.Database("corpus").Collection("jdih")

	// ===== OpenSearch / Elasticsearch =====
	os, err := repository.NewOpenSearch("http://localhost:9200")
	if err != nil {
		log.Fatal("Elasticsearch error:", err)
	}

	log.Println("✅ OpenSearch connected")

	// ===== Handler =====
	h := handler.NewHandler(collection, os)

	// ===== Router =====
	mux := http.NewServeMux()

	mux.HandleFunc("/api/regulasi/search", h.SearchRegulasiHandler)

	// ===== Server =====
	port := getEnv("PORT", "8080")

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	log.Println("🚀 Server running on port", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
