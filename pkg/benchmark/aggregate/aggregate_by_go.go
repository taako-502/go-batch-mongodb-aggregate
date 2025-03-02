package aggregate

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (a *Aggregate) AggregateByGo(ctx context.Context, client *mongo.Client) error {
	points := a.Infrastructure.Find(ctx, client)
	userPoints := make(map[bson.ObjectID]int)
	for _, p := range points {
		userPoints[p.UserID] += p.Point
	}

	var leaderboards []model.Leaderboard
	now := time.Now()
	for userID, point := range userPoints {
		leaderboards = append(leaderboards, model.Leaderboard{
			UserID:     userID,
			Method:     "go",
			TotalPoint: point,
			CreatedAt:  bson.NewDateTimeFromTime(now),
			UpdatedAt:  bson.NewDateTimeFromTime(now),
		})
	}

	sort.Slice(leaderboards, func(i, j int) bool {
		return leaderboards[i].TotalPoint > leaderboards[j].TotalPoint
	})

	for i := range leaderboards {
		leaderboards[i].Rank = i + 1
	}

	for _, l := range leaderboards {
		if _, err := a.Infrastructure.UpsertLeaderboard(ctx, client, &l); err != nil {
			return fmt.Errorf("failed to upsert leaderboard: %w", err)
		}
	}

	return nil
}
