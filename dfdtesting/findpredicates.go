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

// FindPlanByMonthAndYearTestPredicate can be used with SetFindFilter for get or create plan
func FindPlanByMonthAndYearTestPredicate(p *dfdmodels.Plan, m bson.M) bool {
	month := m["month"].(int)
	year := m["year"].(int)
	userid := m["userid"].(primitive.ObjectID)

	return month == p.Month && year == p.Year && userid == p.UserID
}

// SetUserViewFindByUsernamePredicate overrides logic for find users by username
func SetUserViewFindByUsernamePredicate() {
	SetUserViewFindPredicate(func(uv *dfdmodels.UserView, m primitive.M) bool {
		return m["username"] == uv.Username
	})
}
