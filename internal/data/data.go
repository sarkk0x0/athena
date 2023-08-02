package data

type Store struct {
	userStorage map[int]*User
	maxID       int
}

func NewStore() *Store {
	storage := make(map[int]*User)
	return &Store{
		userStorage: storage,
		maxID:       0,
	}
}

func (s *Store) GetNextID() int {
	return s.maxID + 1
}
