package permission

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Service) ensureCodeIsFree(ctx context.Context, code, excludeID string) error {
	existing, err := s.queries.GetPermissionByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	if existing.ID != excludeID {
		return ErrCodeExists
	}
	return nil
}
