package userview

import (
	"context"

	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddRecipe adds recipe to users list
func AddRecipe(ctx context.Context, viewID string, collectionName string, recipeID primitive.ObjectID) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	view, err := FindOne(ctx, viewID)
	if err != nil {
		return err
	}

	if _, ok := view.Collections[collectionName]; !ok {
		view.Collections[collectionName] = &RecipeCollection{Name: collectionName}
	}

	view.Collections[collectionName].addRecipe(recipeID)

	return collection.UpdateByID(ctx, viewID, view)
}

// RemoveRecipe removes recipe from users list
func RemoveRecipe(ctx context.Context, viewID string, collectionName string, recipeID primitive.ObjectID) error {
	view, err := FindOne(ctx, viewID)
	if err != nil {
		return err
	}

	filterFn := func(id *primitive.ObjectID) bool {
		return *id != recipeID
	}

	if viewCollection, ok := view.Collections[collectionName]; !ok {
		viewCollection.Recipes = fridgedoordatabase.Filter(viewCollection.Recipes, filterFn)
	}

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	return collection.UpdateByID(ctx, viewID, view)
}

func (r *RecipeCollection) addRecipe(recipeID primitive.ObjectID) {
	r.Recipes = append(r.Recipes, recipeID)
}
