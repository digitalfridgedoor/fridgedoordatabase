package recipe

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddMethodStep adds new method step to a recipe
func AddMethodStep(ctx context.Context, userID primitive.ObjectID, recipeID *primitive.ObjectID, action string) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	recipe, err := coll.findOne(ctx, recipeID, userID)
	if err != nil {
		return err
	}

	methodStep := dfdmodels.MethodStep{
		Action: action,
	}

	recipe.Method = append(recipe.Method, methodStep)

	return coll.c.UpdateByID(ctx, recipeID, recipe)
}

// UpdateMethodStepByIndex updates method step at index
func UpdateMethodStepByIndex(ctx context.Context, userID primitive.ObjectID, recipeID *primitive.ObjectID, stepIdx int, updates map[string]string) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	recipe, methodStep, err := coll.getMethodStepByID(ctx, recipeID, userID, stepIdx)
	if err != nil {
		fmt.Printf("Error retreiving method step, %v.\n", err)
		return err
	}

	if !CanEdit(recipe, userID) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	recipe.Method[stepIdx] = *updateMethodStep(methodStep, updates)

	return coll.c.UpdateByID(ctx, recipeID, recipe)
}

// RemoveMethodStepByIndex removes method by index
func RemoveMethodStepByIndex(ctx context.Context, userID primitive.ObjectID, recipeID *primitive.ObjectID, stepIdx int) error {

	ok, coll := createCollection(ctx)
	if !ok {
		fmt.Println("Not connected")
		return errNotConnected
	}

	if stepIdx < 0 {
		return errors.New("Invalid index")
	}

	recipe, err := coll.findOne(ctx, recipeID, userID)
	if err != nil {
		return err
	}

	if !CanEdit(recipe, userID) {
		fmt.Println("User not authorised to update recipe")
		return errUnauthorised
	}

	if len(recipe.Method) <= stepIdx {
		return errors.New("Invalid index")
	}

	copy(recipe.Method[stepIdx:], recipe.Method[stepIdx+1:]) // Shift a[i+1:] left one index.
	recipe.Method = recipe.Method[:len(recipe.Method)-1]     // Truncate slice.

	return coll.c.UpdateByID(ctx, recipeID, recipe)
}

func updateMethodStep(methodStep *dfdmodels.MethodStep, updates map[string]string) *dfdmodels.MethodStep {

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

func (coll *collection) getMethodStepByID(ctx context.Context, recipeID *primitive.ObjectID, userID primitive.ObjectID, stepIdx int) (*dfdmodels.Recipe, *dfdmodels.MethodStep, error) {

	if stepIdx < 0 {
		return nil, nil, errors.New("Invalid index, " + strconv.Itoa(stepIdx))
	}

	recipe, err := coll.findOne(ctx, recipeID, userID)
	if err != nil {
		return nil, nil, err
	}

	if len(recipe.Method) <= stepIdx {
		return nil, nil, errors.New("Invalid index, " + strconv.Itoa(stepIdx))
	}

	methodStep := recipe.Method[stepIdx]

	return recipe, &methodStep, nil
}
