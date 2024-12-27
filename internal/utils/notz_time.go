package utils

import (
	"encoding/json"
	"errors"
	"time"
)

// NotzTime is a wrapper around time.Time to handle JSON unmarshalling time without time zone.
type NotzTime struct {
	time.Time
}

func (t *NotzTime) UnmarshalJSON(b []byte) error {
	var timestamp string
	if err := json.Unmarshal(b, &timestamp); err != nil {
		return err
	}

	if timestamp == "" {
		return errors.New("invalid timestamp")
	}

	tt, err := time.Parse("2006-01-02T15:04:05", timestamp)
	if err != nil {
		return err
	}

	t.Time = tt
	return nil
}
