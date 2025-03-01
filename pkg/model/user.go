package model

import "go.mongodb.org/mongo-driver/v2/bson"

// User ユーザーデータを表す構造体
type User struct {
	ID   bson.ObjectID `bson:"_id"`
	Name string        `bson:"name"`
}
