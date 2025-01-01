package stravautil

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
)

func ProcessNewActivities(datestring string, activity_id uint64) {
	slog.Debug("starting call sudo")
	defer slog.Debug("ending call sudo")
	home, err := os.UserHomeDir()
	if err != nil {
		wrappedErr := fmt.Errorf("Error getting my homedir: %w", err)
		slog.Error("got an error:", "wrappedErr", wrappedErr)
		return
	}
	command_argv0 := fmt.Sprintf("%s/new_strava_activity.sh", home)
	cmd := exec.Command(command_argv0, datestring, fmt.Sprintf("%d", activity_id))
	cmd.Env = os.Environ()
	slog.Debug(fmt.Sprintf("cmd: %+v", cmd))
	out, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Error calling sudo", "err", err)
		if ntfyErr := errorNotify(err); ntfyErr != nil {
			log.Println("got ntfy error:", ntfyErr)
		}

	} else {
		slog.Debug("> new_strava_activity.sh:")
		slog.Debug(string(out))
	}
}
