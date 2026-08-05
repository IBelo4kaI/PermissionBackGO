package role

import (
	"context"
	"slices"
)

func (s *Service) hasPermission(ctx context.Context, roleID, permID string) (bool, error) {
	ids, err := s.queries.ListUsedPermissionIDsForRole(ctx, roleID)
	if err != nil {
		return false, err
	}
	return slices.Contains(ids, permID), nil
}
