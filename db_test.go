package stravautil

import (
	"context"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestBasicDB(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, collection, err := getCollection()
	if err != nil {
		log.Println(err)
		t.Errorf("error getting db collection: %s", err)
	}
	defer client.Disconnect(ctx)

	filter := bson.M{}

	log.Println("filter:", filter)

	var activities []*DetailedActivity

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		t.Errorf("error pulling activities from db: %v", err)
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &activities); err != nil {
		log.Println(err)
		t.Errorf("error iterating on cursor for activities from db: %v", err)
	}

}
