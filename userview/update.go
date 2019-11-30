package userview

import (
	"context"

	"github.com/digitalfridgedoor/fridgedoordatabase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddRecipe adds recipe to users list
func (coll *Collection) AddRecipe(ctx context.Context, viewID string, collectionName string, recipeID primitive.ObjectID) error {

	view, err := coll.FindOne(ctx, viewID)
	if err != nil {
		return err
	}

	if _, ok := view.Collections[collectionName]; !ok {
		view.Collections[collectionName] = &RecipeCollection{Name: collectionName}
	}

	view.Collections[collectionName].addRecipe(recipeID)

	return coll.collection.UpdateByID(ctx, viewID, view)
}

// RemoveRecipe removes recipe from users list
func (coll *Collection) RemoveRecipe(ctx context.Context, viewID string, collectionName string, recipeID primitive.ObjectID) error {
	view, err := coll.FindOne(ctx, viewID)
	if err != nil {
		return err
	}

	filterFn := func(id *primitive.ObjectID) bool {
		return *id != recipeID
	}

	if viewCollection, ok := view.Collections[collectionName]; !ok {
		viewCollection.Recipes = fridgedoordatabase.Filter(viewCollection.Recipes, filterFn)
	}

	return coll.collection.UpdateByID(ctx, viewID, view)
}

func (r *RecipeCollection) addRecipe(recipeID primitive.ObjectID) {
	r.Recipes = append(r.Recipes, recipeID)
}
