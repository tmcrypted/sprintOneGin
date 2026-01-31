package service

import (
	"errors"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserService interface {
	CreateUser(name string, age int) (User, error)
	GetUser(id int) (User, error)
	GetAllUsers() []User
}
type inMemoryUserService struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

func NewUserService() UserService {
	return &inMemoryUserService{
		users:  make(map[int]User),
		nextID: 1,
	}
}

func (s *inMemoryUserService) CreateUser(name string, age int) (User, error) {
	if name == "" {
		return User{}, errors.New("name is required")
	}
	if age < 0 || age > 150 {
		return User{}, errors.New("age must be between 0 and 150")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	user := User{ID: s.nextID, Name: name, Age: age}
	s.users[s.nextID] = user
	s.nextID++
	return user, nil
}

func (s *inMemoryUserService) GetUser(id int) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (s *inMemoryUserService) GetAllUsers() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]User, 0, len(s.users))
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}
