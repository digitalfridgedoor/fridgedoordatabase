package recipe

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
)

func TestAddSubRecipe(t *testing.T) {

	ctx := context.Background()

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

	latestRecipe, err := AddSubRecipe(ctx, userID, recipe.ID, subRecipe.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(latestRecipe.Recipes))
	latestSubRecipe := latestRecipe.Recipes[0]

	assert.Equal(t, *subRecipe.ID, latestSubRecipe.RecipeID)
	assert.Equal(t, subRecipe.Name, subRecipeName)

	// Check actual sub recipe
	latestSubRecipeMain, err := FindOne(ctx, subRecipe.ID, userID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(latestSubRecipeMain.ParentIds))

	latestRecipe, err = RemoveSubRecipe(ctx, userID, recipe.ID, subRecipe.ID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(latestRecipe.Recipes))

	// Check actual sub recipe
	latestSubRecipeMain, err = FindOne(ctx, subRecipe.ID, userID)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(latestSubRecipeMain.ParentIds))

	Delete(ctx, recipe.ID)
	Delete(ctx, subRecipe.ID)
}
