package aggregate

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/taako-502/go-batch-mongodb-aggregate/infrastructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AggregateByMongoDB(ctx context.Context, client *mongo.Client, isPrint bool) error {
	results, err := infrastructure.AggregateUserPoints(client, ctx)
	if err != nil {
		return fmt.Errorf("failed to aggregate points: %w", err)
	}

	// rankingデータベースのrankingテーブルを更新
	var leaderboards []infrastructure.Leaderboard
	now := time.Now()
	for _, r := range results {
		leaderboards = append(leaderboards, infrastructure.Leaderboard{
			UserID:     r.UserID,
			Method:     "mongodb",
			TotalPoint: r.TotalPoint,
			CreatedAt:  primitive.NewDateTimeFromTime(now),
			UpdatedAt:  primitive.NewDateTimeFromTime(now),
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
		_, err := infrastructure.UpsertLeaderboard(ctx, client, &l)
		if err != nil {
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
