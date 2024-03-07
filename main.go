package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/taako-502/go-batch-mongodb-aggregate/aggregate"
	"github.com/taako-502/go-batch-mongodb-aggregate/infrastructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// usersをコンソールに表示
	mr := mainReciever{ctx, client}
	mr.printSourceUsers()

	// pointsをコンソールに表示
	mr.printSourcePoints()

	// MongoDBで集計
	startTime := time.Now()
	aggregate.AggregateByMongoDB(ctx, client, true)
	monogDBElapsed := time.Since(startTime)

	// Goで集計
	startTime = time.Now()
	aggregate.AggregateByGo(ctx, client, true)
	goElapsed := time.Since(startTime)

	fmt.Println("") // 改行
	fmt.Printf("Method: MongoDB\t集計にかかった時間: %s\n", monogDBElapsed)
	fmt.Printf("Method: Go\t集計にかかった時間: %s\n", goElapsed)
}

type mainReciever struct {
	ctx    context.Context
	client *mongo.Client
}

type user struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

func (m mainReciever) printSourceUsers() {
	// usersコレクションからデータを取得
	var users []user
	userCollection := m.client.Database("source").Collection("users")
	cursor, err := userCollection.Find(m.ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(m.ctx, &users); err != nil {
		log.Fatal(err)
	}

	// 取得したusersのデータを出力
	fmt.Println("# users:")
	for _, user := range users {
		fmt.Printf("ID: %s\tName: %s\n", user.ID.Hex(), user.Name)
	}
}

func (m mainReciever) printSourcePoints() {
	// pointsコレクションからデータを取得
	points := infrastructure.Find(m.ctx, m.client)

	// 取得したpointsのデータを出力
	fmt.Println("# points:")
	for _, point := range points {
		fmt.Printf("ID: %s\tUserID: %s\tPoint: %d\n", point.ID.Hex(), point.UserID, point.Point)
	}
}
