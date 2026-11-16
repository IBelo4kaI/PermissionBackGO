package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

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
func (q *sessionQueries) DeleteSessionsByUserID(_ context.Context, id string) error {
	q.deletedID = id
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
