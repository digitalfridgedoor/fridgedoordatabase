package recipe

import (
	"context"

	"digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection can find and parse Recipes from mongodb
type Connection struct {
	db fridgedoordatabase.Connection
}

// FindOne does not find one
func (conn *Connection) FindOne(ctx context.Context, id string) (*Recipe, error) {

	collection := conn.db.Collection("recipeapi", "recipes")

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
