package repository

import (
	"context"
	"fmt"
	"log"
	"regulasi-indexer/common"
	"regulasi-indexer/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func GetByID(collection *mongo.Collection, id string) {
	// convert string to ObjectID
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	// filter
	filter := bson.M{"_id": objectID}

	var result bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found")
			return
		}
		log.Fatal(err)
	}

	fmt.Println("Result:", result)
}

func GetByIDs(collection *mongo.Collection, ids []string) ([]model.RegulasiMongo, error) {
	var objectIDs []bson.ObjectID

	// remove duplicates first
	ids = common.UniqueStrings(ids)

	// convert string IDs to ObjectID
	for _, id := range ids {
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objID)
	}

	// filter using $in
	filter := bson.M{
		"_id": bson.M{
			"$in": objectIDs,
		},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []model.RegulasiMongo

	// 🔥 decode directly into struct
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
