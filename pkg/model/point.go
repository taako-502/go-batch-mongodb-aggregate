package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Point struct {
	ID     bson.ObjectID `bson:"_id"`
	UserID bson.ObjectID `bson:"userId"`
	Point  int           `bson:"point"`
}
