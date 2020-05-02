package recipe

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
)

// Create creates a new recipe with given name
func Create(ctx context.Context, userID primitive.ObjectID, name string) (*dfdmodels.Recipe, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	recipe := &Recipe{
		Name:    name,
		AddedOn: time.Now(),
		AddedBy: userID,
	}

	r, err := coll.c.InsertOneAndFind(ctx, recipe, &Recipe{})

	return r.(*dfdmodels.Recipe), nil
}

// Delete removes a recipe
func Delete(ctx context.Context, recipeID *primitive.ObjectID) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	_, err := coll.c.DeleteByID(ctx, recipeID)

	return err
}
