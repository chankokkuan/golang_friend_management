package util

import (
	"friend-management/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson"
)

func GenerateSortFilter(query domain.UserQuery) bson.D {
	var filter bson.D
	if query.Order == "asc" {
		filter = bson.D{
			{Key: "created_at", Value: 1},
			{Key: "_id", Value: 1},
		}
	} else {
		filter = bson.D{
			{Key: "created_at", Value: -1},
			{Key: "_id", Value: -1},
		}
	}
	return filter
}
