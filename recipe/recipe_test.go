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
	r, err := connection.FindOne(context.Background(), "5dbc80036eb36874255e7fcd")

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "5dbc80036eb36874255e7fcd", r.ID.Hex())
	assert.Equal(t, "Nandos chicken", r.Name)
	assert.Equal(t, 2, len(r.Ingredients))
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
