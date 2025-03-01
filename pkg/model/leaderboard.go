package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Leaderboard struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	UserID     bson.ObjectID `bson:"userId"`
	Method     string        `bson:"method"`
	TotalPoint int           `bson:"totalPoint"`
	Rank       int           `bson:"rank"`
	CreatedAt  bson.DateTime `bson:"createdAt"`
	UpdatedAt  bson.DateTime `bson:"updatedAt"`
}
