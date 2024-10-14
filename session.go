package accessgo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

type SessionService struct {
	sessions sync.Map
	cancel   context.CancelFunc
}

func NewSessionService(ctx context.Context) *SessionService {
	ctx, cancel := context.WithCancel(ctx)
	ss := &SessionService{
		cancel: cancel,
	}
	go ss.cleanupRoutine(ctx)
	return ss
}

func (s *SessionService) CreateSession(userID int, longTerm bool) (string, error) {
	sessionID := uuid.NewString()

	expirationTime := time.Now().Add(24 * time.Hour)
	if longTerm {
		expirationTime = time.Now().Add(30 * 24 * time.Hour)
	}

	session := Session{
		ID:         sessionID,
		UserID:     userID,
		CreatedAt:  time.Now(),
		ExpiresAt:  expirationTime,
		IsLongTerm: longTerm,
	}

	s.sessions.Store(sessionID, session)

	return sessionID, nil
}

func (s *SessionService) GetSession(sessionID string) (Session, error) {
	if sessionValue, ok := s.sessions.Load(sessionID); ok {
		session := sessionValue.(Session)
		if time.Now().After(session.ExpiresAt) {
			s.sessions.Delete(sessionID)
			return Session{}, errors.New("session expired")
		}
		return session, nil
	}
	return Session{}, errors.New("session not found")
}

func (s *SessionService) DeleteSession(sessionID string) {
	s.sessions.Delete(sessionID)
}

func (s *SessionService) ExtendSession(sessionID string) error {
	if sessionValue, ok := s.sessions.Load(sessionID); ok {
		session := sessionValue.(Session)
		if time.Now().After(session.ExpiresAt) {
			s.sessions.Delete(sessionID)
			return errors.New("session expired")
		}

		if session.IsLongTerm {
			session.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)
		} else {
			session.ExpiresAt = time.Now().Add(24 * time.Hour)
		}

		s.sessions.Store(sessionID, session)
		return nil
	}
	return errors.New("session not found")
}

func (s *SessionService) cleanupRoutine(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupExpiredSessions()
		case <-ctx.Done():
			return
		}
	}
}

func (s *SessionService) cleanupExpiredSessions() {
	s.sessions.Range(func(key, value interface{}) bool {
		session := value.(Session)
		if time.Now().After(session.ExpiresAt) {
			s.sessions.Delete(key)
		}
		return true
	})
}

func (s *SessionService) Stop() {
	s.cancel()
}
