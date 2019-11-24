package ingredient

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"digitalfridgedoor/fridgedoordatabase"
)

func TestFind(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connect := fridgedoordatabase.Connect(context.Background(), connectionstring)

	connection := Connection{connect}
	ings, err := connection.Find(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, 2, len(ings))
}

func TestFindOne(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connect := fridgedoordatabase.Connect(context.Background(), connectionstring)

	connection := Connection{connect}
	ing, err := connection.FindOne(context.Background(), "5d8f744446106c8aee8cde37")

	assert.Nil(t, err)
	assert.NotNil(t, ing)
	assert.Equal(t, "5dac764fa0b9423b0090a898", ing.ParentID.Hex())
	assert.Equal(t, "5d8f744446106c8aee8cde37", ing.ID.Hex())
	assert.Equal(t, "Chicken thighs", ing.Name)
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
