package logic

import (
	"github.com/google/uuid"
	"random-numbers/internal/common"
)

type SessionWorker interface {
	CreateSession(userID uuid.UUID) (string, error)
	GetUserID(sessionID string) (uuid.UUID, error)
	IsAuthenticated(sessionID string) error
}

type Service interface {
	Generate(from, to, count int64, orderDesc bool) (*common.GenerateResponse, error)
	GetStatistic(userID uuid.UUID) (common.DetailsResponse, error)
	IncrementStatistic(userID uuid.UUID) error
	CreateUser(user common.User) error
	LoginUser(email string, pwd string) (uuid.UUID, error)
}
