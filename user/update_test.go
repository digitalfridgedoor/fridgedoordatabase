package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

func TestUpdate(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connect := fridgedoordatabase.Connect(context.Background(), connectionstring)

	connection := New(connect)
	userID := "5d8f7300a7888700270f7752"
	recipeID := primitive.NewObjectID()
	err := connection.AddRecipe(context.Background(), userID, recipeID)

	assert.Nil(t, err)

	err = connection.RemoveRecipe(context.Background(), userID, recipeID)

	assert.Nil(t, err)
}
