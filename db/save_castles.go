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
			"name":    strings.ToLower(c.FilteredName()),
		}
		// TODO collect fields with values only
		update := bson.M{
			"$set": bson.M{
				"name":             strings.ToLower(c.FilteredName()),
				"link":             c.Link,
				"sources":          c.Sources,
				"country":          strings.ToLower(c.Country.String()),
				"state":            strings.ToLower(c.State),
				"city":             strings.ToLower(c.City),
				"district":         strings.ToLower(c.District),
				"foundationPeriod": c.FoundationPeriod,
				// Only if present
				"propertyCondition": c.PropertyCondition.String(),
				"matchingTags":      c.GetMatchingTags(),
				"pictureURL":        c.PictureLink,
				// Only if present
				"coordinates": c.Coordinates,
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
