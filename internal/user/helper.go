package user

import (
	"context"
	"database/sql"
	"errors"

	repo "permisson/internal/database/sqlc"
)

func (s *Service) ensureUsernameIsFree(ctx context.Context, username, excludeID string) error {
	existing, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	if existing.ID != excludeID {
		return ErrUsernameExists
	}
	return nil
}

func nullableStatusPtr(s *string) repo.NullUsersStatus {
	if s == nil {
		return repo.NullUsersStatus{}
	}
	return repo.NullUsersStatus{UsersStatus: repo.UsersStatus(*s), Valid: true}
}
