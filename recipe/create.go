package recipe

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create creates a new recipe with given name
func Create(ctx context.Context, userID primitive.ObjectID, name string) (*Recipe, error) {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	insertOneOptions := options.InsertOne()

	recipe := &Recipe{
		Name:    name,
		AddedOn: time.Now(),
		AddedBy: userID,
	}

	insertOneResult, err := mongoCollection.InsertOne(ctx, recipe, insertOneOptions)
	if err != nil {
		return nil, err
	}

	insertedID := insertOneResult.InsertedID.(primitive.ObjectID)

	return FindOne(ctx, insertedID.Hex())
}

// Delete removes a recipe
func Delete(ctx context.Context, recipeID primitive.ObjectID) error {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return errNotConnected
	}

	deleteOptions := options.Delete()

	_, err := mongoCollection.DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: recipeID}}, deleteOptions)
	if err != nil {
		return err
	}

	return nil
}
