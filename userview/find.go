package userview

import (
	"context"
	"time"

	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne finds a User by ID
func (coll *Collection) FindOne(ctx context.Context, viewID string) (*View, error) {

	singleResult, err := coll.collection.FindByID(ctx, viewID)
	if err != nil {
		return nil, err
	}

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &View{})
	if err != nil {
		return nil, err
	}

	return ing.(*View), err
}

// GetByUsername tries to get User by username
func (coll *Collection) GetByUsername(ctx context.Context, username string) (*View, error) {

	// Pass these options to the FindOne method
	findOneOptions := options.FindOne()

	singleResult := coll.mongoCollection().FindOne(ctx, bson.D{primitive.E{Key: "username", Value: username}}, findOneOptions)

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &View{})
	if err != nil {
		return nil, err
	}

	return ing.(*View), err
}

// GetLinkedUserViews returns all user views for now
func (coll *Collection) GetLinkedUserViews(ctx context.Context) ([]*View, error) {

	duration3s, _ := time.ParseDuration("3s")
	findctx, cancelFunc := context.WithTimeout(ctx, duration3s)
	defer cancelFunc()

	// Pass these options to the Find method
	findOptions := options.Find()
	findOptions.SetLimit(25)

	cur, err := coll.mongoCollection().Find(findctx, bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}
	return parseRecipe(ctx, cur)
}

func parseRecipe(ctx context.Context, cur *mongo.Cursor) ([]*View, error) {
	ingCh := fridgedoordatabase.Parse(ctx, cur, &View{})

	results := make([]*View, 0)

	for i := range ingCh {
		results = append(results, i.(*View))
	}

	return results, nil
}
