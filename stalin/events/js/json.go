package js

import (
	"encoding/json"
	. "stalin/events"
)

// {"events":[{"time":...},{},{}]}
// {"data_type":"","data":""}
func EventsFromJson(data []byte) (*Events, error) {
	result := NewEvents()
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	return result, nil
}
