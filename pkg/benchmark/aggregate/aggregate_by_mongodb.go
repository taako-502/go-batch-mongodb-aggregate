package aggregate

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/infrastructure"
	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (a *Aggregate) AggregateByMongoDB(ctx context.Context, client *mongo.Client, isPrint bool) error {
	results, err := infrastructure.AggregateUserPoints(client, ctx)
	if err != nil {
		return fmt.Errorf("failed to aggregate points: %w", err)
	}

	// rankingデータベースのrankingテーブルを更新
	var leaderboards []model.Leaderboard
	now := time.Now()
	for _, r := range results {
		leaderboards = append(leaderboards, model.Leaderboard{
			UserID:     r.UserID,
			Method:     "mongodb",
			TotalPoint: r.TotalPoint,
			CreatedAt:  bson.NewDateTimeFromTime(now),
			UpdatedAt:  bson.NewDateTimeFromTime(now),
		})
	}

	// 順位格納
	sort.Slice(leaderboards, func(i, j int) bool {
		return leaderboards[i].TotalPoint > leaderboards[j].TotalPoint
	})

	for i := range leaderboards {
		leaderboards[i].Rank = i + 1
	}

	// rankingデータベースのrankingテーブルを更新
	for _, l := range leaderboards {
		if _, err := a.infrastructure.UpsertLeaderboard(ctx, client, &l); err != nil {
			return fmt.Errorf("failed to upsert leaderboard: %w", err)
		}
	}

	// 結果を表示
	if isPrint {
		fmt.Println("")
		fmt.Println("MongoDBの集計関数で集計した結果")
		for _, l := range leaderboards {
			fmt.Printf("UserID: %v, TotalPoint: %v\n", l.UserID, l.TotalPoint)
		}
	}

	return nil
}
