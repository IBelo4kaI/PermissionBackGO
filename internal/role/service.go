package role

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"time"

	pageHelper "permisson/internal/pkg/page"

	"github.com/google/uuid"

	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/nullable"
)

// detailedTimeLayout — аналог Python strftime("%d.%m.%Y %H:%M") в get_role_detailed.
const detailedTimeLayout = "02.01.2006 15:04"

type Service struct {
	queries repo.Querier
}

func NewService(queries repo.Querier) *Service {
	return &Service{queries: queries}
}

// List — аналог RoleService.get_all.
func (s *Service) List(ctx context.Context, page, limit int) (ListResponse, error) {
	page, limit = pageHelper.NormalizePage(page, limit)
	offset := (page - 1) * limit

	rows, err := s.queries.ListRolesWithCounts(ctx, repo.ListRolesWithCountsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return ListResponse{}, err
	}

	total, err := s.queries.CountRoles(ctx)
	if err != nil {
		return ListResponse{}, err
	}

	items := make([]RoleResponse, 0, len(rows))
	for _, r := range rows {
		items = append(items, fromListWithCountsRow(r))
	}

	return ListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pageHelper.PageCount(total, limit),
	}, nil
}

// ListByServiceID — аналог RoleService.get_all_by_service_id.
func (s *Service) ListByServiceID(ctx context.Context, serviceID string, page, limit int) (ListResponse, error) {
	page, limit = pageHelper.NormalizePage(page, limit)
	offset := (page - 1) * limit

	rows, err := s.queries.ListRolesWithCountsByServiceID(ctx, repo.ListRolesWithCountsByServiceIDParams{
		ServiceID: nullable.String(serviceID),
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return ListResponse{}, err
	}

	total, err := s.queries.CountRolesByServiceID(ctx, nullable.String(serviceID))
	if err != nil {
		return ListResponse{}, err
	}

	items := make([]RoleResponse, 0, len(rows))
	for _, r := range rows {
		items = append(items, fromByServiceWithCountsRow(r))
	}

	return ListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pageHelper.PageCount(total, limit),
	}, nil
}

// GetByID — вспомогательный метод (аналог repo.get_by_id), используется
// Create/Update/AddPermission/RemovePermission/Delete для проверки
// существования роли и сборки ответа после операции.
func (s *Service) GetByID(ctx context.Context, id string) (RoleResponse, error) {
	row, err := s.queries.GetRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RoleResponse{}, ErrNotFound
		}
		return RoleResponse{}, err
	}
	return fromModel(row), nil
}

// Create — аналог RoleService.create.
func (s *Service) Create(ctx context.Context, req UpsertRequest) (RoleResponse, error) {
	id := uuid.NewString()

	err := s.queries.CreateRole(ctx, repo.CreateRoleParams{
		ID:          id,
		ServiceID:   nullable.StringPtr(req.ServiceID),
		Name:        req.Name,
		Description: req.Description,
		IsGlobal:    req.IsGlobal,
	})
	if err != nil {
		return RoleResponse{}, err
	}

	return s.GetByID(ctx, id)
}

// Update — аналог RoleService.update.
func (s *Service) Update(ctx context.Context, id string, req UpsertRequest) (RoleResponse, error) {
	if _, err := s.GetByID(ctx, id); err != nil {
		return RoleResponse{}, err
	}

	err := s.queries.UpdateRole(ctx, repo.UpdateRoleParams{
		ServiceID:   nullable.StringPtr(req.ServiceID),
		Name:        req.Name,
		Description: req.Description,
		IsGlobal:    req.IsGlobal,
		ID:          id,
	})
	if err != nil {
		return RoleResponse{}, err
	}

	return s.GetByID(ctx, id)
}

// Delete — аналог RoleService.delete.
func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	return s.queries.DeleteRole(ctx, id)
}

// AddPermission — аналог RoleService.permission_add.
func (s *Service) AddPermission(ctx context.Context, roleID, permID string) (RoleResponse, error) {
	if _, err := s.GetByID(ctx, roleID); err != nil {
		return RoleResponse{}, err
	}

	if _, err := s.queries.GetPermissionByID(ctx, permID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RoleResponse{}, ErrPermissionNotFound
		}
		return RoleResponse{}, err
	}

	assigned, err := s.hasPermission(ctx, roleID, permID)
	if err != nil {
		return RoleResponse{}, err
	}
	if assigned {
		return RoleResponse{}, ErrPermissionAlreadyAssigned
	}

	if err := s.queries.AddPermissionToRole(ctx, repo.AddPermissionToRoleParams{
		RoleID:       roleID,
		PermissionID: permID,
	}); err != nil {
		return RoleResponse{}, err
	}

	return s.GetByID(ctx, roleID)
}

// RemovePermission — аналог RoleService.permission_remove.
func (s *Service) RemovePermission(ctx context.Context, roleID, permID string) (RoleResponse, error) {
	if _, err := s.GetByID(ctx, roleID); err != nil {
		return RoleResponse{}, err
	}

	if _, err := s.queries.GetPermissionByID(ctx, permID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RoleResponse{}, ErrPermissionNotFound
		}
		return RoleResponse{}, err
	}

	assigned, err := s.hasPermission(ctx, roleID, permID)
	if err != nil {
		return RoleResponse{}, err
	}
	if !assigned {
		return RoleResponse{}, ErrPermissionNotAssigned
	}

	if err := s.queries.RemovePermissionFromRole(ctx, repo.RemovePermissionFromRoleParams{
		RoleID:       roleID,
		PermissionID: permID,
	}); err != nil {
		return RoleResponse{}, err
	}

	return s.GetByID(ctx, roleID)
}

func (s *Service) hasPermission(ctx context.Context, roleID, permID string) (bool, error) {
	ids, err := s.queries.ListUsedPermissionIDsForRole(ctx, roleID)
	if err != nil {
		return false, err
	}
	if slices.Contains(ids, permID) {
		return true, nil
	}
	return false, nil
}

// Detailed — аналог RoleService.get_role_detailed +
// RoleRepository.get_role_with_all_permissions_info.
func (s *Service) Detailed(ctx context.Context, roleID string) (DetailedResponse, error) {
	roleRow, err := s.queries.GetRoleByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DetailedResponse{}, ErrNotFound
		}
		return DetailedResponse{}, err
	}

	allPerms, err := s.allPermissionsForRole(ctx, roleRow)
	if err != nil {
		return DetailedResponse{}, err
	}

	usedIDs, err := s.queries.ListUsedPermissionIDsForRole(ctx, roleID)
	if err != nil {
		return DetailedResponse{}, err
	}
	usedSet := make(map[string]struct{}, len(usedIDs))
	for _, id := range usedIDs {
		usedSet[id] = struct{}{}
	}

	usersRows, err := s.queries.ListUsersWithRole(ctx, roleID)
	if err != nil {
		return DetailedResponse{}, err
	}

	permissionsByService := make(map[string][]PermissionWithUse)
	usedCount := 0

	for _, p := range allPerms {
		_, used := usedSet[p.ID]
		if used {
			usedCount++
		}

		// "global" — тот же fallback-ключ, что и в Python для permissions
		// без service_id.
		key := "global"
		if p.ServiceID.Valid {
			key = p.ServiceID.String
		}

		permissionsByService[key] = append(permissionsByService[key], PermissionWithUse{
			ID:          p.ID,
			ServiceID:   nullable.StringOrNil(p.ServiceID),
			Code:        p.Code,
			Name:        p.Name,
			Description: p.Description,
			CreatedAt:   p.CreatedAt.Format(detailedTimeLayout),
			ServiceName: nullable.StringOrNil(p.ServiceName),
			Use:         used,
		})
	}

	users := make([]UserRoleInfo, 0, len(usersRows))
	for _, u := range usersRows {
		users = append(users, UserRoleInfo{
			ID:       u.ID,
			Name:     u.Name,
			Surname:  u.Surname,
			Username: u.Username,
		})
	}

	return DetailedResponse{
		ID:                   roleRow.ID,
		Name:                 roleRow.Name,
		Description:          roleRow.Description,
		IsGlobal:             roleRow.IsGlobal,
		CreatedAt:            roleRow.CreatedAt.Format(detailedTimeLayout),
		UsedPermissionsCount: usedCount,
		PermissionsByService: permissionsByService,
		Users:                users,
	}, nil
}

// allPermRow — общий промежуточный тип, чтобы Detailed() мог единообразно
// обработать результат двух разных Row-типов sqlc (ListAllPermissionsRow и
// ListAllPermissionsByServiceIDRow), у которых одинаковый набор колонок,
// но разные сгенерированные имена структур.
type allPermRow struct {
	ID          string
	ServiceID   sql.NullString
	Code        string
	Name        string
	Description string
	CreatedAt   time.Time
	ServiceName sql.NullString
}

// allPermissionsForRole — аналог ветвления в
// RoleRepository.get_role_with_all_permissions_info: если у роли задан
// service_id — берём разрешения только этого сервиса, иначе (глобальная
// роль) — вообще все разрешения в системе.
func (s *Service) allPermissionsForRole(ctx context.Context, roleRow repo.Role) ([]allPermRow, error) {
	if roleRow.ServiceID.Valid {
		rows, err := s.queries.ListAllPermissionsByServiceID(ctx, roleRow.ServiceID)
		if err != nil {
			return nil, err
		}
		out := make([]allPermRow, 0, len(rows))
		for _, r := range rows {
			out = append(out, allPermRow{
				ID: r.ID, ServiceID: r.ServiceID, Code: r.Code, Name: r.Name,
				Description: r.Description, CreatedAt: r.CreatedAt, ServiceName: r.ServiceName,
			})
		}
		return out, nil
	}

	rows, err := s.queries.ListAllPermissions(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]allPermRow, 0, len(rows))
	for _, r := range rows {
		out = append(out, allPermRow{
			ID: r.ID, ServiceID: r.ServiceID, Code: r.Code, Name: r.Name,
			Description: r.Description, CreatedAt: r.CreatedAt, ServiceName: r.ServiceName,
		})
	}
	return out, nil
}
