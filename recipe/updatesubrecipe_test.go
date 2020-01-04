package recipe

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

func TestAddSubRecipe(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	ctx := context.Background()
	connected := fridgedoordatabase.Connect(ctx, connectionstring)
	assert.True(t, connected)

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	subRecipeName := "new sub recipe"
	recipe, err := Create(ctx, userID, recipeName)
	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, recipeName, recipe.Name)

	subRecipe, err := Create(ctx, userID, subRecipeName)
	assert.Nil(t, err)
	assert.NotNil(t, subRecipe)
	assert.Equal(t, subRecipeName, subRecipe.Name)

	recipeIDString := recipe.ID.Hex()
	subRecipeIDString := subRecipe.ID.Hex()

	err = AddSubRecipe(ctx, userID, recipeIDString, subRecipeIDString)
	assert.Nil(t, err)

	latestRecipe, err := FindOne(ctx, recipeIDString)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(latestRecipe.Recipes))
	latestSubRecipe := latestRecipe.Recipes[0]

	assert.Equal(t, subRecipe.ID, latestSubRecipe.RecipeID)
	assert.Equal(t, subRecipe.Name, subRecipeName)

	err = RemoveSubRecipe(ctx, userID, recipeIDString, subRecipeIDString)
	assert.Nil(t, err)

	latestRecipe, err = FindOne(ctx, recipeIDString)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(latestRecipe.Recipes))

	Delete(ctx, recipe.ID)
	Delete(ctx, subRecipe.ID)
}
