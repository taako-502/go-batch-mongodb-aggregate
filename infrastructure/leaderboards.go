package infrastructure

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Leaderboard struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	UserID     bson.ObjectID `bson:"userId"`
	Method     string        `bson:"method"`
	TotalPoint int           `bson:"totalPoint"`
	Rank       int           `bson:"rank"`
	CreatedAt  bson.DateTime `bson:"createdAt"`
	UpdatedAt  bson.DateTime `bson:"updatedAt"`
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
	opts := options.UpdateOne().SetUpsert(true)
	result, err := client.Database("aggregate").Collection("leaderboard").UpdateOne(ctx, filter, &update, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to update ranking: %v", err)
	}

	if result.UpsertedID != nil {
		Leaderboard.ID = result.UpsertedID.(bson.ObjectID)
	}

	return Leaderboard, nil
}
