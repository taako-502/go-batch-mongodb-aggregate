package infrastructure

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Leaderboard struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"userId"`
	Method     string             `bson:"method"`
	TotalPoint int                `bson:"totalPoint"`
	Rank       int                `bson:"rank"`
	CreatedAt  primitive.DateTime `bson:"createdAt"`
	UpdatedAt  primitive.DateTime `bson:"updatedAt"`
}

func UpsertLeaderboard(ctx context.Context, client *mongo.Client, Leaderboard *Leaderboard) (*Leaderboard, error) {
	filter := bson.M{"userId": Leaderboard.UserID, "method": Leaderboard.Method}
	update := bson.M{
		"$set": bson.M{
			"totalPoint": Leaderboard.TotalPoint,
			"rank":       Leaderboard.Rank,
			"updatedAt":  Leaderboard.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"createdAt": Leaderboard.CreatedAt,
		},
	}
	upsert := true
	opts := options.UpdateOptions{Upsert: &upsert}
	result, err := client.Database("aggregate").Collection("leaderboard").UpdateOne(ctx, filter, &update, &opts)
	if err != nil {
		return nil, fmt.Errorf("failed to update ranking: %v", err)
	}

	if result.UpsertedID != nil {
		Leaderboard.ID = result.UpsertedID.(primitive.ObjectID)
	}

	return Leaderboard, nil
}
