package plan

import (
	"context"
	"os"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"

	"github.com/digitalfridgedoor/fridgedoordatabase"
)

func TestFindByMonthAndYear(t *testing.T) {
	connectionstring := getEnvironmentVariable("connectionstring")
	connected := fridgedoordatabase.Connect(context.Background(), connectionstring)
	assert.True(t, connected)

	userID := primitive.NewObjectID()

	r, err := findByMonthAndYear(context.Background(), userID, 1, 2020)

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, 0, len(r))
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
