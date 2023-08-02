package data

import "errors"

type User struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	Balance            float64 `json:"balance"`
	VerificationStatus bool    `json:"verification_status"`
}

func (s *Store) CreateUser(user *User) error {
	s.userStorage[user.ID] = user
	s.maxID = user.ID
	return nil
}

func (s *Store) GetUsers() ([]User, error) {
	users := make([]User, 0)
	for _, v := range s.userStorage {
		users = append(users, *v)
	}
	return users, nil
}

func (s *Store) GetUserByID(userID int) (*User, error) {
	user, ok := s.userStorage[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
