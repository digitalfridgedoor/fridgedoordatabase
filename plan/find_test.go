package plan

import (
	"context"
	"testing"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdtesting"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
)

func TestFindByMonthAndYear(t *testing.T) {

	dfdtesting.SetTestCollectionOverride()
	dfdtesting.SetPlanFindPredicate(dfdtesting.FindPlanByMonthAndYearTestPredicate)

	ok, coll := createCollection(context.TODO())
	assert.True(t, ok)

	userID := primitive.NewObjectID()

	r, err := coll.findByMonthAndYear(context.Background(), userID, 1, 2020)

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, 0, len(r))
}
