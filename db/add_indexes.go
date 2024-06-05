package db

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddIndexes(ctx context.Context, collection *mongo.Collection) error {
	isTrue := true
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "country", Value: 1},
			{Key: "name", Value: 1},
		},
		Options: &options.IndexOptions{
			Unique: &isTrue,
		},
	}
	name, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}
	slog.Info("index upserted", "message", fmt.Sprintf("index %s created", name))
	return nil
}
