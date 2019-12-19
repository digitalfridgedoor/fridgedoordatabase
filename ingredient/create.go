package ingredient

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create creates a new ingredient with given name
func Create(ctx context.Context, name string) (*Ingredient, error) {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return nil, errNotConnected
	}

	insertOneOptions := options.InsertOne()

	ingredient := &Ingredient{
		Name:    name,
		AddedOn: time.Now(),
	}

	insertOneResult, err := mongoCollection.InsertOne(ctx, ingredient, insertOneOptions)
	if err != nil {
		return nil, err
	}

	insertedID := insertOneResult.InsertedID.(primitive.ObjectID)

	return FindOne(ctx, insertedID.Hex())
}
