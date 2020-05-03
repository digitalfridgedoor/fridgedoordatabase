package dfdtesting

import (
	"regexp"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindIngredientByNameTestPredicate can be used with SetFindFilter for searching ingredients by name
func FindIngredientByNameTestPredicate(gs *dfdmodels.Ingredient, m primitive.M) bool {
	nameval := m["name"].(bson.M)
	regexval := nameval["$regex"].(primitive.Regex)

	r := regexp.MustCompile(regexval.Pattern)

	return r.MatchString(gs.Name)
}
