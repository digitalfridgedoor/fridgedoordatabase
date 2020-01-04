package recipe

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddIngredient adds an ingredient to a recipe
func AddIngredient(ctx context.Context, user primitive.ObjectID, recipeID string, stepIdx int, ingredientID string, ingredient string) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	recipe, methodStep, err := getMethodStepByID(ctx, recipeID, stepIdx)
	if err != nil {
		fmt.Printf("Error retreiving method step, %v.\n", err)
		return err
	}

	if !canEdit(recipe, user) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	if methodStep.containsIngredient(ingredientID) {
		return errors.New("Duplicate")
	}

	ing := Ingredient{
		Name:         ingredient,
		IngredientID: ingredientID,
	}

	methodStep.Ingredients = append(methodStep.Ingredients, ing)
	recipe.Method[stepIdx] = *methodStep

	return collection.UpdateByID(ctx, recipeID, recipe)
}

// UpdateIngredient removes ingredient from recipe
func UpdateIngredient(ctx context.Context, user primitive.ObjectID, recipeID string, stepIdx int, ingredientID string, updates map[string]string) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	recipe, methodStep, err := getMethodStepByID(ctx, recipeID, stepIdx)
	if err != nil {
		fmt.Printf("Error retreiving method step, %v.\n", err)
		return err
	}

	if !canEdit(recipe, user) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	methodStep.Ingredients = updateByID(methodStep.Ingredients, ingredientID, updates)
	recipe.Method[stepIdx] = *methodStep

	return collection.UpdateByID(ctx, recipeID, recipe)
}

// RemoveIngredient removes ingredient from recipe
func RemoveIngredient(ctx context.Context, user primitive.ObjectID, recipeID string, stepIdx int, ingredientID string) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	recipe, methodStep, err := getMethodStepByID(ctx, recipeID, stepIdx)
	if err != nil {
		fmt.Printf("Error retreiving method step, %v.\n", err)
		return err
	}

	if !canEdit(recipe, user) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	filterFn := func(id *Ingredient) bool {
		return id.IngredientID != ingredientID
	}

	methodStep.Ingredients = filterIngredients(methodStep.Ingredients, filterFn)
	recipe.Method[stepIdx] = *methodStep

	return collection.UpdateByID(ctx, recipeID, recipe)
}

func (r *MethodStep) containsIngredient(ingredientID string) bool {
	for _, ing := range r.Ingredients {
		if ing.IngredientID == ingredientID {
			return true
		}
	}

	return false
}

func filterIngredients(ings []Ingredient, filterFn func(ing *Ingredient) bool) []Ingredient {
	filtered := []Ingredient{}

	for _, ing := range ings {
		if filterFn(&ing) {
			filtered = append(filtered, ing)
		}
	}

	return filtered
}

func updateByID(ings []Ingredient, ingredientID string, updates map[string]string) []Ingredient {
	updated := make([]Ingredient, len(ings))

	for index, ing := range ings {
		if ing.IngredientID == ingredientID {
			if update, ok := updates["name"]; ok {
				ing.Name = update
			}
			if update, ok := updates["amount"]; ok {
				ing.Amount = update
			}
			if update, ok := updates["preperation"]; ok {
				ing.Preperation = update
			}
		}
		updated[index] = ing
	}

	return updated
}
