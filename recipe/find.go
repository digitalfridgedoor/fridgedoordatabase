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
func (coll *Collection) FindOne(ctx context.Context, id string) (*Recipe, error) {

	singleResult, err := coll.collection.FindByID(ctx, id)

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &Recipe{})
	if err != nil {
		return nil, err
	}

	return ing.(*Recipe), err
}

// FindByIds finds recipe by ID
func (coll *Collection) FindByIds(ctx context.Context, ids []*primitive.ObjectID) ([]*Description, error) {

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(25)

	_in := bson.M{"$in": ids}
	idin := bson.M{"_id": _in}

	cur, err := coll.mongoCollection().Find(context.Background(), idin, findOptions)
	if err != nil {
		return make([]*Description, 0), err
	}

	return parseRecipe(ctx, cur)
}

// List lists all the available recipe
func (coll *Collection) List(ctx context.Context) ([]*Description, error) {

	duration3s, _ := time.ParseDuration("3s")
	findctx, cancelFunc := context.WithTimeout(ctx, duration3s)
	defer cancelFunc()

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(25)

	cur, err := coll.mongoCollection().Find(findctx, bson.D{{}}, findOptions)
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
