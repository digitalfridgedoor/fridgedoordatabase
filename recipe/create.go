package recipe

import (
	"context"
	"time"

	"github.com/digitalfridgedoor/fridgedoordatabase/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create creates a new recipe with given name
func (conn *Connection) Create(ctx context.Context, userID string, name string) (*Recipe, error) {

	u := user.New(conn.db)
	userInfo, err := u.FindOne(ctx, userID)
	if err != nil {
		return nil, err
	}

	collection := conn.collection()

	insertOneOptions := options.InsertOne()

	recipe := &Recipe{
		Name:    name,
		AddedOn: time.Now(),
		AddedBy: userInfo.ID,
	}

	insertOneResult, err := collection.InsertOne(ctx, recipe, insertOneOptions)
	if err != nil {
		return nil, err
	}

	insertedID := insertOneResult.InsertedID.(primitive.ObjectID)

	err = u.AddRecipeToUser(ctx, userID, insertedID)
	if err != nil {
		return nil, err
	}

	return conn.FindOne(ctx, insertedID.Hex())
}

// Delete removes a recipe
func (conn *Connection) Delete(ctx context.Context, recipeID primitive.ObjectID) error {

	collection := conn.collection()

	deleteOptions := options.Delete()

	_, err := collection.DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: recipeID}}, deleteOptions)
	if err != nil {
		return err
	}

	return nil
}
