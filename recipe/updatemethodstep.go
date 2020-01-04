package recipe

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddMethodStep adds new method step to a recipe
func AddMethodStep(ctx context.Context, user primitive.ObjectID, recipeID string, action string) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	recipe, err := FindOne(ctx, recipeID)
	if err != nil {
		return err
	}

	methodStep := MethodStep{
		Action: action,
	}

	recipe.Method = append(recipe.Method, methodStep)

	return collection.UpdateByID(ctx, recipeID, recipe)
}

// UpdateMethodStepByIndex updates method step at index
func UpdateMethodStepByIndex(ctx context.Context, user primitive.ObjectID, recipeID string, stepIdx int, updates map[string]string) error {

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

	recipe.Method[stepIdx] = *updateMethodStep(methodStep, updates)

	return collection.UpdateByID(ctx, recipeID, recipe)
}

// RemoveMethodStepByIndex removes method by index
func RemoveMethodStepByIndex(ctx context.Context, user primitive.ObjectID, recipeID string, stepIdx int) error {

	connected, collection := collection()
	if !connected {
		return errNotConnected
	}

	if stepIdx < 0 {
		return errors.New("Invalid index")
	}

	recipe, err := FindOne(ctx, recipeID)
	if err != nil {
		return err
	}

	if !canEdit(recipe, user) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	if len(recipe.Method) <= stepIdx {
		return errors.New("Invalid index")
	}

	copy(recipe.Method[stepIdx:], recipe.Method[stepIdx+1:]) // Shift a[i+1:] left one index.
	recipe.Method = recipe.Method[:len(recipe.Method)-1]     // Truncate slice.

	return collection.UpdateByID(ctx, recipeID, recipe)
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

func getMethodStepByID(ctx context.Context, recipeID string, stepIdx int) (*Recipe, *MethodStep, error) {

	if stepIdx < 0 {
		return nil, nil, errors.New("Invalid index, " + strconv.Itoa(stepIdx))
	}

	recipe, err := FindOne(ctx, recipeID)
	if err != nil {
		return nil, nil, err
	}

	if len(recipe.Method) <= stepIdx {
		return nil, nil, errors.New("Invalid index, " + strconv.Itoa(stepIdx))
	}

	methodStep := recipe.Method[stepIdx]

	return recipe, &methodStep, nil
}
