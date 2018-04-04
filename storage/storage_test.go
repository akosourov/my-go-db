package storage

import (
	"testing"
)

func TestGetString(t *testing.T) {
	s := New()
	s.Set("key", "testValue", 0)
	s.Set("key2", 123, 0)

	value, _ := s.GetString("key")
	if value != "testValue" {
		t.Error("Must be equal")
	}

	if _, ok := s.GetString("key2"); !ok {
		t.Error("Must be not found")
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
