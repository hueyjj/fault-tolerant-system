package store

import (
	"fmt"
)

type Store struct {
	keyvalMap map[string]string
}

func New() *Store {
	return &Store{
		keyvalMap: make(map[string]string),
	}
}

// Put places val into key, and returns true
// if the value was replaced
func (s *Store) Put(key, val string) bool {
	val, exists := s.keyvalMap[key]
	s.keyvalMap[key] = val
	return exists
}

func (s *Store) Get(key string) (string, error) {
	val, exists := s.keyvalMap[key]
	if !exists {
		return "", fmt.Errorf("key %s does not exit in the map", key)
	}
	return val, nil
}

func (s *Store) Delete(key string) error {
	_, exists := s.keyvalMap[key]
	if !exists {
		return fmt.Errorf("key %s does not exit in the map", key)
	}
	delete(s.keyvalMap, key)
	return nil
}

func (s *Store) Count() int {
	return len(s.keyvalMap)
}
