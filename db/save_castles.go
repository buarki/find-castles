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
		update := bson.M{
			"$set": prepareObjectToSave(c),
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

func prepareObjectToSave(c castle.Model) bson.M {
	object := bson.M{
		"name":         strings.ToLower(c.FilteredName()),
		"sources":      c.Sources,
		"country":      strings.ToLower(c.Country.String()),
		"matchingTags": c.GetMatchingTags(),
		"pictureURL":   c.PictureURL,
	}
	if c.State != "" {
		object["state"] = c.State
	}
	if c.City != "" {
		object["city"] = c.City
	}
	if c.District != "" {
		object["district"] = c.District
	}
	if c.FoundationPeriod != "" {
		object["foundationPeriod"] = c.FoundationPeriod
	}
	if c.PropertyCondition != "" {
		object["propertyCondition"] = c.PropertyCondition
	}
	if c.Coordinates != "" {
		object["coordinates"] = c.Coordinates
	}
	if c.Contact != nil {
		object["contact"] = bson.M{
			"phone": c.Contact.Phone,
			"email": c.Contact.Email,
		}
	}
	return object
}
