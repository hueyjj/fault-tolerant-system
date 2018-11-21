package store

import (
	"fmt"
)

type Store struct {
	KeyvalMap map[string]string
}

func New() *Store {
	return &Store{
		KeyvalMap: make(map[string]string),
	}
}

// Put places val into key, and returns true
// if the value was replaced
func (s *Store) Put(key, val string) bool {
	_, exists := s.KeyvalMap[key]
	s.KeyvalMap[key] = val
	return exists
}

func (s *Store) Exists(key string) bool {
	_, exists := s.KeyvalMap[key]
	if !exists {
		return false
	}
	return true
}

func (s *Store) Get(key string) (string, error) {
	val, exists := s.KeyvalMap[key]
	if !exists {
		return "", fmt.Errorf("key %s does not exit in the map", key)
	}
	return val, nil
}

func (s *Store) Delete(key string) error {
	_, exists := s.KeyvalMap[key]
	if !exists {
		return fmt.Errorf("key %s does not exist in the map", key)
	}
	delete(s.KeyvalMap, key)
	return nil
}

func (s *Store) Count() int {
	return len(s.KeyvalMap)
}
