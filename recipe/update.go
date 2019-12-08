package recipe

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

// AddMethodStep adds new method step to a recipe
func (coll *Collection) AddMethodStep(ctx context.Context, recipeID string, action string) error {

	recipe, err := coll.FindOne(ctx, recipeID)
	if err != nil {
		return err
	}

	methodStep := MethodStep{
		Action: action,
	}

	recipe.Method = append(recipe.Method, methodStep)

	return coll.collection.UpdateByID(ctx, recipeID, recipe)
}

// UpdateMethodStepByIndex updates method step at index
func (coll *Collection) UpdateMethodStepByIndex(ctx context.Context, recipeID string, stepIdx int, updates map[string]string) error {

	recipe, methodStep, err := coll.getMethodStepByID(ctx, recipeID, stepIdx)
	if err != nil {
		fmt.Printf("Error retreiving method step, %v.\n", err)
		return err
	}

	recipe.Method[stepIdx] = *updateMethodStep(methodStep, updates)

	return coll.collection.UpdateByID(ctx, recipeID, recipe)
}

// RemoveMethodStepByIndex removes method by index
func (coll *Collection) RemoveMethodStepByIndex(ctx context.Context, recipeID string, stepIdx int) error {

	if stepIdx < 0 {
		return errors.New("Invalid index")
	}

	recipe, err := coll.FindOne(ctx, recipeID)
	if err != nil {
		return err
	}

	if len(recipe.Method) <= stepIdx {
		return errors.New("Invalid index")
	}

	copy(recipe.Method[stepIdx:], recipe.Method[stepIdx+1:]) // Shift a[i+1:] left one index.
	recipe.Method = recipe.Method[:len(recipe.Method)-1]     // Truncate slice.

	return coll.collection.UpdateByID(ctx, recipeID, recipe)
}

// AddIngredient adds an ingredient to a recipe
func (coll *Collection) AddIngredient(ctx context.Context, recipeID string, stepIdx int, ingredientID string, ingredient string) error {

	recipe, methodStep, err := coll.getMethodStepByID(ctx, recipeID, stepIdx)
	if err != nil {
		fmt.Printf("Error retreiving method step, %v.\n", err)
		return err
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

	return coll.collection.UpdateByID(ctx, recipeID, recipe)
}

// UpdateIngredient removes ingredient from recipe
func (coll *Collection) UpdateIngredient(ctx context.Context, recipeID string, stepIdx int, ingredientID string, updates map[string]string) error {

	recipe, methodStep, err := coll.getMethodStepByID(ctx, recipeID, stepIdx)
	if err != nil {
		fmt.Printf("Error retreiving method step, %v.\n", err)
		return err
	}

	methodStep.Ingredients = updateByID(methodStep.Ingredients, ingredientID, updates)
	recipe.Method[stepIdx] = *methodStep

	return coll.collection.UpdateByID(ctx, recipeID, recipe)
}

// RemoveIngredient removes ingredient from recipe
func (coll *Collection) RemoveIngredient(ctx context.Context, recipeID string, stepIdx int, ingredientID string) error {

	recipe, methodStep, err := coll.getMethodStepByID(ctx, recipeID, stepIdx)
	if err != nil {
		fmt.Printf("Error retreiving method step, %v.\n", err)
		return err
	}

	filterFn := func(id *Ingredient) bool {
		return id.IngredientID != ingredientID
	}

	methodStep.Ingredients = filterIngredients(methodStep.Ingredients, filterFn)
	recipe.Method[stepIdx] = *methodStep

	return coll.collection.UpdateByID(ctx, recipeID, recipe)
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

func updateMethodStep(methodStep *MethodStep, updates map[string]string) *MethodStep {

	if update, ok := updates["action"]; ok {
		methodStep.Action = update
	}
	if update, ok := updates["description"]; ok {
		methodStep.Description = update
	}
	if update, ok := updates["time"]; ok {
		methodStep.Time = update
	}

	return methodStep
}

func (coll *Collection) getMethodStepByID(ctx context.Context, recipeID string, stepIdx int) (*Recipe, *MethodStep, error) {

	if stepIdx < 0 {
		return nil, nil, errors.New("Invalid index, " + strconv.Itoa(stepIdx))
	}

	recipe, err := coll.FindOne(ctx, recipeID)
	if err != nil {
		return nil, nil, err
	}

	if len(recipe.Method) <= stepIdx {
		return nil, nil, errors.New("Invalid index, " + strconv.Itoa(stepIdx))
	}

	methodStep := recipe.Method[stepIdx]

	return recipe, &methodStep, nil
}
