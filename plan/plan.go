package plan

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

func collection() (bool, *fridgedoordatabase.Collection) {
	return fridgedoordatabase.CreateCollection("recipeapi", "plans")
}

func mongoCollection() (bool, *mongo.Collection) {
	connected, collection := collection()
	if !connected {
		return false, nil
	}
	return true, collection.MongoCollection
}
