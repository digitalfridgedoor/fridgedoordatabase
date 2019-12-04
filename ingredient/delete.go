package ingredient

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DeleteByID deletes an ingredient
func (coll *Collection) DeleteByID(ctx context.Context, id primitive.ObjectID) error {

	deleteOptions := options.Delete()

	_, err := coll.mongoCollection().DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}, deleteOptions)
	if err != nil {
		return err
	}

	return nil
}
