package stravautil

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"testing"
)

func TestSubUpdate(t *testing.T) {
	str := `{"aspect_type":"delete","event_time":1604072850,"object_id":4222366652,"object_type":"activity","owner_id":3968,"subscription_id":138599,"updates":{}}`
	res := Update{}
	err := json.Unmarshal([]byte(str), &res)
	if err != nil {
		slog.Debug("error", err)
	}
	verifyUpdate(res, t)
}
func verifyUpdate(res Update, t *testing.T) {
	if res.AspectType != "delete" {
		t.Errorf("Couldn't parse aspect type")
	}
	if res.EventTime != 1604072850 {
		t.Errorf("Couldn't parse event time")
	}
	if res.ObjectID != 4222366652 {
		t.Errorf("Couldn't parse object id")
	}
	if res.ObjectType != "activity" {
		t.Errorf("Couldn't parse object type")
	}
	if res.OwnerID != 3968 {
		t.Errorf("Couldn't parse owner id")
	}
	if res.SubscriptionID != 138599 {
		t.Errorf("Couldn't parse subscription id")
	}
	if !reflect.DeepEqual(res.Updates, map[string]string{}) {
		t.Errorf("Couldn't parse updates")
	}
}

func ExampleUpdate_reencode() {
	str := `{"aspect_type":"delete","event_time":1604072850,"object_id":4222366652,"object_type":"activity","owner_id":3968,"subscription_id":138599,"updates":{}}`
	res := Update{}
	json.Unmarshal([]byte(str), &res)
	// fmt.Println(res)
	data, err := json.Marshal(res)
	if err != nil {
		fmt.Println("couldn't remarshall json!")
	}
	fmt.Println(string(data))
	// Output:
	// {"aspect_type":"delete","event_time":1604072850,"object_id":4222366652,"object_type":"activity","owner_id":3968,"subscription_id":138599,"updates":{}}
}

// make the object appear as json without explicit unmarshaling
func ExampleUpdate_implicitjson() {
	str := `{"aspect_type":"delete","event_time":1604072850,"object_id":4222366652,"object_type":"activity","owner_id":3968,"subscription_id":138599,"updates":{}}`
	res := Update{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)
	// Output:
	// {"aspect_type":"delete","event_time":1604072850,"object_id":4222366652,"object_type":"activity","owner_id":3968,"subscription_id":138599,"updates":{}}
}

func TestStreamDecodeIntoObj(t *testing.T) {
	str := `{"aspect_type":"delete","event_time":1604072850,"object_id":4222366652,"object_type":"activity","owner_id":3968,"subscription_id":138599,"updates":{}}`

	decoder := json.NewDecoder(strings.NewReader(str))
	var rawobj Update
	err := decoder.Decode(&rawobj)
	if err != nil {
		slog.Warn("error", err)
		t.Errorf("error decoding: %s", err)
	}
	verifyUpdate(rawobj, t)
}
