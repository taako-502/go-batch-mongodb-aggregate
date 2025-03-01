package model

import "go.mongodb.org/mongo-driver/v2/bson"

type UserPoints struct {
	UserID     bson.ObjectID `bson:"_id"`
	TotalPoint int           `bson:"totalPoint"`
}
