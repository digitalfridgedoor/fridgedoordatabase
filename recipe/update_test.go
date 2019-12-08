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
	connect := fridgedoordatabase.Connect(ctx, connectionstring)
	ingredientID := "5d8f744446106c8aee8cde37"

	connection := New(connect)
	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	recipe, err := connection.Create(ctx, userID, recipeName)

	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, "new recipe", recipe.Name)

	recipeIDString := recipe.ID.Hex()

	err = connection.AddMethodStep(ctx, recipeIDString, "Add to pan")
	assert.Nil(t, err)

	err = connection.AddIngredient(ctx, recipeIDString, 0, ingredientID, "Test ing")
	assert.Nil(t, err)

	latestRecipe, err := connection.FindOne(ctx, recipeIDString)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(latestRecipe.Method))
	method := latestRecipe.Method[0]
	assert.Equal(t, 1, len(method.Ingredients))

	updates := make(map[string]string)
	updates["amount"] = "1 1/2 tsp"
	err = connection.UpdateIngredient(ctx, recipeIDString, 0, ingredientID, updates)
	assert.Nil(t, err)

	latestRecipe, err = connection.FindOne(ctx, recipeIDString)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(latestRecipe.Method))
	method = latestRecipe.Method[0]
	assert.Equal(t, 1, len(method.Ingredients))
	ing := method.Ingredients[0]
	assert.Equal(t, "1 1/2 tsp", ing.Amount)

	err = connection.RemoveIngredient(ctx, recipeIDString, 0, ingredientID)
	assert.Nil(t, err)

	err = connection.RemoveMethodStepByIndex(ctx, recipeIDString, 0)
	assert.Nil(t, err)

	connection.Delete(ctx, recipe.ID)
}

func TestRunCode(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	ctx := context.Background()
	connection := fridgedoordatabase.Connect(ctx, connectionstring)
	recipeID := "5debadc725fbf484aed19ce4"

	collection := New(connection)
	recipe, err := collection.FindOne(ctx, recipeID)

	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, "Fajitas", recipe.Name)

	updates := make(map[string]string)
	updates["preperation"] = "test"
	collection.UpdateIngredient(ctx, recipeID, 0, "5d8f739ba7888700270f775a", updates)
}
