package userview

import (
	"context"
	"testing"
	"time"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
	"github.com/digitalfridgedoor/fridgedoordatabase/dfdtesting"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUsernameCanOnlyBeUsedOnce(t *testing.T) {

	dfdtesting.SetTestCollectionOverride()
	dfdtesting.SetUserViewFindPredicate(func(uv *dfdmodels.UserView, m primitive.M) bool {
		return m["username"] == uv.Username
	})

	duration, _ := time.ParseDuration("10s")
	ctx, cancelFunc := context.WithTimeout(context.Background(), duration)
	defer cancelFunc()

	username := "TestUser"

	view, err := Create(ctx, username)
	assert.NotNil(t, view)
	assert.Nil(t, err)

	view, err = Create(ctx, username)
	assert.NotNil(t, err)
	assert.Equal(t, errUserExists, err)

	err = Delete(ctx, username)
	assert.Nil(t, err)
}
