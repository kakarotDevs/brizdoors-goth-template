package models

import (
	"errors"
	"sync"
	"time"
)

type User struct {
	ID        string // usually UUID or Google sub
	Name      string
	Email     string
	Picture   string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	userByID    = map[string]User{}
	userByEmail = map[string]User{}
	mu          sync.Mutex
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

func CreateUser(user *User) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := userByEmail[user.Email]; exists {
		return ErrUserAlreadyExists
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	userByID[user.ID] = *user
	userByEmail[user.Email] = *user

	return nil
}

// CreateOrUpdateUser creates or updates a user in memory.
func CreateOrUpdateUser(user User) {
	mu.Lock()
	defer mu.Unlock()

	user.UpdatedAt = time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = user.UpdatedAt
	}

	userByID[user.ID] = user
	userByEmail[user.Email] = user
}

func GetUserByEmail(email string) (User, error) {
	mu.Lock()
	defer mu.Unlock()

	user, ok := userByEmail[email]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func GetUserByID(id string) (User, error) {
	mu.Lock()
	defer mu.Unlock()

	user, ok := userByID[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}
