package storage

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Item struct {
	Value          interface{}

	ValueInt       int
	ValueStr       string
	ValueSliceInt  []int
	ValueSliceStr  []string
	expiration     int64
}

type Storage struct {
	mu    *sync.RWMutex
	items map[string]*Item
}

func New() *Storage {
	return &Storage{
		mu:    new(sync.RWMutex),
		items: make(map[string]*Item),
	}
}

func calculateExpiration(ttl int) int64 {
	var exp int64
	if ttl > 0 {
		exp = time.Now().UnixNano() + int64(ttl)*1000*1000*1000
	}
	return exp
}

func (s *Storage) SetString(key, value string, ttl int) {
	exp := calculateExpiration(ttl)

	item := new(Item)
	item.expiration = exp
	item.ValueStr = value

	s.mu.Lock()
	s.items[key] = item
	s.mu.Unlock()
}

func (s *Storage) SetInt(key string, value, ttl int) {
	exp := calculateExpiration(ttl)

	item := new(Item)
	item.expiration = exp
	item.ValueInt = value
	// todo value = 0

	s.mu.Lock()
	s.items[key] = item
	s.mu.Unlock()
}

func (s *Storage) Set(key string, value interface{}, ttl int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var exp int64
	if ttl > 0 {
		exp = time.Now().UnixNano() + int64(ttl)*1000*1000*1000
	}

	item := new(Item)
	item.expiration = exp

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		item.ValueStr = v.Interface().(string)
	case reflect.Int:
		fmt.Println("Value is int")
		item.ValueInt = v.Interface().(int)
	case reflect.Slice:
		if v.Len() > 0 {
			elemV := v.Index(0)
			switch elemV.Kind() {
			case reflect.Int:
				item.ValueSliceInt = v.Interface().([]int)
			case reflect.String:
				item.ValueSliceStr = v.Interface().([]string)
			}
		} else {
			return
		}
	default:
		fmt.Println("Value is %v %T", value, value)
		return
	}

	s.items[key] = item
}


func (s *Storage) GetItem(key string) *Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items[key]
}


func (s *Storage) Get(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item := s.items[key]
	if item != nil {
		if item.ValueStr != "" {
			return item.ValueStr
		} else if item.ValueInt > 0 {
			// todo int == 0
			return item.ValueInt
		}
	}
	return nil
}

func (s *Storage) GetString(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item := s.items[key]
	if item != nil {
		return item.ValueStr, true
	}
	return "", false
}

func (s *Storage) GetInt(key string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item := s.items[key]
	if item != nil {
		return item.ValueInt, true
	}
	return 0, false
}

func (s *Storage) GetFromList(key string, idx int) (interface{}, error) {
	s.mu.RLock()
	item, ok := s.items[key]
	s.mu.RUnlock()
	if !ok {
		return nil, errors.New(fmt.Sprintf("Key: %s does not exist", key))
	}
	value := item.Value
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		if idx >= v.Len() {
			return nil, errors.New("Index out of range")
		}
		slice, ok := value.([]string)
		if !ok {
			return nil, errors.New("Unsupported type in storage")
		}
		return slice[idx], nil
	}
	return nil, errors.New(fmt.Sprintf("Item: %v is not a list", item.Value))
}


func (s *Storage) GetFromDict(key, dkey string) (interface{}, error) {
	s.mu.Lock()
	item, ok := s.items[key]
	s.mu.RUnlock()
	if !ok {
		return nil, errors.New(fmt.Sprintf("Key: %s does not exist", key))
	}
	v := reflect.ValueOf(item.Value)
	switch v.Kind() {
	case reflect.Map:
		m, ok := item.Value.(map[string]string)
		if !ok {
			return nil, errors.New("Unsupportd type in storage")
		}
		return m[dkey], nil
	}
	return nil, errors.New(fmt.Sprintf("Item: %v is not a list", item.Value))
}

func (s *Storage) Remove(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.items[key]
	if !ok {
		return errors.New(fmt.Sprintf("Key: %s does not exist", key))
	}
	delete(s.items, key)
	return nil
}

func (s *Storage) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := []string{}
	for k := range s.items {
		keys = append(keys, k)
	}
	return keys
}

func (s *Storage) DeleteExpired() {
	now := time.Now().UnixNano()
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, item := range s.items {
		if item.expiration > 0 && item.expiration <= now {
			delete(s.items, k)
		}
	}
}
