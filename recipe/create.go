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

	recipe := &dfdmodels.Recipe{
		Name:    name,
		AddedOn: time.Now(),
		AddedBy: userID,
	}

	r, err := coll.c.InsertOneAndFind(ctx, recipe, &dfdmodels.Recipe{})
	if err != nil {
		fmt.Printf("Error creating recipe, %v\n", err)
		return nil, err
	}

	return r.(*dfdmodels.Recipe), nil
}

// Delete removes a recipe
func Delete(ctx context.Context, recipeID *primitive.ObjectID) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	return coll.c.DeleteByID(ctx, recipeID)
}
