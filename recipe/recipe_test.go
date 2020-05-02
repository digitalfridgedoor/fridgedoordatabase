package recipe

import (
	"context"
	"os"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
)

func TestFindOne(t *testing.T) {

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	assert.Nil(t, err)
	id, err := primitive.ObjectIDFromHex("5dbc814c6eb36874255e7fd0")
	assert.Nil(t, err)

	r, err := FindOne(context.Background(), &id, userID)

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "5dbc814c6eb36874255e7fd0", r.ID.Hex())
	assert.Equal(t, "Macho peas", r.Name)
	assert.Equal(t, 0, len(r.Method))
	assert.Equal(t, 0, len(r.Recipes))
}

func TestFindStartingWith(t *testing.T) {

	userID, err := primitive.ObjectIDFromHex("5de28cfd7633c82c6982cd0a")
	assert.Nil(t, err)

	results, err := FindByName(context.Background(), "fi", userID, 10)

	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 1, len(results))
}

func TestFindByTags(t *testing.T) {

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")

	recipes, err := FindByTags(context.Background(), userID, []string{}, []string{}, 5)

	assert.Nil(t, err)
	assert.NotNil(t, recipes)
	assert.LessOrEqual(t, len(recipes), 5)
}

func TestCreate(t *testing.T) {

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	recipe, err := Create(context.Background(), userID, recipeName)

	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, "new recipe", recipe.Name)

	Delete(context.Background(), recipe.ID)
}

func TestAddAndRemove(t *testing.T) {
	ctx := context.Background()
	ingredientID := "5d8f744446106c8aee8cde37"

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	recipe, err := Create(ctx, userID, recipeName)

	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, "new recipe", recipe.Name)

	err = AddMethodStep(ctx, userID, recipe.ID, "Add to pan")
	assert.Nil(t, err)

	err = AddIngredient(ctx, userID, recipe.ID, 0, ingredientID, "Test ing")
	assert.Nil(t, err)

	latestRecipe, err := FindOne(ctx, recipe.ID, userID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(latestRecipe.Method))
	method := latestRecipe.Method[0]
	assert.Equal(t, 1, len(method.Ingredients))

	err = RemoveIngredient(ctx, userID, recipe.ID, 0, ingredientID)
	assert.Nil(t, err)

	err = RemoveMethodStepByIndex(ctx, userID, recipe.ID, 0)
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
