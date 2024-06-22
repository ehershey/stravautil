package stravautil

import (
	"encoding/json"
	"log/slog"
)

// {"aspect_type":"delete","event_time":1604072850,"object_id":4222366652,"object_type":"activity","owner_id":3968,"subscription_id":138599,"updates":{}}
// {"aspect_type":"update","event_time":1603066007,"object_id":4206062481,"object_type":"activity","owner_id":3968,"subscription_id":138599,"updates":{"title":"NYCRUNS Falling Leaves 5K 2020"}}

// Update post from strava web hook system
type Update struct {
	AspectType     string            `json:"aspect_type"`
	EventTime      uint64            `json:"event_time"`
	ObjectID       uint64            `json:"object_id"`
	ObjectType     string            `json:"object_type"`
	OwnerID        uint64            `json:"owner_id"`
	SubscriptionID uint64            `json:"subscription_id"`
	Updates        map[string]string `json:"updates"`
}

// res := response2{}
// json.Unmarshal([]byte(str), &res)
// fmt.Println(res)
// fmt.Println(res.Fruits[0])

func (update Update) String() string {
	data, err := json.Marshal(update)
	if err != nil {
		slog.Debug("couldn't marshall json!", err)
		return ""
	}
	return string(data)
}
