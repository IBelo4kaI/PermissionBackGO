package auth

import (
	"context"
	"database/sql"
	"errors"

	"permisson/internal/pkg/caller"
	"permisson/internal/pkg/nullable"
	"permisson/internal/pkg/token"
)

// CallerFrom определяет идентичность запроса по значению cookie "session":
//  1. сначала пробуем пользовательскую сессию;
//  2. не нашли — пробуем API-ключ сервиса (хэшируется SHA-256, как в Python);
//  3. иначе — ErrInvalidSession. Пустой токен — ErrSessionTokenMissing.
func (s *Service) CallerFrom(ctx context.Context, rawToken string) (caller.Caller, error) {
	if rawToken == "" {
		return caller.Caller{}, ErrSessionTokenMissing
	}

	session, err := s.queries.GetSessionByTokenHash(ctx, token.Hash(rawToken))
	if err == nil {
		return caller.Caller{Type: caller.User, UserID: session.UserID}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return caller.Caller{}, err
	}

	service, err := s.queries.GetServiceByAPIKeyHash(ctx, nullable.String(token.Hash(rawToken)))
	if err == nil {
		return caller.Caller{Type: caller.Service, ServiceID: service.ID}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return caller.Caller{}, err
	}

	return caller.Caller{}, ErrInvalidSession
}
