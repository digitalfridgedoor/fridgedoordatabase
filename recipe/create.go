package recipe

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create creates a new recipe with given name
func (coll *Collection) Create(ctx context.Context, userID primitive.ObjectID, name string) (*Recipe, error) {

	insertOneOptions := options.InsertOne()

	recipe := &Recipe{
		Name:    name,
		AddedOn: time.Now(),
		AddedBy: userID,
	}

	insertOneResult, err := coll.mongoCollection().InsertOne(ctx, recipe, insertOneOptions)
	if err != nil {
		return nil, err
	}

	insertedID := insertOneResult.InsertedID.(primitive.ObjectID)

	return coll.FindOne(ctx, insertedID.Hex())
}

// Delete removes a recipe
func (coll *Collection) Delete(ctx context.Context, recipeID primitive.ObjectID) error {

	deleteOptions := options.Delete()

	_, err := coll.mongoCollection().DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: recipeID}}, deleteOptions)
	if err != nil {
		return err
	}

	return nil
}
