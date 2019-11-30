package userview

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errUserExists = errors.New("User exists")

// Create creates a new userview for a user
func (coll *Collection) Create(ctx context.Context, username string) (*View, error) {
	_, err := coll.GetByUsername(ctx, username)
	if err == nil {
		// found user with that username
		return nil, errUserExists
	}

	insertOneOptions := options.InsertOne()

	collections := make(map[string]*RecipeCollection)
	view := &View{
		Username:    username,
		Collections: collections,
	}

	insertOneResult, err := coll.mongoCollection().InsertOne(ctx, view, insertOneOptions)
	if err != nil {
		return nil, err
	}

	insertedID := insertOneResult.InsertedID.(primitive.ObjectID)

	return coll.FindOne(ctx, insertedID.Hex())
}

// Delete removes a userview for a user
func (coll *Collection) Delete(ctx context.Context, username string) error {
	deleteOptions := options.Delete()

	view, err := coll.GetByUsername(ctx, username)
	if err != nil {
		return err
	}

	viewID := view.ID

	_, err = coll.mongoCollection().DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: viewID}}, deleteOptions)
	if err != nil {
		return err
	}

	return nil
}
