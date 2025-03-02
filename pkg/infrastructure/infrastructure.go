package infrastructure

import "go.mongodb.org/mongo-driver/v2/mongo"

type Infrastructure struct {
	SourcePointCol          *mongo.Collection
	SourceUserCol           *mongo.Collection
	AggregateLeaderboardCol *mongo.Collection
}

func NewInfrastructure(client *mongo.Client) *Infrastructure {
	return &Infrastructure{
		SourcePointCol:          client.Database("source").Collection("points"),
		SourceUserCol:           client.Database("source").Collection("users"),
		AggregateLeaderboardCol: client.Database("aggregate").Collection("leaderboard"),
	}
}
