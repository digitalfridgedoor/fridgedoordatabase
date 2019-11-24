package recipe

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

func TestFindOne(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connect := fridgedoordatabase.Connect(context.Background(), connectionstring)

	connection := New(connect)
	ing, err := connection.FindOne(context.Background(), "5dbc80036eb36874255e7fcd")

	assert.Nil(t, err)
	assert.NotNil(t, ing)
	assert.Equal(t, "5dbc80036eb36874255e7fcd", ing.ID.Hex())
	assert.Equal(t, "Nandos chicken", ing.Name)
	assert.Equal(t, 2, len(ing.Ingredients))
	assert.Equal(t, 1, len(ing.Recipes))
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
