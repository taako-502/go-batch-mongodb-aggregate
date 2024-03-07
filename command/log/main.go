package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type logEntry struct {
	ID             primitive.ObjectID `bson:"_id"`
	Method         string             `bson:"method"`
	NumberOfUsers  int                `bson:"numberOfUsers"`
	NumberOfPoints int                `bson:"numberOfPoints"`
	Elapsed        int                `bson:"elapsed"`
}

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

	cursor, err := client.Database("aggregate").Collection("logs").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var logs []logEntry
	if err = cursor.All(ctx, &logs); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID\tMethod\tNumber Of Users\tNumber Of Points\tElapsed Time\t\n")
	for _, log := range logs {
		fmt.Printf("%s\t%s\t%d\t%d\t%d\t\n", log.ID.Hex(), log.Method, log.NumberOfUsers, log.NumberOfPoints, log.Elapsed)
	}
}
