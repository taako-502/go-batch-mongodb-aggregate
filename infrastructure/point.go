package infrastructure

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Point struct {
	ID     bson.ObjectID `bson:"_id"`
	UserID bson.ObjectID `bson:"userId"`
	Point  int           `bson:"point"`
}

func Find(ctx context.Context, client *mongo.Client) []Point {
	var points []Point
	pointsCollection := client.Database("source").Collection("points")
	cursor, err := pointsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &points); err != nil {
		log.Fatal(err)
	}
	return points
}
