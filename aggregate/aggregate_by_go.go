package aggregate

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: 順位も計算してデータベースに登録する
func AggregateByGo(ctx context.Context, client *mongo.Client) {
	pointsCollection := client.Database("source").Collection("points")

	// MongoDBからすべてのpointsドキュメントを取得
	cursor, err := pointsCollection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	userPoints := make(map[interface{}]int)
	for cursor.Next(ctx) {
		var point bson.M
		if err := cursor.Decode(&point); err != nil {
			log.Fatal(err)
		}

		userId := point["userId"]
		pointValue, ok := point["point"].(int)
		if !ok {
			continue
		}

		if _, exists := userPoints[userId]; !exists {
			userPoints[userId] = pointValue
		} else {
			userPoints[userId] += pointValue
		}
	}

	// 結果を表示
	for userId, totalPoint := range userPoints {
		fmt.Printf("UserID: %v, TotalPoint: %v\n", userId, totalPoint)
	}

	// rankingデータベースのrankingテーブルを更新
	rankingCollection := client.Database("aggregate").Collection("leaderboard")
	for userId, totalPoint := range userPoints {
		filter := bson.M{"userId": userId}
		update := bson.M{"$set": bson.M{"totalPoint": totalPoint}}
		upsert := true
		opts := options.UpdateOptions{Upsert: &upsert}
		_, err := rankingCollection.UpdateOne(ctx, filter, update, &opts)
		if err != nil {
			log.Fatal(err)
		}
	}
}
