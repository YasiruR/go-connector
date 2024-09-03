package memory

import (
	"github.com/YasiruR/connector/domain/pkg"
	"sync"
)

type Store struct {
	maps []pkg.Collection
}

func NewStore(log pkg.Log) *Store {
	log.Info("initialized in-memory database with sync.Map as collections")
	return &Store{maps: make([]pkg.Collection, 0)}
}

func (s *Store) NewCollection() pkg.Collection {
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
		return nil, nil
	}

	return val, nil
}

func (m *Map) GetAll() ([]any, error) {
	data := make([]any, 0)
	m.data.Range(func(key, val any) bool {
		data = append(data, val)
		return true
	})
	return data, nil
}
