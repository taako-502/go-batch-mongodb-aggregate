package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/taako-502/go-batch-mongodb-aggregate/aggregate"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	startTime := time.Now()
	aggregate.AggregateByMongoDB(ctx, client)
	elapsed := time.Since(startTime)
	fmt.Printf("MongoDBによる集計にかかった時間: %s\n", elapsed)

	startTime = time.Now()
	aggregate.AggregateByGo()
	elapsed = time.Since(startTime)
	fmt.Printf("Golangによる集計にかかった時間: %s\n", elapsed)

	fmt.Println("順位表を更新しました")
}
