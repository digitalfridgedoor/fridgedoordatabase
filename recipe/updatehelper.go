package recipe

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func canEdit(recipe *Recipe, userID primitive.ObjectID) bool {
	return recipe.AddedBy == userID
}
