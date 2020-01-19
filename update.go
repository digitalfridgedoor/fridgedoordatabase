package fridgedoordatabase

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne inserts a new document into the collection
func (coll *Collection) InsertOne(ctx context.Context, document interface{}) (*primitive.ObjectID, error) {
	insertOneOptions := options.InsertOne()

	insertOneResult, err := coll.MongoCollection.InsertOne(ctx, document, insertOneOptions)
	if err != nil {
		return nil, err
	}

	insertedID := insertOneResult.InsertedID.(primitive.ObjectID)
	return &insertedID, nil
}

// UpdateByID finds and updates an object by ID
func (coll *Collection) UpdateByID(ctx context.Context, id string, obj interface{}) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	o := options.FindOneAndReplace()

	filter := bson.D{primitive.E{Key: "_id", Value: objID}}

	singleResult := coll.MongoCollection.FindOneAndReplace(ctx, filter, obj, o)

	return singleResult.Err()
}
