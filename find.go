package fridgedoordatabase

import (
	"context"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindByID finds a single result using the collection ID
func (coll *Collection) FindByID(ctx context.Context, id string) (*mongo.SingleResult, error) {
	findOneOptions := options.FindOne()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	singleResult := coll.Collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: objID}}, findOneOptions)

	return singleResult, nil
}

// ParseSingleResult parses result returned by FindOne
func ParseSingleResult(singleResult *mongo.SingleResult, obj interface{}) (interface{}, error) {

	err := singleResult.Err()
	if err != nil {
		return nil, err
	}

	err = singleResult.Decode(obj)

	if err != nil {
		return nil, err
	}

	return obj, nil
}

// Parse parses cursor returned by Find
func Parse(ctx context.Context, cur *mongo.Cursor, obj interface{}) <-chan interface{} {

	objectType := reflect.TypeOf(obj).Elem()

	ch := make(chan interface{})

	go func() {
		defer close(ch)

		// Finding multiple documents returns a cursor
		// Iterating through the cursor allows us to decode documents one at a time
		for cur.Next(ctx) {

			// create a value into which the single document can be decoded
			result := reflect.New(objectType).Interface()
			err := cur.Decode(result)

			if err != nil {
				log.Fatal(err)
			}

			ch <- result
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		// Close the cursor once finished
		cur.Close(ctx)
	}()

	return ch
}
