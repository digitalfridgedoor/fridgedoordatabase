package recipe

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func canEdit(recipe *Recipe, userID primitive.ObjectID) bool {
	return recipe.AddedBy == userID
}

func appendString(current []string, value string) []string {
	hasValue := false

	for _, v := range current {
		if v == value {
			hasValue = true
		}
	}

	if !hasValue {
		current = append(current, value)
	}

	return current
}

func removeString(current []string, removeValue string) []string {
	filtered := []string{}

	for _, v := range current {
		if v != removeValue {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

func appendID(current []primitive.ObjectID, value primitive.ObjectID) []primitive.ObjectID {
	hasValue := false

	for _, v := range current {
		if v == value {
			hasValue = true
		}
	}

	if !hasValue {
		current = append(current, value)
	}

	return current
}

func removeID(current []primitive.ObjectID, removeValue primitive.ObjectID) []primitive.ObjectID {
	filtered := []primitive.ObjectID{}

	for _, v := range current {
		if v != removeValue {
			filtered = append(filtered, v)
		}
	}

	return filtered
}
