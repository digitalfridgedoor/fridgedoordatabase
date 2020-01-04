package recipe

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

func TestUpdate(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	ctx := context.Background()
	ingredientID := "5d8f744446106c8aee8cde37"
	connected := fridgedoordatabase.Connect(ctx, connectionstring)
	assert.True(t, connected)

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	recipe, err := Create(ctx, userID, recipeName)

	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, recipeName, recipe.Name)

	recipeIDString := recipe.ID.Hex()

	err = AddMethodStep(ctx, userID, recipeIDString, "Add to pan")
	assert.Nil(t, err)

	err = AddIngredient(ctx, userID, recipeIDString, 0, ingredientID, "Test ing")
	assert.Nil(t, err)

	latestRecipe, err := FindOne(ctx, recipeIDString)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(latestRecipe.Method))
	method := latestRecipe.Method[0]
	assert.Equal(t, 1, len(method.Ingredients))

	updates := make(map[string]string)
	updates["amount"] = "1 1/2 tsp"
	err = UpdateIngredient(ctx, userID, recipeIDString, 0, ingredientID, updates)
	assert.Nil(t, err)

	latestRecipe, err = FindOne(ctx, recipeIDString)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(latestRecipe.Method))
	method = latestRecipe.Method[0]
	assert.Equal(t, 1, len(method.Ingredients))
	ing := method.Ingredients[0]
	assert.Equal(t, "1 1/2 tsp", ing.Amount)

	err = RemoveIngredient(ctx, userID, recipeIDString, 0, ingredientID)
	assert.Nil(t, err)

	err = RemoveMethodStepByIndex(ctx, userID, recipeIDString, 0)
	assert.Nil(t, err)

	Delete(ctx, recipe.ID)
}
