package recipe

import (
	"context"
	"fmt"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddSubRecipe adds a link between the recipe and the subrecipe
func AddSubRecipe(ctx context.Context, user primitive.ObjectID, recipeID *primitive.ObjectID, subRecipeID *primitive.ObjectID) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	if *recipeID == *subRecipeID {
		return errSubRecipe
	}

	recipe, err := coll.findOne(ctx, recipeID)
	if err != nil {
		return err
	}

	if !CanEdit(recipe, user) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	// todo: append parent so we know when we unlink?
	if len(recipe.ParentIds) > 0 {
		fmt.Println("Cannot add subrecipe to subrecipe")
		return errSubRecipe
	}

	if hasSubRecipe(recipe, subRecipeID) {
		return errDuplicate
	}

	subRecipe, err := FindOne(ctx, subRecipeID)
	if err != nil {
		return err
	}

	if len(subRecipe.Recipes) != 0 {
		return errSubRecipe
	}

	subRecipe.ParentIds = appendParentRecipeID(subRecipe.ParentIds, *recipe.ID)
	err = coll.c.UpdateByID(ctx, subRecipeID, subRecipe)
	if err != nil {
		fmt.Printf("Error updating subrecipe: %v\n", err)
		return err
	}

	recipe.Recipes = append(recipe.Recipes, dfdmodels.SubRecipe{
		RecipeID: *subRecipe.ID,
		Name:     subRecipe.Name,
	})

	return coll.c.UpdateByID(ctx, recipeID, recipe)
}

// RemoveSubRecipe the link between the recipe/subrecipe
func RemoveSubRecipe(ctx context.Context, user primitive.ObjectID, recipeID *primitive.ObjectID, subRecipeID *primitive.ObjectID) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	recipe, err := coll.findOne(ctx, recipeID)
	if err != nil {
		return err
	}

	if !CanEdit(recipe, user) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	filterFn := func(id *dfdmodels.SubRecipe) bool {
		return id.RecipeID != *subRecipeID
	}

	recipe.Recipes = filterSubRecipes(recipe.Recipes, filterFn)

	subRecipe, err := coll.findOne(ctx, subRecipeID)
	if err == nil {
		subRecipe.ParentIds = removeParentRecipeID(subRecipe.ParentIds, *recipe.ID)
		err = coll.c.UpdateByID(ctx, subRecipeID, subRecipe)
		if err != nil {
			fmt.Printf("Error updating subrecipe: %v.", err)
		}
	} else {
		fmt.Printf("Could not find subrecipe with id=%v.\n", subRecipeID)
	}

	return coll.c.UpdateByID(ctx, recipeID, recipe)
}

func hasSubRecipe(r *dfdmodels.Recipe, subRecipeID *primitive.ObjectID) bool {
	for _, subrecipe := range r.Recipes {
		if subrecipe.RecipeID == *subRecipeID {
			return true
		}
	}

	return false
}

func filterSubRecipes(subRecipes []dfdmodels.SubRecipe, filterFn func(ing *dfdmodels.SubRecipe) bool) []dfdmodels.SubRecipe {
	filtered := []dfdmodels.SubRecipe{}

	for _, sr := range subRecipes {
		if filterFn(&sr) {
			filtered = append(filtered, sr)
		}
	}

	return filtered
}

func appendParentRecipeID(parentIds []primitive.ObjectID, parentID primitive.ObjectID) []primitive.ObjectID {
	hasParentID := false

	for _, id := range parentIds {
		if id == parentID {
			hasParentID = true
		}
	}

	if !hasParentID {
		parentIds = append(parentIds, parentID)
	}

	return parentIds
}

func removeParentRecipeID(parentIds []primitive.ObjectID, parentID primitive.ObjectID) []primitive.ObjectID {
	filtered := []primitive.ObjectID{}

	for _, id := range parentIds {
		if id != parentID {
			filtered = append(filtered, id)
		}
	}

	return filtered
}
