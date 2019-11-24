package recipe

// import (
// 	"context"
// 	"digitalfridgedoor/fridgedoordatabase"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // Create creates a new recipe with given name
// func (conn *Connection) Create(ctx context.Context, userID string, name string) (*Recipe, error) {

// 	collection := conn.collection()

// 	insertOneOptions := options.InsertOne()

// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	singleResult := collection.InsertOne(ctx, bson.D{primitive.E{Key: "_id", Value: objID}}, findOneOptions)

// 	ing, err := fridgedoordatabase.ParseSingleResult(singleResult, &Recipe{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ing.(*Recipe), err
// }
