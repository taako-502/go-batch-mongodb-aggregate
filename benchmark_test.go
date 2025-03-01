package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/benchmark"
	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/benchmark/aggregate"
	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

type pattern struct {
	numberOfUsers  int
	numberOfPoints int
}

var benchmarkPattern = []pattern{
	{numberOfUsers: 1000, numberOfPoints: 100},
	{numberOfUsers: 1000, numberOfPoints: 1000},
	// {numberOfUsers: 1000, numberOfPoints: 10000},
	// {numberOfUsers: 10000, numberOfPoints: 100},
	// {numberOfUsers: 10000, numberOfPoints: 1000},
	// {numberOfUsers: 10000, numberOfPoints: 10000},
	// {numberOfUsers: 100000, numberOfPoints: 100},
	// {numberOfUsers: 100000, numberOfPoints: 1000},
	// {numberOfUsers: 100000, numberOfPoints: 10000},
	// {numberOfUsers: 1000000, numberOfPoints: 100},
	// {numberOfUsers: 1000000, numberOfPoints: 1000},
	// {numberOfUsers: 1000000, numberOfPoints: 10000},
	// {numberOfUsers: 10000000, numberOfPoints: 10000},
	// {numberOfUsers: 10000000, numberOfPoints: 100000},
	// {numberOfUsers: 10000000, numberOfPoints: 1000000},
	// {numberOfUsers: 100000000, numberOfPoints: 1000000},
}

func init() {
	var err error
	client, err = mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
}

// BenchmarkAggregationPipeline Aggregation Pipelineで集計する
func BenchmarkAggregationPipeline(b *testing.B) {
	ctx := context.Background()
	// コレクションを初期化
	benchmark.Cleanup(ctx, client.Database("source").Collection("users"))
	benchmark.Cleanup(ctx, client.Database("source").Collection("points"))
	benchmark.Cleanup(ctx, client.Database("aggregate").Collection("leaderboard"))

	for _, n := range benchmarkPattern {
		b.Run("Benchmark_"+fmt.Sprint(n), func(b *testing.B) {
			// ユーザーデータを挿入
			users := generateUsers(n.numberOfUsers)
			if _, err := client.Database("source").Collection("users").InsertMany(ctx, users); err != nil {
				b.Fatal(err)
			}

			// ポイントデータを挿入
			var userIDs []bson.ObjectID
			for _, u := range users {
				userIDs = append(userIDs, u.ID)
			}
			points := generatePoints(userIDs, n.numberOfPoints)
			if _, err := client.Database("source").Collection("points").InsertMany(ctx, points); err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for b.Loop() {
				aggregate.AggregateByMongoDB(ctx, client, false)
			}
		})
	}
}

// BenchmarkGo Goで集計する
func BenchmarkGo(b *testing.B) {
	ctx := context.Background()
	benchmark.Cleanup(ctx, client.Database("source").Collection("users"))
	benchmark.Cleanup(ctx, client.Database("source").Collection("points"))
	benchmark.Cleanup(ctx, client.Database("aggregate").Collection("leaderboard"))

	for _, n := range benchmarkPattern {
		b.Run("Benchmark_"+fmt.Sprint(n), func(b *testing.B) {
			users := generateUsers(n.numberOfUsers)
			if _, err := client.Database("source").Collection("users").InsertMany(ctx, users); err != nil {
				b.Fatal(err)
			}

			var userIDs []bson.ObjectID
			for _, u := range users {
				userIDs = append(userIDs, u.ID)
			}
			points := generatePoints(userIDs, n.numberOfPoints)
			if _, err := client.Database("source").Collection("points").InsertMany(ctx, points); err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for b.Loop() {
				aggregate.AggregateByGo(ctx, client, false)
			}
		})
	}
}

func generateUsers(numberOfUsers int) []model.User {
	var users []model.User
	for i := range numberOfUsers {
		users = append(users, model.User{
			ID:   bson.NewObjectID(),
			Name: fmt.Sprintf("user%v", i),
		})
	}
	return users
}

func generatePoints(userIDs []bson.ObjectID, numberOfPoints int) []interface{} {
	var generatedPoints []interface{}
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for _, ID := range userIDs {
		for range numberOfPoints {
			generatedPoints = append(generatedPoints, model.Point{
				ID:     bson.NewObjectID(),
				UserID: ID,
				Point:  r.Intn(2000) + 1, // 1〜2000のランダムな値
			})
		}
	}
	return generatedPoints
}
