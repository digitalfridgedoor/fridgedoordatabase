package plan

import (
	"context"
	"testing"

	"github.com/digitalfridgedoor/fridgedoordatabase/dfdmodels"
	"github.com/digitalfridgedoor/fridgedoordatabase/dfdtesting"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
)

func TestFindByMonthAndYear(t *testing.T) {

	dfdtesting.SetTestCollectionOverride()
	dfdtesting.SetPlanFindPredicate(func(p *dfdmodels.Plan, m bson.M) bool {
		month := m["month"].(int)
		year := m["year"].(int)
		userid := m["userid"].(primitive.ObjectID)

		return month == p.Month && year == p.Year && userid == p.UserID
	})

	ok, coll := createCollection(context.TODO())
	assert.True(t, ok)

	userID := primitive.NewObjectID()

	r, err := coll.findByMonthAndYear(context.Background(), userID, 1, 2020)

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, 0, len(r))
}
