package stravautil

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func ProcessNewActivities(datestring string, activity_id uint64) {
	log.Println("starting call sudo")
	home, err := os.UserHomeDir()
	if err != nil {
		wrappedErr := fmt.Errorf("Error getting my homedir: %v", err)
		log.Println("got an error:", wrappedErr)
		return
	}
	command_argv0 := fmt.Sprintf("%s/new_strava_activity.sh", home)
	cmd := exec.Command(command_argv0, datestring, fmt.Sprintf("%d", activity_id))
	log.Println("cmd:", cmd)
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("> new_strava_activity.sh:")
		log.Println(string(out))
	}
	log.Println("ending call sudo")
}
