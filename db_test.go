package stravautil

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestBasicDB(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, collection, err := getCollection()
	if err != nil {
		slog.Debug("error", "err", err)
		t.Errorf("error getting db collection: %s", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			fmt.Printf("Error disconnecting mongo: %v\n", err)
		}
	}()

	filter := bson.M{}

	slog.Debug(fmt.Sprintf("filter: %+v", filter))

	var activities []*DetailedActivity

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		slog.Debug("error", "err", err)
		t.Errorf("error pulling activities from db: %v", err)
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &activities); err != nil {
		slog.Debug("error", "err", err)
		t.Errorf("error iterating on cursor for activities from db: %v", err)
	}

}

func TestCheckIndexes(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, coll, err := getCollection()
	if err != nil {
		slog.Debug("error", "err", err)
		t.Errorf("error getting db collection: %s", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			fmt.Printf("Error disconnecting mongo: %v\n", err)
		}
	}()

	if err := CheckIndexes(ctx); err != nil {
		t.Errorf("error checking indexes: %v", err)
	}

	cursor, err := coll.Indexes().List(ctx, &options.ListIndexesOptions{})
	if err != nil {
		t.Errorf("error getting cursor to list indexes: %v", err)
	}
	if err := CheckIndexes(ctx); err != nil {
		t.Errorf("error listing indexes: %v", err)
	}

	var indexes []bson.M
	if err = cursor.All(ctx, &indexes); err != nil {
		t.Errorf("error advancing cursor for listing indexes: %v", err)
	}
	fmt.Printf("indexes:\n")
	fmt.Println(indexes)
	spew.Dump(indexes)
	fmt.Printf("len(indexes):\n")
	fmt.Println(len(indexes))
	if len(indexes) != 2 {
		t.Errorf("unexpected length of indexes on activities collection (expected 2): %v", len(indexes))
	}
}

//
//
// }
