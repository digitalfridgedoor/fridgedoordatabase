package ingredient

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteByID deletes an ingredient
func DeleteByID(ctx context.Context, id *primitive.ObjectID) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	return coll.c.DeleteByID(ctx, id)
}
