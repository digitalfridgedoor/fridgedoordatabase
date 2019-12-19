package ingredient

import (
	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/mongo"
)

func collection() (bool, *fridgedoordatabase.Collection) {
	return fridgedoordatabase.CreateCollection("recipeapi", "ingredients")
}

func mongoCollection() (bool, *mongo.Collection) {
	connected, collection := collection()
	if !connected {
		return false, nil
	}
	return true, collection.MongoCollection
}
