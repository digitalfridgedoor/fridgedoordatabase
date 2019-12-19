package ingredient

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DeleteByID deletes an ingredient
func DeleteByID(ctx context.Context, id primitive.ObjectID) error {

	connected, mongoCollection := mongoCollection()
	if !connected {
		return errNotConnected
	}

	deleteOptions := options.Delete()

	_, err := mongoCollection.DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}, deleteOptions)
	if err != nil {
		return err
	}

	return nil
}
