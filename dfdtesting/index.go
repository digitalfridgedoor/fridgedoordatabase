package dfdtesting

import (
	"github.com/maisiesadler/theilliminationgame/database"
	"github.com/maisiesadler/theilliminationgame/models"
	"go.mongodb.org/mongo-driver/bson"
)

var overrides = make(map[string]*TestCollection)

// SetTestCollectionOverride sets a the database package to use a TestCollection
func SetTestCollectionOverride() {
	database.SetOverride(overrideDb)
}

// SetUserViewFindPredicate overrides the logic to get the result for Find
func SetUserViewFindPredicate(predicate func(*models.UserView, bson.M) bool) bool {
	fn := func(value interface{}, filter bson.M) bool {
		uv := value.(*models.UserView)
		return predicate(uv, filter)
	}

	coll := getOrAddTestCollection("recipeapi", "userviews")
	coll.findPredicate = fn
	return true
}

func overrideDb(database string, collection string) database.ICollection {
	return getOrAddTestCollection(database, collection)
}

func getOrAddTestCollection(database string, collection string) *TestCollection {
	key := database + "_" + collection
	if val, ok := overrides[key]; ok {
		return val
	}
	overrides[key] = CreateTestCollection()
	return overrides[key]
}
