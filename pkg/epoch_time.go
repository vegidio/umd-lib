package pkg

import (
	"encoding/json"
	"errors"
	"time"
)

// EpochTime is a wrapper around time.Time to handle JSON unmarshalling of epoch time.
type EpochTime struct {
	time.Time
}

func (t *EpochTime) UnmarshalJSON(b []byte) error {
	var epoch float64
	if err := json.Unmarshal(b, &epoch); err != nil {
		return err
	}

	if epoch < 0 {
		return errors.New("invalid epoch time")
	}

	t.Time = time.Unix(int64(epoch), 0)
	return nil
}
