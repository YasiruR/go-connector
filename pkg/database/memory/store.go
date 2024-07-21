package memory

import (
	"fmt"
	"github.com/YasiruR/connector/core/pkg"
	"sync"
)

type Store struct {
	maps []pkg.DataStore
}

func NewStore() *Store {
	return &Store{maps: make([]pkg.DataStore, 0)}
}

func (s *Store) NewDataStore() pkg.DataStore {
	m := &Map{data: new(sync.Map)}
	s.maps = append(s.maps, m)
	return m
}

type Map struct {
	data *sync.Map
}

func (m *Map) Set(key string, value any) error {
	m.data.Store(key, value)
	return nil
}

func (m *Map) Get(key string) (any, error) {
	val, ok := m.data.Load(key)
	if val == nil || !ok {
		return nil, fmt.Errorf("key (%s) not found in memory database", key)
	}

	return val, nil
}
