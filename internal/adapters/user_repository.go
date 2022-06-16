package adapters

import (
	"errors"
	"github.com/google/uuid"
	"random-numbers/internal/common"
	"sync"
)

type userPersistence struct {
	// using old maps with sync, as sync.Map
	// will make code complicated to read
	users map[uuid.UUID]common.User
	*sync.Mutex
}

func NewUserPersistence() UserPersistence {
	return &userPersistence{
		users: make(map[uuid.UUID]common.User),
		Mutex: &sync.Mutex{},
	}
}

func (u *userPersistence) Create(user common.User) error {
	u.Lock()
	defer u.Unlock()

	for _, v := range u.users {
		if v.Email == user.Email {
			return errors.New("user already exists")
		}
	}

	u.users[user.ID] = user

	return nil
}

func (u *userPersistence) ByCreds(email, pwd string) (*common.User, error) {
	u.Lock()
	defer u.Unlock()

	for _, user := range u.users {
		if user.Email == email && user.PasswordHash == pwd {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (u *userPersistence) ByID(userID uuid.UUID) (*common.User, error) {
	u.Lock()
	defer u.Unlock()

	user, ok := u.users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (u *userPersistence) Update(user common.User) error {
	u.Lock()
	defer u.Unlock()

	_, found := u.users[user.ID]

	if !found {
		return errors.New("not found")
	}

	u.users[user.ID] = user

	return nil
}
