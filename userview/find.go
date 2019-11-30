package userview

import (
	"context"

	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
