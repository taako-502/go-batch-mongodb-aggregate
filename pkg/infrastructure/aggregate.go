package infrastructure

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserPoints struct {
	UserID     bson.ObjectID `bson:"_id"`
	TotalPoint int           `bson:"totalPoint"`
}

func AggregateUserPoints(client *mongo.Client, ctx context.Context) ([]UserPoints, error) {
	pointsCollection := client.Database("source").Collection("points")

	// 集計ステージの定義
	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$userId"},
		{Key: "totalPoint", Value: bson.D{{Key: "$sum", Value: "$point"}}},
	}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "totalPoint", Value: -1}}}}

	// 集計クエリを実行
	cursor, err := pointsCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, sortStage}, options.Aggregate())
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate points: %w", err)
	}
	defer cursor.Close(ctx)

	// 結果をスライスに格納
	var results []UserPoints
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to get all results: %w", err)
	}

	return results, nil
}
