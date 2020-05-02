package recipe

import (
	"context"
	"errors"
	"fmt"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddIngredient adds an ingredient to a recipe
func AddIngredient(ctx context.Context, userID primitive.ObjectID, recipeID *primitive.ObjectID, stepIdx int, ingredientID string, ingredient string) (*dfdmodels.Recipe, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	recipe, methodStep, err := coll.getMethodStepByID(ctx, recipeID, userID, stepIdx)
	if err != nil {
		fmt.Printf("Error retreiving method step, %v.\n", err)
		return nil, err
	}

	if !CanEdit(recipe, userID) {
		fmt.Println("User not authorised to update recipe")
		return nil, errUnauthorised
	}

	if containsIngredient(methodStep, ingredientID) {
		return nil, errors.New("Duplicate")
	}

	ing := dfdmodels.Ingredient{
		Name:         ingredient,
		IngredientID: ingredientID,
	}

	methodStep.Ingredients = append(methodStep.Ingredients, ing)
	recipe.Method[stepIdx] = *methodStep

	err = coll.c.UpdateByID(ctx, recipeID, recipe)
	if err != nil {
		return nil, err
	}

	return coll.findOne(ctx, recipeID, userID)
}

// UpdateIngredient removes ingredient from recipe
func UpdateIngredient(ctx context.Context, userID primitive.ObjectID, recipeID *primitive.ObjectID, stepIdx int, ingredientID string, updates map[string]string) (*dfdmodels.Recipe, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	recipe, methodStep, err := coll.getMethodStepByID(ctx, recipeID, userID, stepIdx)
	if err != nil {
		fmt.Printf("Error retreiving method step, %v.\n", err)
		return nil, err
	}

	if !CanEdit(recipe, userID) {
		fmt.Println("User not authorised to update recipe")
		return nil, errUnauthorised
	}

	methodStep.Ingredients = updateByID(methodStep.Ingredients, ingredientID, updates)
	recipe.Method[stepIdx] = *methodStep

	err = coll.c.UpdateByID(ctx, recipeID, recipe)
	if err != nil {
		return nil, err
	}

	return coll.findOne(ctx, recipeID, userID)
}

// RemoveIngredient removes ingredient from recipe
func RemoveIngredient(ctx context.Context, userID primitive.ObjectID, recipeID *primitive.ObjectID, stepIdx int, ingredientID string) (*dfdmodels.Recipe, error) {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return nil, errNotConnected
	}

	recipe, methodStep, err := coll.getMethodStepByID(ctx, recipeID, userID, stepIdx)
	if err != nil {
		fmt.Printf("Error retreiving method step, %v.\n", err)
		return nil, err
	}

	if !CanEdit(recipe, userID) {
		fmt.Println("User not authorised to update recipe")
		return nil, errUnauthorised
	}

	filterFn := func(id *dfdmodels.Ingredient) bool {
		return id.IngredientID != ingredientID
	}

	methodStep.Ingredients = filterIngredients(methodStep.Ingredients, filterFn)
	recipe.Method[stepIdx] = *methodStep

	err = coll.c.UpdateByID(ctx, recipeID, recipe)

	if err != nil {
		return nil, err
	}

	return coll.findOne(ctx, recipeID, userID)
}

func containsIngredient(r *dfdmodels.MethodStep, ingredientID string) bool {
	for _, ing := range r.Ingredients {
		if ing.IngredientID == ingredientID {
			return true
		}
	}

	return false
}

func filterIngredients(ings []dfdmodels.Ingredient, filterFn func(ing *dfdmodels.Ingredient) bool) []dfdmodels.Ingredient {
	filtered := []dfdmodels.Ingredient{}

	for _, ing := range ings {
		if filterFn(&ing) {
			filtered = append(filtered, ing)
		}
	}

	return filtered
}

func updateByID(ings []dfdmodels.Ingredient, ingredientID string, updates map[string]string) []dfdmodels.Ingredient {
	updated := make([]dfdmodels.Ingredient, len(ings))

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
