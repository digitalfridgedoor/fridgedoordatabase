package userview

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

func TestUpdate(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")

	duration3s, _ := time.ParseDuration("10s")
	ctx, cancelFunc := context.WithTimeout(context.Background(), duration3s)
	defer cancelFunc()
	connect := fridgedoordatabase.Connect(ctx, connectionstring)
	assert.True(t, connect)

	username := "TestUser"

	view, err := Create(context.Background(), username)

	assert.Nil(t, err)

	assert.Equal(t, username, view.Username)
	assert.NotNil(t, view.Collections)

	viewID := view.ID
	recipeCollectionName := "MyNewCollection"
	recipeID, _ := primitive.ObjectIDFromHex("5d8f7300a7888700270f7752")
	err = AddRecipe(context.Background(), view.ID.Hex(), recipeCollectionName, recipeID)
	assert.Nil(t, err)

	view, err = GetByUsername(context.Background(), username)
	assert.Nil(t, err)
	assert.Equal(t, viewID, view.ID)

	recipeCollection, ok := view.Collections[recipeCollectionName]
	assert.True(t, ok)
	assert.Equal(t, recipeCollectionName, recipeCollection.Name)
	assert.Equal(t, 1, len(recipeCollection.Recipes))

	recipe := recipeCollection.Recipes[0]
	assert.Equal(t, recipeID, recipe)

	err = RemoveRecipe(context.Background(), viewID.Hex(), recipeCollectionName, recipeID)
	assert.Nil(t, err)

	err = Delete(ctx, username)
	assert.Nil(t, err)

	_, err = GetByUsername(context.Background(), username)
	assert.NotNil(t, err)
}

func TestUsernameCanOnlyBeUsedOnce(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")

	duration3s, _ := time.ParseDuration("10s")
	ctx, cancelFunc := context.WithTimeout(context.Background(), duration3s)
	defer cancelFunc()
	connect := fridgedoordatabase.Connect(ctx, connectionstring)
	assert.True(t, connect)

	username := "TestUser"

	view, err := Create(context.Background(), username)
	assert.NotNil(t, view)
	assert.Nil(t, err)

	view, err = Create(context.Background(), username)
	assert.NotNil(t, err)
	assert.Equal(t, errUserExists, err)

	err = Delete(ctx, username)
	assert.Nil(t, err)
}
