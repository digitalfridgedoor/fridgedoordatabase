package userview

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/digitalfridgedoor/fridgedoordatabase"

	"github.com/stretchr/testify/assert"
)

func TestFindByUsername(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connect := fridgedoordatabase.Connect(context.Background(), connectionstring)
	assert.True(t, connect)

	r, err := GetByUsername(context.Background(), "Maisie")

	assert.NotNil(t, err)
	assert.Nil(t, r)
}

func TestFindLinked(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connect := fridgedoordatabase.Connect(context.Background(), connectionstring)
	assert.True(t, connect)

	r, err := GetLinkedUserViews(context.Background())

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Greater(t, len(r), 0)
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
