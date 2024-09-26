package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func New(uri, dbName, collectionName string) (Repository, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return Repository{}, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	return Repository{client: client, collection: collection}, nil
}

func (r Repository) CountryNCityByIP(ctx context.Context, ip string) (country, city string, err error) {
	var result struct {
		Country string `bson:"country"`
		City    string `bson:"city"`
	}

	filter := bson.M{"ip": ip}
	err = r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", "", fmt.Errorf("no document found for IP: %s", ip)
		}
		return "", "", fmt.Errorf("failed to find document: %w", err)
	}

	return result.Country, result.City, nil
}
