package permission

import (
	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/nullable"
)

func fromModel(p repo.Permission) Permission {
	return Permission{
		ID:          p.ID,
		ServiceID:   nullable.StringOrNil(p.ServiceID),
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

func fromListRow(r repo.ListPermissionsRow) Permission {
	return Permission{
		ID:          r.ID,
		ServiceID:   nullable.StringOrNil(r.ServiceID),
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		ServiceName: nullable.StringOr(r.ServiceName, ""),
	}
}

func fromByServiceRow(r repo.ListPermissionsByServiceIDRow) Permission {
	return Permission{
		ID:          r.ID,
		ServiceID:   nullable.StringOrNil(r.ServiceID),
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		ServiceName: nullable.StringOr(r.ServiceName, ""),
	}
}

func fromGetByIDRow(r repo.GetPermissionByIDRow) Permission {
	return Permission{
		ID:          r.ID,
		ServiceID:   nullable.StringOrNil(r.ServiceID),
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		ServiceName: nullable.StringOr(r.ServiceName, ""),
	}
}

func fromGetByCodeRow(r repo.GetPermissionByCodeRow) Permission {
	return Permission{
		ID:          r.ID,
		ServiceID:   nullable.StringOrNil(r.ServiceID),
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		ServiceName: nullable.StringOr(r.ServiceName, ""),
	}
}
