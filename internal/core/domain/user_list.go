package domain

import "time"

type UserQuery struct {
	Page                int64
	PerPage             int64
	Order               string
	RangeStartCreatedAt time.Time
	RangeEndCreatedAt   time.Time
}

type MetaUsers struct {
	Meta struct {
		Page    int64 `bson:"page" json:"page"`
		PerPage int64 `bson:"per_page" json:"per_page"`
	} `bson:"meta" json:"meta"`
	Users []User `json:"users"`
}
