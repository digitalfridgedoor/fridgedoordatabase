package ingredient

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteByID deletes an ingredient
func DeleteByID(ctx context.Context, id primitive.ObjectID) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	return collection.DeleteByID(ctx, id)
}
