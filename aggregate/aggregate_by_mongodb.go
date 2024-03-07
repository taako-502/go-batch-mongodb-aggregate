package aggregate

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: 順位も計算してデータベースに登録する
func AggregateByMongoDB(ctx context.Context, client *mongo.Client, isPrint bool) {
	// pointsコレクションからデータを集計
	pointsCollection := client.Database("source").Collection("points")
	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$userId"}, {Key: "totalPoint", Value: bson.D{{Key: "$sum", Value: "$point"}}}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "totalPoint", Value: -1}}}}

	cursor, err := pointsCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, sortStage})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	// 結果を表示
	if isPrint {
		fmt.Println("MongoDBの集計関数で集計した結果")
		for _, totalPoint := range results {
			ID := strings.Replace(fmt.Sprintf("%v", totalPoint["_id"]), "ObjectID(\"", "", -1)
			ID = strings.Replace(ID, "\")", "", -1)
			fmt.Printf("UserID: %v, TotalPoint: %v\n", ID, totalPoint["totalPoint"])
		}
	}

	// rankingデータベースのrankingテーブルを更新
	rankingCollection := client.Database("aggregate").Collection("leaderboard")
	for _, result := range results {
		filter := bson.M{"userId": result["_id"]}
		update := bson.M{"$set": bson.M{"totalPoint": result["totalPoint"]}}
		options := options.Update().SetUpsert(true)
		_, err := rankingCollection.UpdateOne(ctx, filter, update, options)
		if err != nil {
			log.Fatal(err)
		}
	}
}
