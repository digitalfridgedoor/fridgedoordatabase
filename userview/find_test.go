package userview

import (
	"context"
	"testing"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
	"github.com/digitalfridgedoor/fridgedoordatabase/dfdtesting"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
)

func TestFindByUsername(t *testing.T) {

	dfdtesting.SetTestCollectionOverride()
	dfdtesting.SetUserViewFindPredicate(func(uv *dfdmodels.UserView, m primitive.M) bool {
		return m["username"] == uv.Username
	})

	r, err := GetByUsername(context.Background(), "Maisie")

	assert.NotNil(t, err)
	assert.Nil(t, r)
}

func TestFindLinked(t *testing.T) {

	dfdtesting.SetTestCollectionOverride()
	dfdtesting.SetUserViewFindPredicate(func(uv *dfdmodels.UserView, m primitive.M) bool {
		return true
	})

	username := "TestUser"
	Create(context.TODO(), username)

	r, err := GetLinkedUserViews(context.Background())

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Greater(t, len(r), 0)

	err = Delete(context.TODO(), username)
	assert.Nil(t, err)
}
