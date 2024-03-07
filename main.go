package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/taako-502/go-batch-mongodb-aggregate/aggregate"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type pattern struct {
	numberOfUsers  int
	numberOfPoints int
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDBに接続
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// 計測パターン
	patterns := []pattern{
		{numberOfUsers: 1000, numberOfPoints: 100},
		{numberOfUsers: 1000, numberOfPoints: 1000},
		{numberOfUsers: 1000, numberOfPoints: 10000},
		{numberOfUsers: 10000, numberOfPoints: 100},
		{numberOfUsers: 10000, numberOfPoints: 1000},
		{numberOfUsers: 10000, numberOfPoints: 10000},
		{numberOfUsers: 100000, numberOfPoints: 100},
		{numberOfUsers: 100000, numberOfPoints: 1000},
		{numberOfUsers: 100000, numberOfPoints: 10000},
		{numberOfUsers: 1000000, numberOfPoints: 100},
		{numberOfUsers: 1000000, numberOfPoints: 1000},
		{numberOfUsers: 1000000, numberOfPoints: 10000},
		// 極端に増やしてみる
		{numberOfUsers: 100000000, numberOfPoints: 100000},
	}

	for _, p := range patterns {
		// コレクションを初期化
		mr := mainReciever{ctx, client}
		mr.cleanUp("source", "users")
		mr.cleanUp("source", "points")
		mr.cleanUp("aggregate", "leaderboard")

		// ユーザーデータを挿入
		users := generateUsers(p.numberOfUsers)
		client.Database("source").Collection("users").InsertMany(ctx, users)

		// ポイントデータを挿入
		var userIDs []primitive.ObjectID
		for _, item := range users {
			if p, ok := item.(points); ok {
				userIDs = append(userIDs, p.userID)
			}
		}
		client.Database("source").Collection("points").InsertMany(ctx, generatePoints(userIDs, p.numberOfPoints))

		// MongoDBで集計
		startTime := time.Now()
		aggregate.AggregateByMongoDB(ctx, client)
		elapsed := time.Since(startTime)
		fmt.Printf("Method: MongoDB	ユーザ数: %d 1ユーザあたりのポイント数: %d	集計にかかった時間: %s\n", p.numberOfUsers, p.numberOfPoints, elapsed)

		// Goで集計
		startTime = time.Now()
		aggregate.AggregateByGo(ctx, client)
		elapsed = time.Since(startTime)
		fmt.Printf("Method: Go	ユーザ数: %d 1ユーザあたりのポイント数: %d	集計にかかった時間: %s\n", p.numberOfUsers, p.numberOfPoints, elapsed)
	}
}

type mainReciever struct {
	ctx    context.Context
	client *mongo.Client
}

func (m mainReciever) cleanUp(databaseName, collectionName string) {
	collection := m.client.Database(databaseName).Collection(collectionName)
	_, err := collection.DeleteMany(m.ctx, map[string]interface{}{})
	if err != nil {
		log.Fatal(err)
	}
}

type users struct {
	id   primitive.ObjectID `bson:"id"`
	name string             `bson:"name"`
}

func generateUsers(numberOfUsers int) []interface{} {
	var generatedUsers []interface{}
	for i := range numberOfUsers {
		user := users{id: primitive.NewObjectID(), name: fmt.Sprintf("user%v", i)}
		generatedUsers = append(generatedUsers, user)
	}
	return generatedUsers
}

type points struct {
	id     primitive.ObjectID `bson:"id"`
	userID primitive.ObjectID `bson:"userId"`
	point  int                `bson:"point"`
}

func generatePoints(userIDs []primitive.ObjectID, numberOfPoints int) []interface{} {
	var generatedPoints []interface{}
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for _, ID := range userIDs {
		for range numberOfPoints {
			generatedPoints = append(generatedPoints, points{id: primitive.NewObjectID(), userID: ID, point: r.Intn(2000) + 1}) // 1〜2000のランダムな値
		}
	}
	return generatedPoints
}
