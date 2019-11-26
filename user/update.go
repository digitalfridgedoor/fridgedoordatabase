package user

import (
	"context"
	"digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddRecipe adds recipe to users list
func (coll *Collection) AddRecipe(ctx context.Context, userID string, recipeID primitive.ObjectID) error {

	user, err := coll.FindOne(ctx, userID)
	if err != nil {
		return err
	}

	user.Recipes = append(user.Recipes, recipeID)

	return coll.collection.UpdateByID(ctx, userID, user)
}

// RemoveRecipe removes recipe from users list
func (coll *Collection) RemoveRecipe(ctx context.Context, userID string, recipeID primitive.ObjectID) error {
	user, err := coll.FindOne(ctx, userID)
	if err != nil {
		return err
	}

	filterFn := func(id *primitive.ObjectID) bool {
		return *id != recipeID
	}

	user.Recipes = fridgedoordatabase.Filter(user.Recipes, filterFn)

	return coll.collection.UpdateByID(ctx, userID, user)
}
