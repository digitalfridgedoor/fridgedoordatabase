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
	connect := fridgedoordatabase.Connect(context.Background(), connectionstring)

	connection := New(connect)
	r, err := connection.FindOne(context.Background(), "5dbc80036eb36874255e7fcd")

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "5dbc80036eb36874255e7fcd", r.ID.Hex())
	assert.Equal(t, "Nandos chicken", r.Name)
	assert.Equal(t, 0, len(r.Method))
	assert.Equal(t, 1, len(r.Recipes))
}

func TestList(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connect := fridgedoordatabase.Connect(context.Background(), connectionstring)

	connection := New(connect)
	recipes, err := connection.List(context.Background())

	assert.Nil(t, err)
	assert.NotNil(t, recipes)
	assert.LessOrEqual(t, len(recipes), 25)
}

func TestCreate(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connect := fridgedoordatabase.Connect(context.Background(), connectionstring)

	connection := New(connect)
	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	recipe, err := connection.Create(context.Background(), userID, recipeName)

	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, "new recipe", recipe.Name)

	connection.Delete(context.Background(), recipe.ID)
}

func TestAddAndRemove(t *testing.T) {
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

	err = connection.RemoveIngredient(ctx, recipeIDString, 0, ingredientID)
	assert.Nil(t, err)

	err = connection.RemoveMethodStepByIndex(ctx, recipeIDString, 0)
	assert.Nil(t, err)

	connection.Delete(ctx, recipe.ID)
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
