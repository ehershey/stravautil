package stravautil

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func checkIndexes(ctx context.Context) error {

	client, coll, err := getCollection()
	if err != nil {
		return fmt.Errorf("Error getting collection for indexes: %w", err)
	}
	defer client.Disconnect(ctx)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"strava_id", 1}},
		Options: options.Index().SetUnique(true),
	}
	name, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("Error creating index: %w", err)
	}
	fmt.Println("Name of Index Created: " + name)
	return nil
}
