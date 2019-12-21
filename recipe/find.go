package recipe

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne finds a recipe by ID
func FindOne(ctx context.Context, id string) (*Recipe, error) {

	connected, collection := collection()
	if !connected {
		return nil, errNotConnected
	}

	singleResult, err := collection.FindByID(ctx, id)

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &Recipe{})
	if err != nil {
		return nil, err
	}

	return ing.(*Recipe), err
}

// FindByIds finds recipe by ID
func FindByIds(ctx context.Context, ids []primitive.ObjectID) ([]*Description, error) {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(25)

	_in := bson.M{"$in": ids}
	idin := bson.M{"_id": _in}

	cur, err := mongoCollection.Find(context.Background(), idin, findOptions)
	if err != nil {
		return make([]*Description, 0), err
	}

	return parseRecipe(ctx, cur)
}

// FindByName finds recipes starting with the given letter
func FindByName(ctx context.Context, startsWith string, userID primitive.ObjectID) ([]*Recipe, error) {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(20)

	regex := bson.M{"$regex": primitive.Regex{Pattern: "\\b" + startsWith, Options: "i"}}
	startsWithBson := bson.M{"name": regex}
	addedByBson := bson.M{"addedBy": userID}
	andBson := bson.M{"$and": []bson.M{startsWithBson, addedByBson}}

	cur, err := mongoCollection.Find(ctx, andBson, findOptions)
	if err != nil {
		return make([]*Recipe, 0), err
	}

	recipeCh := fridgedoordatabase.Parse(ctx, cur, &Recipe{})

	results := make([]*Recipe, 0)

	for i := range recipeCh {
		results = append(results, i.(*Recipe))
	}

	return results, nil
}

// List lists all the available recipe
func List(ctx context.Context) ([]*Description, error) {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	duration3s, _ := time.ParseDuration("3s")
	findctx, cancelFunc := context.WithTimeout(ctx, duration3s)
	defer cancelFunc()

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(25)

	cur, err := mongoCollection.Find(findctx, bson.D{{}}, findOptions)
	if err != nil {
		return make([]*Description, 0), err
	}

	return parseRecipe(ctx, cur)
}

func parseRecipe(ctx context.Context, cur *mongo.Cursor) ([]*Description, error) {
	ingCh := fridgedoordatabase.Parse(ctx, cur, &Description{})

	results := make([]*Description, 0)

	for i := range ingCh {
		results = append(results, i.(*Description))
	}

	return results, nil
}
