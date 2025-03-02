package model

import (
	"math/rand/v2"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Point ユーザーのポイントを格納する
type Point struct {
	ID     bson.ObjectID `bson:"_id,omitempty"`
	UserID bson.ObjectID `bson:"userId"`
	Point  int           `bson:"point"`
}

func GeneratePoints(userIDs []bson.ObjectID, numberOfPoints int) []Point {
	var generatedPoints []Point
	for _, ID := range userIDs {
		for range numberOfPoints {
			generatedPoints = append(generatedPoints, Point{
				UserID: ID,
				Point:  rand.IntN(2000) + 1, // 1〜2000のランダムな値
			})
		}
	}
	return generatedPoints
}
