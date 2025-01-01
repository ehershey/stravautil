package stravautil

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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
		wrappedErr := fmt.Errorf("Error parsing db_uri: %w", err)
		return nil, nil, wrappedErr
	}
	slog.Debug("connecting to mongodb", "url", url.Redacted())
	client, err := mongo.NewClient(clientoptions)
	if err != nil {
		wrappedErr := fmt.Errorf("Error from mongo.NewClient: %w", err)
		return nil, nil, wrappedErr
	}
	slog.Debug(fmt.Sprintf("got client: %+v", client))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	slog.Debug("connecting")
	err = client.Connect(ctx)
	if err != nil {
		wrappedErr := fmt.Errorf("Error from client.Connect: %w", err)
		return nil, nil, wrappedErr
	}
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
	defer func() {
		slog.Debug("Disconnecting client")
		if err := client.Disconnect(ctx); err != nil {
			slog.Warn("error disconnecting client", "err", err)
		}
	}()

	// filter := bson.D{{Key: "strava_id", Value: activity_id}}
	filter := bson.D{{Key: "strava_id", Value: activity_id}}
	//filter := bson.D{}

	slog.Debug(fmt.Sprintf("filter: %+v", filter))

	var old_activity Activity

	err = collection.FindOne(ctx, filter).Decode(&old_activity)

	if err != nil {
		return nil, err
	}

	old_activity_json, err := json.Marshal(old_activity)
	if err != nil {
		return nil, err
	}
	slog.Debug(fmt.Sprintf("got old_activity: %+v", old_activity))
	slog.Debug(fmt.Sprintf("old_activity_json: %s", old_activity_json))

	// do the delete
	result, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return nil, err
	}
	slog.Debug(fmt.Sprintf("DeleteOne removed %v document(s)\n", result.DeletedCount))

	end_date_year, end_date_month, end_date_day := old_activity.End_date_local.Date()
	end_date_string := fmt.Sprintf("%d-%02d-%02d", end_date_year, end_date_month, end_date_day)

	start_date_year, start_date_month, start_date_day := old_activity.Start_date_local.Date()
	start_date_string := fmt.Sprintf("%d-%02d-%02d", start_date_year, start_date_month, start_date_day)

	// trigger processing new activities
	//
	slog.Debug("starting call goroutine for start date", "start_date_string", start_date_string)
	go ProcessNewActivities(start_date_string, 0)

	if end_date_string != start_date_string {
		slog.Debug("starting call second goroutine for end date", "end_date_string", end_date_string)
		go ProcessNewActivities(end_date_string, 0)
	}
	slog.Debug("ending call goroutine")

	return &old_activity, nil
}

func GetActivities(weeks_back int) ([]*DetailedActivity, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, collection, err := getCollection()
	if err != nil {
		return nil, err
	}

	defer func() {
		slog.Debug("Disconnecting client")
		if err := client.Disconnect(ctx); err != nil {
			slog.Warn("error disconnecting client", "err", err)
		}
	}()

	minimum_timestamp := time.Now().AddDate(0, 0, int(7*-1)*int(weeks_back))

	// filter := bson.D{{Key: "strava_id", Value: activity_id}}
	//filter := bson.D{{Key: "strava_id", Value: activity_id}}
	filter := bson.M{"start_date_local": bson.M{"$gte": minimum_timestamp}}

	slog.Debug(fmt.Sprintf("filter: %+v", filter))

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
