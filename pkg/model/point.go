package model

import (
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Point ユーザーのポイントを格納する
type Point struct {
	ID     bson.ObjectID `bson:"_id"`
	UserID bson.ObjectID `bson:"userId"`
	Point  int           `bson:"point"`
}

func GeneratePoints(userIDs []bson.ObjectID, numberOfPoints int) []Point {
	var generatedPoints []Point
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for _, ID := range userIDs {
		for range numberOfPoints {
			generatedPoints = append(generatedPoints, Point{
				ID:     bson.NewObjectID(),
				UserID: ID,
				Point:  r.Intn(2000) + 1, // 1〜2000のランダムな値
			})
		}
	}
	return generatedPoints
}
