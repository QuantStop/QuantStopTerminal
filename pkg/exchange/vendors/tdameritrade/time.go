package tdameritrade

import (
	"fmt"
	"strings"
	"time"
)

// Timestamps
//
// No specific documentation on time formats other than 'valid ISO 8601 format'

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		var ptr Time
		*t = ptr
		return nil
	}

	layouts := []string{
		"2006-01-02T15:04:05.000000Z",
		"2006-01-02 15:04:05+00",
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05.999999+00",
		time.RFC3339Nano,
		time.RFC3339,
	}
	var parsedTime time.Time
	var err error
	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, strings.ReplaceAll(string(data), "\"", ""))
		if err != nil {
			continue
		}
		break
	}
	if err != nil {
		return fmt.Errorf("time %s in unhandled format", data)
	}
	*t = Time(parsedTime)
	return nil
}

func (t *Time) Time() time.Time {
	return time.Time(*t)
}
