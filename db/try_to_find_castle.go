package db

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/buarki/find-castles/castle"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCastleNotFound = errors.New("castle not found")
)

func TryToFindCastle(ctx context.Context, collection *mongo.Collection, c castle.Model) (castle.Model, error) {
	var result castle.Model

	matchingTags := c.GetMatchingTags()
	c.GetMatchingTags()

	regexPattern := fmt.Sprintf(".*%s.*", regexp.QuoteMeta(c.FilteredName()))
	filter := bson.M{
		"country":      c.Country.String(),
		"name":         bson.M{"$regex": regexPattern, "$options": "i"},
		"matchingTags": bson.M{"$in": matchingTags},
	}
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return castle.Model{}, ErrCastleNotFound
		}
		return castle.Model{}, err
	}

	return result, nil
}

func TryToFindCastles(ctx context.Context, collection *mongo.Collection, castles []castle.Model) ([]castle.Model, error) {
	var results []castle.Model
	var filters []bson.M

	for _, c := range castles {
		matchingTags := c.GetMatchingTags()

		regexPattern := fmt.Sprintf(".*%s.*", regexp.QuoteMeta(c.FilteredName()))
		filter := bson.M{
			"country": c.Country.String(),
			"$or": []bson.M{
				{"name": bson.M{"$regex": regexPattern, "$options": "i"}},
				{"name": bson.M{"$in": matchingTags}},
				{"name": bson.M{"$in": strings.Split(strings.ToLower(c.Name), " ")}},
			},
		}
		filters = append(filters, filter)
	}

	query := bson.M{
		"$or": filters,
	}

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find query: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result castle.Model
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode castle: %w", err)
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return results, nil
}
