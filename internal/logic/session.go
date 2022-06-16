package logic

import (
	"github.com/google/uuid"
	"random-numbers/internal/adapters"
)

type sessionWorker struct {
	sessions adapters.SessionPersistence
}

func NewSessionWorker(sessions adapters.SessionPersistence) SessionWorker {
	return &sessionWorker{
		sessions: sessions,
	}
}

func (s *sessionWorker) CreateSession(userID uuid.UUID) (string, error) {
	return s.sessions.CreateSession(userID)
}

func (s *sessionWorker) GetUserID(sessionID string) (uuid.UUID, error) {
	return s.sessions.GetSessionUserID(sessionID)
}

func (s *sessionWorker) IsAuthenticated(sessionID string) error {
	_, err := s.sessions.GetSessionUserID(sessionID)

	return err
}
