package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regulasi-indexer/repository"
)

func main() {
	// index := flag.Int("index", 0, "do indexing?")
	// flag.Parse()

	client, err := repository.ConnectMongo("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	db := client.Database("corpus")
	collection := db.Collection("jdih")

	fmt.Println("✅ DB:", db.Name())
	fmt.Println("✅ Collection:", collection.Name())

	osClient, err := repository.NewOpenSearch("http://localhost:9200")
	if err != nil {
		panic(err)
	}

	err = repository.Ping(osClient)
	if err != nil {
		log.Fatal("OpenSearch not reachable:", err)
	}

	// indexing process
	err = repository.ParseMongoToOpenSearch(collection, osClient)
	if err != nil {
		log.Fatal(err)
	}

}

func PrintJSON(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(string(bytes))
}
