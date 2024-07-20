package memory

import (
	"fmt"
	"sync"
)

type Store struct {
	data *sync.Map
}

func NewStore() *Store {
	return &Store{data: new(sync.Map)}
}

func (s *Store) Set(key string, value any) error {
	s.data.Store(key, value)
	return nil
}

func (s *Store) Get(key string) (any, error) {
	val, ok := s.data.Load(key)
	if val == nil || !ok {
		return nil, fmt.Errorf("key (%s) not found in memory database", key)
	}

	return val, nil
}
