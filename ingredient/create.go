package ingredient

import (
	"context"
	"fmt"
	"time"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
)

// Create creates a new ingredient with given name
func Create(ctx context.Context, name string) (*dfdmodels.Ingredient, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	ingredient := &dfdmodels.Ingredient{
		Name:    name,
		AddedOn: time.Now(),
	}

	insertedID, err := coll.c.InsertOne(ctx, ingredient)
	if err != nil {
		return nil, err
	}

	return coll.findOne(ctx, insertedID)
}
