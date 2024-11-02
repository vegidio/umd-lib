package pkg

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnmarshalJSON_ValidEpochTime(t *testing.T) {
	var et EpochTime
	err := json.Unmarshal([]byte("1633072800"), &et)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := time.Unix(1633072800, 0)
	if !et.Time.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, et.Time)
	}
}

func TestUnmarshalJSON_InvalidEpochTime(t *testing.T) {
	var et EpochTime
	err := json.Unmarshal([]byte("-1"), &et)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	expectedErr := "invalid epoch time"
	if err.Error() != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err.Error())
	}
}

func TestUnmarshalJSON_InvalidJSON(t *testing.T) {
	var et EpochTime
	err := json.Unmarshal([]byte("\"invalid\""), &et)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
