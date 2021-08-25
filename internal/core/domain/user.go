package domain

import "time"

type Friend struct {
	ID    string `bson:"id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
}

type User struct {
	ID         string    `bson:"_id" json:"id"`
	Name       string    `bson:"name" json:"name"`
	Email      string    `bson:"email" json:"email"`
	Friends    []Friend  `bson:"friends" json:"friends"`
	VersionRev string    `bson:"version_rev" json:"version_rev"`
	VersionSeq int       `bson:"version_seq" json:"version_seq"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}
