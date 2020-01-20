package plan

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateDayPlanRequest is the request for updating a day plan
type UpdateDayPlanRequest struct {
	UserID     primitive.ObjectID
	Month      int
	Year       int
	Day        int
	MealIndex  int
	RecipeName string
	RecipeID   primitive.ObjectID
}

// Update updates a Plan for a user
func Update(ctx context.Context, updateRequest *UpdateDayPlanRequest) (*primitive.ObjectID, error) {

	connected, collection := collection()
	if !connected {
		return nil, errNotConnected
	}

	plan, isNew, err := FindOrCreateOne(ctx, updateRequest.UserID, updateRequest.Month, updateRequest.Year)
	if err != nil {
		return nil, err
	}

	if len(plan.Days) < updateRequest.Day {
		fmt.Printf("Invalid day (%v) for month with %v days.\n", updateRequest.Day, len(plan.Days))
		return nil, errInvalidInput
	}

	currentPlanLength := len(plan.Days[updateRequest.Day-1].Meal)

	if currentPlanLength == 0 {
		plan.Days[updateRequest.Day-1].Meal = make([]Meal, updateRequest.MealIndex+1)
	} else if currentPlanLength <= updateRequest.MealIndex {
		diff := updateRequest.MealIndex + 1 - currentPlanLength
		plan.Days[updateRequest.Day-1].Meal = append(plan.Days[updateRequest.Day-1].Meal, make([]Meal, diff)...)
	}

	plan.Days[updateRequest.Day-1].Meal[updateRequest.MealIndex].Name = updateRequest.RecipeName
	plan.Days[updateRequest.Day-1].Meal[updateRequest.MealIndex].RecipeID = updateRequest.RecipeID

	if isNew {
		return collection.InsertOne(ctx, plan)
	}

	err = collection.UpdateByID(ctx, plan.ID.Hex(), plan)
	return &plan.ID, err
}

func create(userID primitive.ObjectID, month int, year int) (bool, *Plan) {
	ok, dayLength := days(month, year)
	if !ok {
		return false, nil
	}

	days := make([]Day, dayLength)
	return true, &Plan{
		UserID: userID,
		Month:  month,
		Year:   year,
		Days:   days,
	}
}

func days(month int, year int) (bool, int) {
	if month > 12 || month < 1 {
		return false, 0
	}
	if year < 2000 {
		return false, 0
	}

	switch month {
	case 4:
		return true, 30
	case 6:
		return true, 30
	case 9:
		return true, 30
	case 11:
		return true, 30
	case 2:
		if year%4 == 0 {
			return true, 29
		}
		return true, 28
	default:
		return true, 31
	}
}
