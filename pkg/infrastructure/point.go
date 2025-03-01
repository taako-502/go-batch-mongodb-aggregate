package infrastructure

import (
	"context"
	"log"

	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Find(ctx context.Context, client *mongo.Client) []model.Point {
	var points []model.Point
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
