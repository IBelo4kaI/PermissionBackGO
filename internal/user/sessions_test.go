package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	repo "permisson/internal/database/sqlc"
)

type sessionQueries struct {
	repo.Querier
	lookupErr error
	deleteErr error
	deletedID string
}

func (q *sessionQueries) GetUserByID(_ context.Context, _ string) (repo.GetUserByIDRow, error) {
	return repo.GetUserByIDRow{}, q.lookupErr
}
func (q *sessionQueries) DeleteSessionsByUserIDBefore(_ context.Context, arg repo.DeleteSessionsByUserIDBeforeParams) error {
	q.deletedID = arg.UserID
	return q.deleteErr
}

func TestDeleteSessions(t *testing.T) {
	dbErr := errors.New("database error")
	for _, tc := range []struct {
		name                   string
		lookup, deletion, want error
		deleted                bool
	}{
		{name: "existing user", deleted: true},
		{name: "missing user", lookup: sql.ErrNoRows, want: ErrNotFound},
		{name: "lookup failure", lookup: dbErr, want: dbErr},
		{name: "delete failure", deletion: dbErr, want: dbErr, deleted: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			q := &sessionQueries{lookupErr: tc.lookup, deleteErr: tc.deletion}
			s := &Service{queries: q}
			err := s.DeleteSessions(context.Background(), "user-1")
			if !errors.Is(err, tc.want) {
				t.Fatalf("error = %v, want %v", err, tc.want)
			}
			if (q.deletedID == "user-1") != tc.deleted {
				t.Fatalf("deleted ID = %q", q.deletedID)
			}
		})
	}
}

func TestSessionExpiryCutoff(t *testing.T) {
	for _, tc := range []struct{ now, want string }{
		{"2026-11-16T12:00:00Z", "2026-11-17T21:00:00Z"},
		{"2026-11-16T22:00:00Z", "2026-11-18T21:00:00Z"},
	} {
		now, _ := time.Parse(time.RFC3339, tc.now)
		if got := sessionExpiryCutoff(now).Format(time.RFC3339); got != tc.want {
			t.Errorf("cutoff for %s = %s, want %s", tc.now, got, tc.want)
		}
	}
}
