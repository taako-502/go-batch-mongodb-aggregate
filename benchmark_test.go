package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/benchmark"
	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/benchmark/aggregate"
	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/infrastructure"
	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client
var a *aggregate.Aggregate

type pattern struct {
	numberOfUsers  int
	numberOfPoints int
}

var benchmarkPattern = []pattern{
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
	{numberOfUsers: 10000000, numberOfPoints: 10000},
	{numberOfUsers: 10000000, numberOfPoints: 100000},
	{numberOfUsers: 10000000, numberOfPoints: 1000000},
	{numberOfUsers: 100000000, numberOfPoints: 1000000},
}

func init() {
	var err error
	client, err = mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	a = aggregate.NewAggregate(
		infrastructure.NewInfrastructure(client),
	)
}

func cleanup(ctx context.Context) error {
	if err := benchmark.Cleanup(ctx, client.Database("source").Collection("users")); err != nil {
		return err
	}
	if err := benchmark.Cleanup(ctx, client.Database("source").Collection("points")); err != nil {
		return err
	}
	if err := benchmark.Cleanup(ctx, client.Database("aggregate").Collection("leaderboard")); err != nil {
		return err
	}
	return nil
}

// BenchmarkAggregationPipeline Aggregation Pipelineで集計する
func BenchmarkAggregationPipeline(b *testing.B) {
	ctx := context.Background()
	for _, n := range benchmarkPattern {
		b.Run("Benchmark_"+fmt.Sprint(n), func(b *testing.B) {
			if err := cleanup(ctx); err != nil {
				b.Fatal(err)
			}

			// ユーザーデータを挿入
			users := model.GenerateUsers(n.numberOfUsers)
			if _, err := client.Database("source").Collection("users").InsertMany(ctx, users); err != nil {
				b.Fatal(err)
			}

			// ポイントデータを挿入
			var userIDs []bson.ObjectID
			for _, u := range users {
				userIDs = append(userIDs, u.ID)
			}
			points := model.GeneratePoints(userIDs, n.numberOfPoints)
			if _, err := client.Database("source").Collection("points").InsertMany(ctx, points); err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for b.Loop() {
				a.AggregateByMongoDB(ctx, client, false)
			}
		})
	}
}

// BenchmarkGo Goで集計する
func BenchmarkGo(b *testing.B) {
	ctx := context.Background()
	for _, n := range benchmarkPattern {
		b.Run("Benchmark_"+fmt.Sprint(n), func(b *testing.B) {
			if err := cleanup(ctx); err != nil {
				b.Fatal(err)
			}

			users := model.GenerateUsers(n.numberOfUsers)
			if _, err := client.Database("source").Collection("users").InsertMany(ctx, users); err != nil {
				b.Fatal(err)
			}

			var userIDs []bson.ObjectID
			for _, u := range users {
				userIDs = append(userIDs, u.ID)
			}
			points := model.GeneratePoints(userIDs, n.numberOfPoints)
			if _, err := client.Database("source").Collection("points").InsertMany(ctx, points); err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for b.Loop() {
				a.AggregateByGo(ctx, client, false)
			}
		})
	}
}
