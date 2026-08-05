package role

import (
	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/nullable"
)

func fromModel(r repo.Role) RoleResponse {
	return RoleResponse{
		ID:          r.ID,
		ServiceID:   nullable.StringOrNil(r.ServiceID),
		Name:        r.Name,
		Description: r.Description,
		IsGlobal:    r.IsGlobal,
		CreatedAt:   r.CreatedAt,
	}
}

func fromListWithCountsRow(r repo.ListRolesWithCountsRow) RoleResponse {
	return RoleResponse{
		ID:               r.ID,
		ServiceID:        nullable.StringOrNil(r.ServiceID),
		Name:             r.Name,
		Description:      r.Description,
		IsGlobal:         r.IsGlobal,
		CreatedAt:        r.CreatedAt,
		UserCount:        r.UserCount,
		PermissionsCount: r.PermissionCount,
	}
}

func fromByServiceWithCountsRow(r repo.ListRolesWithCountsByServiceIDRow) RoleResponse {
	return RoleResponse{
		ID:               r.ID,
		ServiceID:        nullable.StringOrNil(r.ServiceID),
		Name:             r.Name,
		Description:      r.Description,
		IsGlobal:         r.IsGlobal,
		CreatedAt:        r.CreatedAt,
		UserCount:        r.UserCount,
		PermissionsCount: r.PermissionCount,
	}
}
