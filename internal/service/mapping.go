package service

import (
	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/nullable"
)

func fromListRow(r repo.ListServicesRow, permissionsCount int64) ServiceResponse {
	return ServiceResponse{
		ID:               r.ID,
		Name:             r.Name,
		ServiceName:      r.Name,
		Description:      r.Description,
		ImageURL:         nullable.StringOrNil(r.ImageUrl),
		URL:              nullable.StringOrNil(r.Url),
		Theme:            nullable.StringOrNil(r.Theme),
		Prefix:           r.Prefix,
		CreatedAt:        r.CreatedAt,
		PermissionsCount: int(permissionsCount),
	}
}

func fromGetByIDRow(r repo.GetServiceByIDRow, permissionsCount int64) ServiceResponse {
	return ServiceResponse{
		ID:               r.ID,
		Name:             r.Name,
		ServiceName:      r.Name,
		Description:      r.Description,
		ImageURL:         nullable.StringOrNil(r.ImageUrl),
		URL:              nullable.StringOrNil(r.Url),
		Theme:            nullable.StringOrNil(r.Theme),
		Prefix:           r.Prefix,
		CreatedAt:        r.CreatedAt,
		PermissionsCount: int(permissionsCount),
	}
}
