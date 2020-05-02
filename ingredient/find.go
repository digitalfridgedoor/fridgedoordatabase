package ingredient

import (
	"context"
	"fmt"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindByName finds ingredients starting with the given letter
func FindByName(ctx context.Context, startsWith string) ([]*dfdmodels.Ingredient, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(20)

	regex := bson.M{"$regex": primitive.Regex{Pattern: "\\b" + startsWith, Options: "i"}}
	startsWithBson := bson.M{"name": regex}

	ch, err := coll.c.Find(ctx, startsWithBson, findOptions, &dfdmodels.Ingredient{})
	if err != nil {
		return make([]*dfdmodels.Ingredient, 0), err
	}

	results := make([]*dfdmodels.Ingredient, 0)

	for i := range ch {
		results = append(results, i.(*dfdmodels.Ingredient))
	}

	return results, nil
}

// FindOne does not find one
func FindOne(ctx context.Context, id *primitive.ObjectID) (*dfdmodels.Ingredient, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	return coll.findOne(ctx, id)
}

func (coll *collection) findOne(ctx context.Context, id *primitive.ObjectID) (*dfdmodels.Ingredient, error) {

	ing, err := coll.c.FindByID(ctx, id, &dfdmodels.Ingredient{})

	return ing.(*dfdmodels.Ingredient), err
}
