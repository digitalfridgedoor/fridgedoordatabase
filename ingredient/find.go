package ingredient

import (
	"context"

	"github.com/digitalfridgedoor/fridgedoordatabase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindByName finds ingredients starting with the given letter
func FindByName(ctx context.Context, startsWith string) ([]*Ingredient, error) {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(20)

	regex := bson.M{"$regex": primitive.Regex{Pattern: "\\b" + startsWith, Options: "i"}}
	startsWithBson := bson.M{"name": regex}

	cur, err := mongoCollection.Find(ctx, startsWithBson, findOptions)
	if err != nil {
		return make([]*Ingredient, 0), err
	}

	ingCh := fridgedoordatabase.Parse(ctx, cur, &Ingredient{})

	results := make([]*Ingredient, 0)

	for i := range ingCh {
		results = append(results, i.(*Ingredient))
	}

	return results, nil
}

// FindOne does not find one
func FindOne(ctx context.Context, id string) (*Ingredient, error) {

	connected, collection := collection()
	if !connected {
		return nil, errNotConnected
	}

	singleResult, err := collection.FindByID(ctx, id)

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &Ingredient{})

	return ing.(*Ingredient), err
}
