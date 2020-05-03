package plan

import (
	"context"

	"github.com/digitalfridgedoor/fridgedoordatabase/database"
)

type collection struct {
	c database.ICollection
}

func createCollection(ctx context.Context) (bool, *collection) {
	if ok, coll := database.Plan(ctx); ok {
		return true, &collection{coll}
	}
	return false, nil
}
