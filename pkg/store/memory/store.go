package memory

type Store struct{}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Set(key string, value interface{}) error {
	return nil
}

func (s *Store) Get(key string) (interface{}, error) {
	return nil, nil
}
