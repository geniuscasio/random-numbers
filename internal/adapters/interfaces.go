package adapters

import (
	"github.com/google/uuid"
	"random-numbers/internal/common"
)

type NumbersGenerator interface {
	Get(start, stop, count int64) ([]int64, error)
}

type SessionPersistence interface {
	GetSessionUserID(sessionID string) (uuid.UUID, error)
	CreateSession(userID uuid.UUID) (string, error)
}

type UserPersistence interface {
	Create(common.User) error
	ByCreds(email, pwd string) (*common.User, error)
	ByID(userID uuid.UUID) (*common.User, error)
	Update(user common.User) error
}
