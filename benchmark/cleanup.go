package benchmark

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Cleanup(ctx context.Context, collection *mongo.Collection) error {
	if _, err := collection.DeleteMany(ctx, bson.M{}); err != nil {
		return fmt.Errorf("failed to delete all documents: %w", err)
	}
	return nil
}
