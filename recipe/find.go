package recipe

import (
	"context"
	"fmt"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindOne finds a recipe by ID
func FindOne(ctx context.Context, id *primitive.ObjectID, userID primitive.ObjectID) (*dfdmodels.Recipe, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	return coll.findOne(ctx, id, userID)
}

func (coll *collection) findOne(ctx context.Context, id *primitive.ObjectID, userID primitive.ObjectID) (*dfdmodels.Recipe, error) {

	r, err := coll.c.FindByID(ctx, id, &dfdmodels.Recipe{})

	if err != nil {
		return nil, err
	}

	re := r.(*dfdmodels.Recipe)
	if !CanView(re, userID) {
		return nil, errUnauthorised
	}

	return re, err
}
