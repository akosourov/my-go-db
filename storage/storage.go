package storage

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Item struct {
	Value      interface{}
	expiration int64
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

func (s *Storage) Set(key string, value interface{}, ttl int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var exp int64
	if ttl > 0 {
		exp = time.Now().UnixNano() + ttl*1000*1000*1000
	}
	item := Item{
		Value:      value,
		expiration: exp,
	}
	s.items[key] = &item
}

func (s *Storage) Get(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item := s.items[key]
	if item != nil {
		return item.Value
	} else {
		return nil
	}
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
