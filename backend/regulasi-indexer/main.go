package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

func main() {
	index := flag.Int("index", 0, "do indexing?")
	flag.Parse()

	client, err := ConnectMongo("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	db := client.Database("corpus")
	collection := db.Collection("jdih")

	fmt.Println("✅ DB:", db.Name())
	fmt.Println("✅ Collection:", collection.Name())

	osClient, err := NewOpenSearch()
	if err != nil {
		panic(err)
	}

	err = Ping(osClient)
	if err != nil {
		log.Fatal("OpenSearch not reachable:", err)
	}

	// indexing process
	if *index == 1 {
		err = ParseMongoToOpenSearch(collection, osClient)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Testing Search
	mongoIDs, err := SearchRegulasi(osClient, "renstra")
	log.Println(mongoIDs)

	// 🔥 sanity check: ambil 1 data
	// var result map[string]interface{}

	// err = collection.FindOne(context.Background(), map[string]interface{}{}).Decode(&result)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("✅ Sample data found")
	// PrintJSON(result)
}

func PrintJSON(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(string(bytes))
}
