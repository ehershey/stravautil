package stravautil

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckIndexes(ctx context.Context) error {

	client, coll, err := getCollection()
	if err != nil {
		return fmt.Errorf("Error getting collection for indexes: %w", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			fmt.Printf("Error disconnecting mongo: %v\n", err)
		}
	}()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "strava_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	name, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("Error creating index: %w", err)
	}
	fmt.Println("Name of Index Created: " + name)
	return nil
}
