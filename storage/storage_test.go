package storage

import (
	"testing"
)

func TestGetString(t *testing.T) {
	s := New()
	s.Set("key", "testValue", 0)
	value := s.Get("key")
	valueStr := value.(string)
	if valueStr != "testValue" {
		t.Error()
	}
}

func TestGetFromList(t *testing.T) {
	s := New()
	s.Set("key", []string{"a", "b", "c"}, 0)
	valueFromList, err := s.GetFromList("key", 0)
	if err != nil {
		t.Error("Err must be nil: %v", err)
	}
	if valueFromList != "a" {
		t.Error("Must be equal to `a`")
	}
}

func TestGetFromDict(t *testing.T) {
	s := New()
	s.Set("key", map[string]string{"a": "abc", "b": "bcd"}, 0)
	value, err := s.GetFromDict("key", "a")
	if err != nil {
		t.Error("Err must be nil: %v", err)
	}
	if value != "abc" {
		t.Error("Must be equal to `abc`")
	}
}