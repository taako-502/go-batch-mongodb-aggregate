package aggregate

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/taako-502/go-batch-mongodb-aggregate/pkg/infrastructure"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func AggregateByGo(ctx context.Context, client *mongo.Client, isPrint bool) error {
	// MongoDBからすべてのpointsドキュメントを取得
	points := infrastructure.Find(ctx, client)

	userPoints := make(map[bson.ObjectID]int)
	for _, p := range points {
		userPoints[p.UserID] += p.Point
	}

	var leaderboards []infrastructure.Leaderboard
	now := time.Now()
	for userID, point := range userPoints {
		leaderboards = append(leaderboards, infrastructure.Leaderboard{
			UserID:     userID,
			Method:     "go",
			TotalPoint: point,
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
		_, err := infrastructure.UpsertLeaderboard(ctx, client, &l)
		if err != nil {
			return fmt.Errorf("failed to upsert leaderboard: %w", err)
		}
	}

	// 結果を表示
	if isPrint {
		fmt.Println("")
		fmt.Println("Goで集計した結果")
		for _, l := range leaderboards {
			fmt.Printf("UserID: %v, TotalPoint: %v\n", l.UserID, l.TotalPoint)
		}
	}

	return nil
}
