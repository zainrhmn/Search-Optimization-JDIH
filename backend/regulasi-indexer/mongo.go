package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongo(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// 🔥 Ping untuk memastikan koneksi
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("✅ Connected to MongoDB (v2)")

	return client, nil
}