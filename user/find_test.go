package user

import (
	"context"
	"digitalfridgedoor/fridgedoordatabase"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindOne(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connect := fridgedoordatabase.Connect(context.Background(), connectionstring)

	connection := New(connect)
	r, err := connection.GetByUsername(context.Background(), "Maisie")

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "5d8f7300a7888700270f7752", r.ID.Hex())
	assert.Greater(t, len(r.Recipes), 0)
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
