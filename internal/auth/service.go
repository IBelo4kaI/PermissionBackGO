package auth

import (
	"context"
	"database/sql"
	"errors"
	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/passwordhash"
	"permisson/internal/pkg/token"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	queries    repo.Querier
	sessionTTL time.Duration
}

func NewService(queries repo.Querier, sessionTTL time.Duration) *Service {
	return &Service{
		queries:    queries,
		sessionTTL: sessionTTL,
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResult, error) {
	user, err := s.queries.GetUserByUsername(ctx, req.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, err
	}

	if !passwordhash.Verify(user.Password, req.Password) {
		return LoginResult{}, ErrInvalidCredentials
	}

	if user.Status != repo.UsersStatusActive {
		return LoginResult{}, ErrUserInactive
	}

	rawToken, err := token.Generate()
	if err != nil {
		return LoginResult{}, err
	}

	expiresAt := time.Now().UTC().Add(s.sessionTTL)

	err = s.queries.CreateSession(ctx, repo.CreateSessionParams{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: token.Hash(rawToken),
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: rawToken, ExpiresAt: expiresAt}, nil
}

func (s *Service) ValidateSession(ctx context.Context, rawToken string) (bool, error) {
	if rawToken == "" {
		return false, nil
	}

	_, err := s.queries.GetSessionByTokenHash(ctx, token.Hash(rawToken))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Service) SessionUserID(ctx context.Context, rawToken string) (string, error) {
	if rawToken == "" {
		return "", ErrInvalidSession
	}

	session, err := s.queries.GetSessionByTokenHash(ctx, token.Hash(rawToken))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidSession
		}
		return "", err
	}

	return session.UserID, nil
}
