package permission

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/nullable"
	pageHelper "permisson/internal/pkg/page"
	"permisson/internal/pkg/response"

	"github.com/google/uuid"
)

type Service struct {
	queries repo.Querier
}

func NewService(queries repo.Querier) *Service {
	return &Service{
		queries: queries,
	}
}

func (s *Service) List(ctx context.Context, page, limit int) (*response.Page[Permission], error) {
	page, limit = pageHelper.NormalizePage(page, limit)
	offset := (page - 1) * limit

	rows, err := s.queries.ListPermissions(ctx, repo.ListPermissionsParams{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}

	total, err := s.queries.CountPermissions(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]Permission, 0, len(rows))
	for _, r := range rows {
		items = append(items, fromListRow(r))
	}

	fmt.Printf("%v\n", rows)

	return &response.Page[Permission]{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pageHelper.PageCount(total, limit),
	}, nil
}

func (s *Service) ListByServiceID(ctx context.Context, serviceID string, page, limit int) ([]Permission, error) {
	page, limit = pageHelper.NormalizePage(page, limit)
	offset := (page - 1) * limit

	rows, err := s.queries.ListPermissionsByServiceID(ctx, repo.ListPermissionsByServiceIDParams{
		ServiceID: nullable.String(serviceID),
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, err
	}

	items := make([]Permission, 0, len(rows))
	for _, r := range rows {
		items = append(items, fromByServiceRow(r))
	}
	return items, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (Permission, error) {
	row, err := s.queries.GetPermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Permission{}, ErrNotFound
		}
		return Permission{}, err
	}
	return fromGetByIDRow(row), nil
}

func (s *Service) GetByCode(ctx context.Context, code string) (Permission, error) {
	row, err := s.queries.GetPermissionByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Permission{}, ErrNotFound
		}
		return Permission{}, err
	}
	return fromGetByCodeRow(row), nil
}

// Используется другими сущностями (например, role), поэтому не привязан к HTTP-роуту.
func (s *Service) ListForUser(ctx context.Context, userID string) ([]Permission, error) {
	rows, err := s.queries.ListPermissionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	items := make([]Permission, 0, len(rows))
	for _, r := range rows {
		items = append(items, fromModel(r))
	}
	return items, nil
}

// serviceName нужен отдельно от serviceID, т.к. wildcard-коды матчатся по имени
// сервиса ("shop:all:all"), а не по его ID — вызывающая сторона (например, middleware)
// должна сама подтянуть имя сервиса перед вызовом.
func (s *Service) ListForUserAndService(ctx context.Context, userID, serviceID, serviceName string) ([]Permission, error) {
	rows, err := s.queries.ListPermissionsByUserIDAndServiceID(ctx, repo.ListPermissionsByUserIDAndServiceIDParams{
		UserID:      userID,
		ServiceID:   nullable.String(serviceID),
		ServiceName: serviceName,
	})
	if err != nil {
		return nil, err
	}
	items := make([]Permission, 0, len(rows))
	for _, r := range rows {
		items = append(items, fromModel(r))
	}
	return items, nil
}

// Основа для будущего middleware.RequirePermission(entity, action).
func (s *Service) ExistsForUser(ctx context.Context, userID, service, entity, action string) (bool, error) {
	return s.queries.ExistsUserPermission(ctx, repo.ExistsUserPermissionParams{
		UserID:  userID,
		Service: service,
		Entity:  entity,
		Action:  action,
	})
}

func (s *Service) Create(ctx context.Context, req UpsertRequest) (Permission, error) {
	if err := s.ensureCodeIsFree(ctx, req.Code, ""); err != nil {
		return Permission{}, err
	}

	id := uuid.NewString()
	err := s.queries.CreatePermission(ctx, repo.CreatePermissionParams{
		ID:          id,
		ServiceID:   nullable.StringPtr(req.ServiceID),
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return Permission{}, err
	}

	return s.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, req UpsertRequest) (Permission, error) {
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return Permission{}, err
	}

	if req.Code != existing.Code {
		if err := s.ensureCodeIsFree(ctx, req.Code, id); err != nil {
			return Permission{}, err
		}
	}

	err = s.queries.UpdatePermission(ctx, repo.UpdatePermissionParams{
		ServiceID:   nullable.StringPtr(req.ServiceID),
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		ID:          id,
	})
	if err != nil {
		return Permission{}, err
	}

	return s.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	return s.queries.DeletePermission(ctx, id)
}
