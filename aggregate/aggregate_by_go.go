package aggregate

import (
	"context"
	"fmt"
	"log"

	"github.com/taako-502/go-batch-mongodb-aggregate/infrastructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: 順位も計算してデータベースに登録する
func AggregateByGo(ctx context.Context, client *mongo.Client, isPrint bool) {
	// MongoDBからすべてのpointsドキュメントを取得
	points := infrastructure.Find(ctx, client)

	userPoints := make(map[string]int)
	for _, p := range points {
		userPoints[p.UserID] += p.Point
	}

	// 結果を表示
	if isPrint {
		fmt.Println("")
		fmt.Println("Goで集計した結果")
		for userId, totalPoint := range userPoints {
			fmt.Printf("UserID: %v, TotalPoint: %v\n", userId, totalPoint)
		}
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
