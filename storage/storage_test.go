package storage

import (
	"testing"
)

func TestStorage_SetInt(t *testing.T) {
	s := New()

	s.SetInt("key", 1, 0)
	item := s.GetItem("key")
	if item == nil {
		t.Error("Must contains key")
	}
	if item.Int != 1 {
		t.Error("Must be equal 1")
	}

	s.SetInt("key2", 0, 0)
	item = s.GetItem("key2")
	if item == nil {
		t.Error("Must contains key")
	}
	if item.Int != 0 {
		t.Error("Must be equal 1")
	}
}

func TestStorage_SetInt_Zero(t *testing.T) {
	s := New()

	s.SetInt("key", 0, 0)
	item := s.GetItem("key")
	if item == nil {
		t.Error("Must contains key")
	}
	if item.Int != 0 {
		t.Error("Must be equal 0")
	}
}

func TestStorage_SetString(t *testing.T) {
	s := New()

	s.SetString("key", "val", 0)
	item := s.GetItem("key")
	if item == nil {
		t.Error("Must contains key")
	}
	if item.String != "val" {
		t.Error("Must be equal `val`")
	}
}

func TestStorage_SetString_Empty(t *testing.T) {
	s := New()

	s.SetString("key", "", 0)
	item := s.GetItem("key")
	if item != nil {
		t.Error("Must be not found")
	}
}

func TestStorage_SetIntSlice(t *testing.T) {
	s := New()

	s.SetIntSlice("key", []int{1,2,3,4}, 0)
	item := s.GetItem("key")
	if item == nil {
		t.Error("Must contains key")
	}
	if item.IntSlice == nil {
		t.Error("Must contains slice of int")
	}
	if len(item.IntSlice) != len([]int{1,2,3,4}) {
		t.Fatal("Must be equal", len(item.IntSlice), len([]int{1,2,3,4}))
	}
	for i, v := range []int{1,2,3,4} {
		if item.IntSlice[i] != v {
			t.Error("Must be equal")
		}
	}
}

func TestStorage_SetStringMap(t *testing.T) {
	s := New()

	value := map[string]string{
		"ab": "qwe",
		"bc": "asd",
	}
	s.SetStringMap("key", value, 0)

	item := s.GetItem("key")
	if item == nil {
		t.Fatal("Must contains key")
	}
	if item.StringMap == nil {
		t.Fatal("Must contains value")
	}
	if len(value) != len(item.StringMap) {
		t.Fatal("Must be equal", len(value), len(item.StringMap))
	}
	for k, v := range item.StringMap {
		if v != item.StringMap[k] {
			t.Fatal("Must be equal", v, item.StringMap[k])
		}
	}
}

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

func TestStorage_GetIntFromList(t *testing.T) {
	s := New()

	value := []int{1,2,3,4}
	s.SetIntSlice("key", value, 0)

	v, ok := s.GetIntFromList("key", 0)
	if !ok {
		t.Error("Must return", value[0])
	}
	if v != value[0] {
		t.Error("Must be equal", v, value[0])
	}

	v, ok = s.GetIntFromList("key", -1)
	if ok {
		t.Error("Must return false")
	}

	v, ok = s.GetIntFromList("key", 100)
	if ok {
		t.Error("Must return false")
	}
}