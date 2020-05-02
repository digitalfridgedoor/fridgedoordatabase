package recipe

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {

	ctx := context.Background()

	ingredientID := "5d8f744446106c8aee8cde37"
	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	recipe, err := Create(ctx, userID, recipeName)

	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, recipeName, recipe.Name)

	err = AddMethodStep(ctx, userID, recipe.ID, "Add to pan")
	assert.Nil(t, err)

	latestRecipe, err := AddIngredient(ctx, userID, recipe.ID, 0, ingredientID, "Test ing")
	assert.Nil(t, err)
	method := latestRecipe.Method[0]
	assert.Equal(t, 1, len(method.Ingredients))

	updates := make(map[string]string)
	updates["amount"] = "1 1/2 tsp"

	latestRecipe, err = UpdateIngredient(ctx, userID, recipe.ID, 0, ingredientID, updates)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(latestRecipe.Method))
	method = latestRecipe.Method[0]
	assert.Equal(t, 1, len(method.Ingredients))
	ing := method.Ingredients[0]
	assert.Equal(t, "1 1/2 tsp", ing.Amount)

	latestRecipe, err = RemoveIngredient(ctx, userID, recipe.ID, 0, ingredientID)
	assert.Nil(t, err)

	err = RemoveMethodStepByIndex(ctx, userID, recipe.ID, 0)
	assert.Nil(t, err)

	Delete(ctx, recipe.ID)
}
