package auth

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/token"
)

type validationQueries struct {
	repo.Querier
	sessionErr    error
	apiKeyErr     error
	sessionHashes []string
	apiKeyHashes  []sql.NullString
}

func (q *validationQueries) GetSessionByTokenHash(_ context.Context, hash string) (repo.GetSessionByTokenHashRow, error) {
	q.sessionHashes = append(q.sessionHashes, hash)
	return repo.GetSessionByTokenHashRow{}, q.sessionErr
}

func (q *validationQueries) GetServiceByAPIKeyHash(_ context.Context, hash sql.NullString) (repo.GetServiceByAPIKeyHashRow, error) {
	q.apiKeyHashes = append(q.apiKeyHashes, hash)
	return repo.GetServiceByAPIKeyHashRow{}, q.apiKeyErr
}

func TestValidateSession(t *testing.T) {
	dbErr := errors.New("database unavailable")
	for _, tc := range []struct {
		name                              string
		rawToken                          string
		sessionErr, apiKeyErr             error
		wantValid                         bool
		wantErr                           error
		wantSessionCalls, wantAPIKeyCalls int
	}{
		{name: "empty token"},
		{name: "valid session", rawToken: "session-token", wantValid: true, wantSessionCalls: 1},
		{name: "valid API key", rawToken: "api-key", sessionErr: sql.ErrNoRows, wantValid: true, wantSessionCalls: 1, wantAPIKeyCalls: 1},
		{name: "invalid token", rawToken: "unknown", sessionErr: sql.ErrNoRows, apiKeyErr: sql.ErrNoRows, wantSessionCalls: 1, wantAPIKeyCalls: 1},
		{name: "session database error", rawToken: "session-token", sessionErr: dbErr, wantErr: dbErr, wantSessionCalls: 1},
		{name: "API key database error", rawToken: "api-key", sessionErr: sql.ErrNoRows, apiKeyErr: dbErr, wantErr: dbErr, wantSessionCalls: 1, wantAPIKeyCalls: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			queries := &validationQueries{sessionErr: tc.sessionErr, apiKeyErr: tc.apiKeyErr}
			s := NewService(queries, 0, nil, "")
			valid, err := s.ValidateSession(context.Background(), tc.rawToken)
			if valid != tc.wantValid || !errors.Is(err, tc.wantErr) {
				t.Fatalf("ValidateSession() = (%v, %v), want (%v, %v)", valid, err, tc.wantValid, tc.wantErr)
			}
			if len(queries.sessionHashes) != tc.wantSessionCalls || len(queries.apiKeyHashes) != tc.wantAPIKeyCalls {
				t.Fatalf("unexpected queries: sessions=%d, API keys=%d", len(queries.sessionHashes), len(queries.apiKeyHashes))
			}
			for _, hash := range queries.sessionHashes {
				if hash != token.Hash(tc.rawToken) {
					t.Fatal("session lookup must use the token hash")
				}
			}
			for _, hash := range queries.apiKeyHashes {
				if !hash.Valid || hash.String != token.Hash(tc.rawToken) {
					t.Fatal("API key lookup must use the token hash")
				}
			}
		})
	}
}
