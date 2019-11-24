package recipe

import (
	"context"
	"digitalfridgedoor/fridgedoordatabase"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne finds a recipe by ID
func (conn *Connection) FindOne(ctx context.Context, id string) (*Recipe, error) {

	collection := conn.collection()

	// Pass these options to the FindOne method
	findOneOptions := options.FindOne()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	singleResult := collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: objID}}, findOneOptions)

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &Recipe{})
	if err != nil {
		return nil, err
	}

	return ing.(*Recipe), err
}

// List lists all the available recipe
func (conn *Connection) List(ctx context.Context) ([]*Description, error) {
	collection := conn.collection()

	duration3s, _ := time.ParseDuration("3s")
	findctx, cancelFunc := context.WithTimeout(ctx, duration3s)
	defer cancelFunc()

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(25)

	cur, err := collection.Find(findctx, bson.D{{}}, findOptions)
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
