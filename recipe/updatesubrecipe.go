package recipe

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddSubRecipe adds a link between the recipe and the subrecipe
func AddSubRecipe(ctx context.Context, recipeID string, subRecipeID string) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	if recipeID == subRecipeID {
		return errSubRecipe
	}

	recipe, err := FindOne(ctx, recipeID)
	if err != nil {
		return err
	}

	// todo: append parent so we know when we unlink?
	if recipe.IsSubRecipe {
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

	subRecipe.IsSubRecipe = true
	err = collection.UpdateByID(ctx, recipeID, recipe)
	if err != nil {
		fmt.Printf("Error updating subrecipe: %v\n", err)
		return err
	}

	recipe.Recipes = append(recipe.Recipes, SubRecipe{
		RecipeID: subRecipe.ID,
		Name:     subRecipe.Name,
	})

	return collection.UpdateByID(ctx, recipeID, recipe)
}

// RemoveSubRecipe the link between the recipe/subrecipe
func RemoveSubRecipe(ctx context.Context, recipeID string, subRecipeID string) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	srID, err := primitive.ObjectIDFromHex(subRecipeID)
	if err != nil {
		fmt.Printf("Invalid object ID, %v.\n", subRecipeID)
		return err
	}

	recipe, err := FindOne(ctx, recipeID)
	if err != nil {
		return err
	}

	filterFn := func(id *SubRecipe) bool {
		return id.RecipeID != srID
	}

	recipe.Recipes = filterSubRecipes(recipe.Recipes, filterFn)

	return collection.UpdateByID(ctx, recipeID, recipe)
}

func hasSubRecipe(r *Recipe, subRecipeID string) bool {
	for _, subrecipe := range r.Recipes {
		if subrecipe.RecipeID.Hex() == subRecipeID {
			return true
		}
	}

	return false
}

func filterSubRecipes(subRecipes []SubRecipe, filterFn func(ing *SubRecipe) bool) []SubRecipe {
	filtered := []SubRecipe{}

	for _, sr := range subRecipes {
		if filterFn(&sr) {
			filtered = append(filtered, sr)
		}
	}

	return filtered
}
