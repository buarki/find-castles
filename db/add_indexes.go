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
	isFalse := false

	nameAndCountry := mongo.IndexModel{
		Keys: bson.D{
			{Key: "country", Value: 1},
			{Key: "name", Value: 1},
		},
		Options: &options.IndexOptions{
			Unique: &isTrue,
		},
	}

	matchingTags := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "matchingTags",
				Value: 1,
			},
		},
	}

	countryIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "country", Value: 1},
		},
		Options: &options.IndexOptions{
			Unique: &isFalse,
		},
	}

	webNameIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "webName", Value: 1},
		},
		Options: &options.IndexOptions{
			Unique: &isFalse,
		},
	}

	indexes := []mongo.IndexModel{
		nameAndCountry,
		matchingTags,
		countryIndex,
		webNameIndex,
	}

	name, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}
	slog.Info("index upserted", "message", fmt.Sprintf("index %s created", name))
	return nil
}
