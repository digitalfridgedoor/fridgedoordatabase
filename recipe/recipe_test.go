package recipe

import (
	"context"
	"os"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

func TestFindOne(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connected := fridgedoordatabase.Connect(context.Background(), connectionstring)
	assert.True(t, connected)

	r, err := FindOne(context.Background(), "5dbc814c6eb36874255e7fd0")

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "5dbc814c6eb36874255e7fd0", r.ID.Hex())
	assert.Equal(t, "Macho peas", r.Name)
	assert.Equal(t, 0, len(r.Method))
	assert.Equal(t, 0, len(r.Recipes))
}

func TestFindStartingWith(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connected := fridgedoordatabase.Connect(context.Background(), connectionstring)
	assert.True(t, connected)

	userID, err := primitive.ObjectIDFromHex("5de28cfd7633c82c6982cd0a")
	assert.Nil(t, err)

	results, err := FindByName(context.Background(), "fi", userID)

	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 1, len(results))
}

func TestUserRecipes(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connected := fridgedoordatabase.Connect(context.Background(), connectionstring)
	assert.True(t, connected)

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")

	recipes, err := UserRecipes(context.Background(), userID)

	assert.Nil(t, err)
	assert.NotNil(t, recipes)
	assert.LessOrEqual(t, len(recipes), 25)
}

func TestCreate(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connected := fridgedoordatabase.Connect(context.Background(), connectionstring)
	assert.True(t, connected)

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	recipe, err := Create(context.Background(), userID, recipeName)

	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, "new recipe", recipe.Name)

	Delete(context.Background(), recipe.ID)
}

func TestAddAndRemove(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	ctx := context.Background()
	ingredientID := "5d8f744446106c8aee8cde37"
	connected := fridgedoordatabase.Connect(context.Background(), connectionstring)
	assert.True(t, connected)

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	recipe, err := Create(ctx, userID, recipeName)

	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, "new recipe", recipe.Name)

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

	err = RemoveIngredient(ctx, userID, recipeIDString, 0, ingredientID)
	assert.Nil(t, err)

	err = RemoveMethodStepByIndex(ctx, userID, recipeIDString, 0)
	assert.Nil(t, err)

	Delete(ctx, recipe.ID)
}

func getEnvironmentVariable(key string) string {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == key {
			return pair[1]
		}
	}

	os.Exit(1)
	return ""
}
