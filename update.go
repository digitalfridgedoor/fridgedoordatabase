package fridgedoordatabase

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateByID finds and updates an object by ID
func (coll *Collection) UpdateByID(ctx context.Context, id string, obj interface{}) error {

	o := options.FindOneAndReplace()

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	singleResult := coll.Collection.FindOneAndReplace(ctx, filter, obj, o)

	return singleResult.Err()
}
