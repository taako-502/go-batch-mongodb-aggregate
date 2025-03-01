package infrastructure

import (
	"context"
	"fmt"

	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func UpsertLeaderboard(ctx context.Context, client *mongo.Client, Leaderboard *model.Leaderboard) (*model.Leaderboard, error) {
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
