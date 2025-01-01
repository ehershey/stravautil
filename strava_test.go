package stravautil

import (
	"log"
	"testing"

	"log/slog"
)

func TestDelete(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	stravaID := uint64(5545947515)
	// insert it first
	//
	ProcessNewActivities("", stravaID)
	// then delete it
	//
	log.Println("test")
	old_activity, err := Delete_activity(stravaID)
	if err != nil {
		t.Errorf("error deleting activity: %v", err)
	}
	log.Println("old_activity:", old_activity)
}
