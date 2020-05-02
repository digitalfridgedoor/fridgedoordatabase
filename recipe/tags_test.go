package recipe

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

func TestTags(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	ctx := context.Background()
	connected := fridgedoordatabase.Connect(ctx, connectionstring)
	assert.True(t, connected)

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	recipe, err := Create(ctx, userID, recipeName)

	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, recipeName, recipe.Name)

	recipeIDString := recipe.ID.Hex()
	tag := primitive.NewObjectID().Hex()

	updates := make(map[string]string)
	updates["tag_add"] = tag
	err = UpdateMetadata(ctx, userID, recipeIDString, updates)
	assert.Nil(t, err)

	results, err := FindByTags(ctx, userID, []string{tag}, []string{}, 20)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(results))

	updates = make(map[string]string)
	updates["tag_remove"] = tag
	err = UpdateMetadata(ctx, userID, recipeIDString, updates)

	results, err = FindByTags(ctx, userID, []string{tag}, []string{}, 20)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(results))

	Delete(ctx, recipe.ID)
}

func TestNinTags(t *testing.T) {

	connectionstring := getEnvironmentVariable("connectionstring")
	ctx := context.Background()
	connected := fridgedoordatabase.Connect(ctx, connectionstring)
	assert.True(t, connected)

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	recipe, err := Create(ctx, userID, recipeName)
	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, recipeName, recipe.Name)

	recipeIDString := recipe.ID.Hex()

	tag := primitive.NewObjectID().Hex()

	updates := make(map[string]string)
	updates["tag_add"] = tag
	err = UpdateMetadata(ctx, userID, recipeIDString, updates)
	assert.Nil(t, err)

	results, err := FindByTags(ctx, userID, []string{}, []string{tag}, 20)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(results), 1)

	recipeInResult := false
	for _, r := range results {
		if r.ID == recipe.ID {
			recipeInResult = true
		}
	}

	assert.False(t, recipeInResult)

	Delete(ctx, recipe.ID)
}

func TestIncludeAndNinTags(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	ctx := context.Background()
	connected := fridgedoordatabase.Connect(ctx, connectionstring)
	assert.True(t, connected)

	userID, err := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	recipeName := "new recipe"
	recipe, err := Create(ctx, userID, recipeName)
	assert.Nil(t, err)
	assert.NotNil(t, recipe)
	assert.Equal(t, recipeName, recipe.Name)

	recipeIDString := recipe.ID.Hex()

	tag := primitive.NewObjectID().Hex()
	anothertag := primitive.NewObjectID().Hex()

	updates := make(map[string]string)
	updates["tag_add"] = tag
	err = UpdateMetadata(ctx, userID, recipeIDString, updates)
	assert.Nil(t, err)

	r, err := FindOne(ctx, recipeIDString)
	assert.Nil(t, err)
	assert.NotNil(t, r)

	results, err := FindByTags(ctx, userID, []string{tag}, []string{anothertag}, 20)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, recipeName, results[0].Name)

	updates = make(map[string]string)
	updates["tag_add"] = anothertag
	err = UpdateMetadata(ctx, userID, recipeIDString, updates)
	assert.Nil(t, err)

	results, err = FindByTags(ctx, userID, []string{tag}, []string{anothertag}, 20)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(results))

	updates = make(map[string]string)
	updates["tag_remove"] = anothertag
	err = UpdateMetadata(ctx, userID, recipeIDString, updates)
	results, err = FindByTags(ctx, userID, []string{tag}, []string{anothertag}, 20)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, recipeName, results[0].Name)

	Delete(ctx, recipe.ID)
}
