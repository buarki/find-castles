package db

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/buarki/find-castles/castle"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveCastles(ctx context.Context, collection *mongo.Collection, castles []castle.Model) error {
	var operations []mongo.WriteModel

	for _, c := range castles {
		filter := bson.M{
			"country": strings.ToLower(c.Country.String()),
			"name":    strings.ToLower(c.Name),
		}
		update := bson.M{
			"$set": bson.M{
				"name":             strings.ToLower(c.Name),
				"link":             c.Link,
				"country":          strings.ToLower(c.Country.String()),
				"state":            c.State,
				"city":             c.City,
				"district":         c.District,
				"foundationPeriod": c.FoundationPeriod,
			},
		}
		operation := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		operations = append(operations, operation)
	}

	if len(operations) > 0 {
		_, err := collection.BulkWrite(ctx, operations)
		if err != nil {
			log.Printf("failed to upsert castles: %v", err)
			return fmt.Errorf("failed to upsert [%d] castles, got %v", len(operations), err)
		} else {
			log.Printf("successfully upserted [%d] castles", len(castles))
		}
	}
	return nil
}
