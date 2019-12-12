package ingredient

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create creates a new ingredient with given name
func (coll *Collection) Create(ctx context.Context, name string) (*Ingredient, error) {

	insertOneOptions := options.InsertOne()

	ingredient := &Ingredient{
		Name:    name,
		AddedOn: time.Now(),
	}

	insertOneResult, err := coll.mongoCollection().InsertOne(ctx, ingredient, insertOneOptions)
	if err != nil {
		return nil, err
	}

	insertedID := insertOneResult.InsertedID.(primitive.ObjectID)

	return coll.FindOne(ctx, insertedID.Hex())
}
