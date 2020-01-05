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
func FindOne(ctx context.Context, viewID string) (*View, error) {

	connected, collection := collection()
	if !connected {
		return nil, errNotConnected
	}

	singleResult, err := collection.FindByID(ctx, viewID)
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
func GetByUsername(ctx context.Context, username string) (*View, error) {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	// Pass these options to the FindOne method
	findOneOptions := options.FindOne()

	singleResult := mongoCollection.FindOne(ctx, bson.D{primitive.E{Key: "username", Value: username}}, findOneOptions)

	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &View{})
	if err != nil {
		return nil, err
	}

	return ing.(*View), err
}

// GetLinkedUserViews returns all user views for now
func GetLinkedUserViews(ctx context.Context) ([]*View, error) {

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
		return nil, err
	}
	return parseView(ctx, cur)
}

func parseView(ctx context.Context, cur *mongo.Cursor) ([]*View, error) {
	ingCh := fridgedoordatabase.Parse(ctx, cur, &View{})

	results := make([]*View, 0)

	for i := range ingCh {
		results = append(results, i.(*View))
	}

	return results, nil
}
