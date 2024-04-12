package stravautil

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Activity struct {
	Strava_id        int       `bson:"strava_id" json:"strava_id"`
	External_id      string    `bson:"external_id" json:"external_id"`
	Start_date_local time.Time `bson:"start_date_local" json:"start_date_local"`
	End_date_local   time.Time `bson:"end_date_local" json:"end_date_local"`
}

type DetailedActivity struct {
	StravaId         int       `bson:"strava_id" json:"strava_id"`
	External_id      string    `bson:"external_id" json:"external_id"`
	StartDate        time.Time `bson:"start_date" json:"start_date"`
	Distance         float64   `bson:"distance" json:"distance"`
	ElapsedTime      int       `bson:"elapsed_time" json:"elapsed_time"`
	Start_date_local time.Time `bson:"start_date_local" json:"start_date_local"`
	End_date_local   time.Time `bson:"end_date_local" json:"end_date_local"`
	Name             string    `bson:"name" json:"name"`
	Type             string    `bson:"type" json:"type"`
}

const db_name = "strava"
const collection_name = "activities"
const default_db_uri = "mongodb://localhost:27017"

func getCollection() (*mongo.Client, *mongo.Collection, error) {
	strava_env_uri := os.Getenv("STRAVA_MONGODB_URI") // prefer this
	env_uri := os.Getenv("MONGODB_URI")               // second choice
	db_uri := default_db_uri                          // third choice
	if strava_env_uri != "" {
		db_uri = strava_env_uri
	} else if env_uri != "" {
		db_uri = env_uri
	}
	clientoptions := options.Client().ApplyURI(db_uri)

	url, err := url.Parse(db_uri)
	if err != nil {
		log.Println("got an error:", err)
		return nil, nil, err
	}
	log.Printf("connecting to mongodb at: %s\n", url.Redacted())
	client, err := mongo.NewClient(clientoptions)
	if err != nil {
		log.Println("got an error:", err)
		return nil, nil, err
	}
	log.Println("got client:", client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("connecting")
	err = client.Connect(ctx)
	collection := client.Database(db_name).Collection(collection_name)
	return client, collection, nil
}

func Delete_activity(activity_id uint64) (*Activity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, collection, err := getCollection()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	// filter := bson.D{{Key: "strava_id", Value: activity_id}}
	filter := bson.D{{Key: "strava_id", Value: activity_id}}
	//filter := bson.D{}

	log.Println("filter:", filter)

	var old_activity Activity

	err = collection.FindOne(ctx, filter).Decode(&old_activity)

	if err != nil {
		return nil, err
	}

	old_activity_json, err := json.Marshal(old_activity)
	if err != nil {
		return nil, err
	}
	log.Println("got old_activity:", old_activity)
	log.Printf("old_activity_json: %s\n", old_activity_json)

	// do the delete
	result, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return nil, err
	}
	log.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)

	end_date_year, end_date_month, end_date_day := old_activity.End_date_local.Date()
	end_date_string := fmt.Sprintf("%d-%02d-%02d", end_date_year, end_date_month, end_date_day)

	start_date_year, start_date_month, start_date_day := old_activity.Start_date_local.Date()
	start_date_string := fmt.Sprintf("%d-%02d-%02d", start_date_year, start_date_month, start_date_day)

	// trigger processing new activities
	//
	log.Println("starting call goroutine for start date:", start_date_string)
	go ProcessNewActivities(start_date_string, 0)

	if end_date_string != start_date_string {
		log.Println("starting call second goroutine for end date:", end_date_string)
		go ProcessNewActivities(end_date_string, 0)
	}
	log.Println("ending call goroutine")

	return &old_activity, nil
}
func GetActivities(weeks_back int) ([]*DetailedActivity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, collection, err := getCollection()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	minimum_timestamp := time.Now().AddDate(0, 0, int(7*-1)*int(weeks_back))

	// filter := bson.D{{Key: "strava_id", Value: activity_id}}
	//filter := bson.D{{Key: "strava_id", Value: activity_id}}
	filter := bson.M{"start_date_local": bson.M{"$gte": minimum_timestamp}}

	log.Println("filter:", filter)

	var activities []*DetailedActivity

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error pulling activities from db: %w", err)
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &activities); err != nil {
		return nil, fmt.Errorf("error iterating on cursor for activities from db: %w", err)
	}

	return activities, nil
}
