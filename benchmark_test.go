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
	{numberOfUsers: 1000, numberOfPoints: 1},
	{numberOfUsers: 1000, numberOfPoints: 10},
	{numberOfUsers: 1000, numberOfPoints: 100},
	{numberOfUsers: 1000, numberOfPoints: 1000},
	{numberOfUsers: 1000, numberOfPoints: 10000},
	{numberOfUsers: 1000, numberOfPoints: 100000},
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
	if err := benchmark.Cleanup(ctx, a.Infrastructure.SourceUserCol); err != nil {
		return err
	}
	if err := benchmark.Cleanup(ctx, a.Infrastructure.SourcePointCol); err != nil {
		return err
	}
	if err := benchmark.Cleanup(ctx, a.Infrastructure.AggregateLeaderboardCol); err != nil {
		return err
	}
	return nil
}

func seed(ctx context.Context, numberOfUsers, numberOfPoints int) error {
	users := model.GenerateUsers(numberOfUsers)
	if _, err := a.Infrastructure.SourceUserCol.InsertMany(ctx, users); err != nil {
		return err
	}

	var userIDs []bson.ObjectID
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	points := model.GeneratePoints(userIDs, numberOfPoints)
	if _, err := a.Infrastructure.SourcePointCol.InsertMany(ctx, points); err != nil {
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

			if err := seed(ctx, n.numberOfUsers, n.numberOfPoints); err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for b.Loop() {
				a.AggregateByMongoDB(ctx, client)
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

			if err := seed(ctx, n.numberOfUsers, n.numberOfPoints); err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for b.Loop() {
				a.AggregateByGo(ctx, client)
			}
		})
	}
}
