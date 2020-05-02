package recipe

import (
	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CanView returns whether the specfified user is authorized to view the recipe
func CanView(r *dfdmodels.Recipe, userID primitive.ObjectID) bool {
	if r.AddedBy == userID || r.Metadata.ViewableBy.Everyone {
		return true
	}

	for _, id := range r.Metadata.ViewableBy.Users {
		if id == userID {
			return true
		}
	}

	return false
}

// CanEdit returns whether the specfified user is authorized to edit the recipe
func CanEdit(r *dfdmodels.Recipe, userID primitive.ObjectID) bool {
	return r.AddedBy == userID
}
