package infrastructure

import "go.mongodb.org/mongo-driver/v2/mongo"

type Infrastructure struct {
	sourcePointCol          *mongo.Collection
	aggregateLeaderboardCol *mongo.Collection
}

func NewInfrastructure(client *mongo.Client) *Infrastructure {
	return &Infrastructure{
		sourcePointCol:          client.Database("source").Collection("points"),
		aggregateLeaderboardCol: client.Database("aggregate").Collection("leaderboard"),
	}
}
