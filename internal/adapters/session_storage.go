package adapters

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

type sessionPersistence struct {
	sessions sync.Map
}

func NewSessionPersistence() SessionPersistence {
	return &sessionPersistence{
		sessions: sync.Map{},
	}
}

func (s *sessionPersistence) GetSessionUserID(sessionID string) (uuid.UUID, error) {
	userID, ok := s.sessions.Load(sessionID)
	if !ok {
		return uuid.Nil, errors.New("user session not found")
	}

	return userID.(uuid.UUID), nil
}

func (s *sessionPersistence) CreateSession(userID uuid.UUID) (string, error) {
	sessionID := uuid.New()

	s.sessions.Store(sessionID.String(), userID)

	return sessionID.String(), nil
}
