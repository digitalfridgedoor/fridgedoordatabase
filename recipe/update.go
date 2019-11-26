package recipe

import (
	"context"
	"errors"
)

// AddIngredient adds recipe to users list
func (coll *Collection) AddIngredient(ctx context.Context, recipeID string, ingredientID string, ingredient string) error {

	recipe, err := coll.FindOne(ctx, recipeID)
	if err != nil {
		return err
	}

	if recipe.containsIngredient(ingredientID) {
		return errors.New("Duplicate")
	}

	ing := Ingredient{
		Name:         ingredient,
		IngredientID: ingredientID,
	}

	recipe.Ingredients = append(recipe.Ingredients, ing)

	return coll.collection.UpdateByID(ctx, recipeID, recipe)
}

// RemoveIngredient removes ingredient from recipe
func (coll *Collection) RemoveIngredient(ctx context.Context, recipeID string, ingredientID string) error {
	recipe, err := coll.FindOne(ctx, recipeID)
	if err != nil {
		return err
	}

	filterFn := func(id *Ingredient) bool {
		return id.IngredientID != ingredientID
	}

	recipe.Ingredients = filterIngredients(recipe.Ingredients, filterFn)

	return coll.collection.UpdateByID(ctx, recipeID, recipe)
}

func (r *Recipe) containsIngredient(ingredientID string) bool {
	for _, ing := range r.Ingredients {
		if ing.IngredientID == ingredientID {
			return true
		}
	}

	return false
}

func filterIngredients(ings []Ingredient, filterFn func(ing *Ingredient) bool) []Ingredient {
	filtered := []Ingredient{}

	for ing := range iterate(ings) {
		if filterFn(ing) {
			filtered = append(filtered, *ing)
		}
	}

	return filtered
}

func iterate(ings []Ingredient) <-chan *Ingredient {
	ch := make(chan *Ingredient)

	go func() {
		defer close(ch)
		for _, id := range ings {
			ch <- &id
		}
	}()

	return ch
}
