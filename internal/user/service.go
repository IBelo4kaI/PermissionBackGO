package user

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"

	"github.com/google/uuid"

	repo "permisson/internal/database/sqlc"
	"permisson/internal/permission"
	"permisson/internal/pkg/nullable"
	pageHelper "permisson/internal/pkg/page"
	"permisson/internal/pkg/passwordhash"
	"permisson/internal/pkg/response"
)

type Service struct {
	queries repo.Querier
}

// ListGenders — справочник полов для формы создания/редактирования
// пользователя. Без пагинации: справочник маленький и статичный.
func (s *Service) ListGenders(ctx context.Context) ([]repo.Gender, error) {
	return s.queries.ListGenders(ctx)
}

func NewService(queries repo.Querier) *Service {
	return &Service{queries: queries}
}

func (s *Service) List(ctx context.Context, page, limit int, search *string, sortBy, sortDir string) (*response.Page[UserResponse], error) {
	page, limit = pageHelper.NormalizePage(page, limit)
	offset := (page - 1) * limit

	rows, err := s.queries.ListUsers(ctx, repo.ListUsersParams{
		Search:  search,
		SortBy:  sortBy,
		SortDir: sortDir,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	total, err := s.queries.CountUsers(ctx, repo.CountUsersParams{Search: search})
	if err != nil {
		return nil, err
	}

	items := make([]UserResponse, 0, len(rows))
	for _, r := range rows {
		resp, err := s.buildResponse(ctx, fromListUsersRow(r))
		if err != nil {
			return nil, err
		}
		items = append(items, resp)
	}

	return &response.Page[UserResponse]{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pageHelper.PageCount(total, limit),
	}, nil
}

func (s *Service) ListByServiceID(ctx context.Context, serviceID string, page, limit int, search *string, sortBy, sortDir string) (*response.Page[UserResponse], error) {
	page, limit = pageHelper.NormalizePage(page, limit)
	offset := (page - 1) * limit

	rows, err := s.queries.ListUsersByServiceID(ctx, repo.ListUsersByServiceIDParams{
		ServiceID: nullable.String(serviceID),
		Search:    search,
		SortBy:    sortBy,
		SortDir:   sortDir,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, err
	}

	total, err := s.queries.CountUsersByServiceID(ctx, repo.CountUsersByServiceIDParams{
		ServiceID: nullable.String(serviceID),
		Search:    search,
	})
	if err != nil {
		return nil, err
	}

	items := make([]UserResponse, 0, len(rows))
	for _, r := range rows {
		resp, err := s.buildResponseByService(ctx, fromListUsersByServiceIDRow(r), serviceID)
		if err != nil {
			return nil, err
		}
		items = append(items, resp)
	}

	return &response.Page[UserResponse]{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pageHelper.PageCount(total, limit),
	}, nil
}

// ListAll: список без пагинации (весь users), поэтому search/sort применяются
// в Go до сборки ответа — так и фильтрация корректна (нет LIMIT/OFFSET, с
// которым её нужно было бы согласовывать), и N+1 запросы gender/roles в
// buildResponse не тратятся на строки, которые всё равно отфильтруются.
func (s *Service) ListAll(ctx context.Context, search *string, sortBy, sortDir string) ([]UserResponse, error) {
	rows, err := s.queries.ListAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	all := make([]userRow, 0, len(rows))
	for _, r := range rows {
		all = append(all, fromListAllUsersRow(r))
	}

	all = filterUserRows(all, search)
	sortUserRows(all, sortBy, sortDir)

	items := make([]UserResponse, 0, len(all))
	for _, u := range all {
		resp, err := s.buildResponse(ctx, u)
		if err != nil {
			return nil, err
		}
		items = append(items, resp)
	}
	return items, nil
}

// GetByID — базовый метод: используется остальными операциями для проверки
// существования пользователя и сборки ответа после мутации.
func (s *Service) GetByID(ctx context.Context, id string) (UserResponse, error) {
	row, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserResponse{}, ErrNotFound
		}
		return UserResponse{}, err
	}
	return s.buildResponse(ctx, fromGetUserRow(row))
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (UserResponse, error) {
	if _, err := s.queries.GetGenderByID(ctx, req.GenderID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserResponse{}, ErrGenderNotFound
		}
		return UserResponse{}, err
	}

	if err := s.ensureUsernameIsFree(ctx, req.Username, ""); err != nil {
		return UserResponse{}, err
	}

	password, err := passwordhash.Hash(req.Password)
	if err != nil {
		return UserResponse{}, err
	}

	id := uuid.NewString()
	err = s.queries.CreateUser(ctx, repo.CreateUserParams{
		ID:         id,
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: nullable.StringPtr(&req.Patronymic),
		Username:   req.Username,
		GenderID:   req.GenderID,
		Birthday:   req.Birthday,
		Password:   password,
		Status:     repo.UsersStatusActive,
	})
	if err != nil {
		return UserResponse{}, err
	}

	return s.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (UserResponse, error) {
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}

	if req.GenderID != nil {
		if _, err := s.queries.GetGenderByID(ctx, *req.GenderID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return UserResponse{}, ErrGenderNotFound
			}
			return UserResponse{}, err
		}
	}

	if req.Username != nil && *req.Username != existing.Username {
		if err := s.ensureUsernameIsFree(ctx, *req.Username, id); err != nil {
			return UserResponse{}, err
		}
	}

	var password sql.NullString
	if req.Password != nil {
		hash, err := passwordhash.Hash(*req.Password)
		if err != nil {
			return UserResponse{}, err
		}
		password = nullable.String(hash)
	}

	err = s.queries.UpdateUser(ctx, repo.UpdateUserParams{
		Name:       nullable.StringPtr(req.Name),
		Surname:    nullable.StringPtr(req.Surname),
		Patronymic: nullable.StringPtr(req.Patronymic),
		Username:   nullable.StringPtr(req.Username),
		Birthday:   nullable.TimePtr(req.Birthday),
		Status:     nullableStatusPtr(req.Status),
		GenderID:   nullable.StringPtr(req.GenderID),
		Password:   password,
		ID:         id,
	})
	if err != nil {
		return UserResponse{}, err
	}

	return s.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}

	// Сначала удаляем сессии пользователя, иначе FK (sessions.user_id -> users.id) блокирует удаление.
	if err := s.queries.DeleteSessionsByUserID(ctx, id); err != nil {
		return err
	}
	return s.queries.DeleteUser(ctx, id)
}

func (s *Service) AddRole(ctx context.Context, userID, roleID string) (UserResponse, error) {
	if _, err := s.GetByID(ctx, userID); err != nil {
		return UserResponse{}, err
	}

	if _, err := s.queries.GetRoleByID(ctx, roleID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserResponse{}, ErrRoleNotFound
		}
		return UserResponse{}, err
	}

	roles, err := s.queries.ListRolesForUser(ctx, userID)
	if err != nil {
		return UserResponse{}, err
	}
	for _, r := range roles {
		if r.ID == roleID {
			return UserResponse{}, ErrRoleAlreadyAssigned
		}
	}

	if err := s.queries.AddRoleToUser(ctx, repo.AddRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	}); err != nil {
		return UserResponse{}, err
	}

	return s.GetByID(ctx, userID)
}

func (s *Service) RemoveRole(ctx context.Context, userID, roleID string) (UserResponse, error) {
	if _, err := s.GetByID(ctx, userID); err != nil {
		return UserResponse{}, err
	}

	if _, err := s.queries.GetRoleByID(ctx, roleID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserResponse{}, ErrRoleNotFound
		}
		return UserResponse{}, err
	}

	roles, err := s.queries.ListRolesForUser(ctx, userID)
	if err != nil {
		return UserResponse{}, err
	}
	assigned := false
	for _, r := range roles {
		if r.ID == roleID {
			assigned = true
			break
		}
	}
	if !assigned {
		return UserResponse{}, ErrRoleNotAssigned
	}

	if err := s.queries.RemoveRoleFromUser(ctx, repo.RemoveRoleFromUserParams{
		UserID: userID,
		RoleID: roleID,
	}); err != nil {
		return UserResponse{}, err
	}

	return s.GetByID(ctx, userID)
}

// MePermissions — разрешения текущего пользователя для сервиса, с учётом
// wildcard-кодов вида "all:all:all" и "<service_name>:all:all".
// Если сервис не найден — пустой список (как в Python-версии).
//
// search/sort — как в ListAll: запрос ListPermissionsByUserIDAndServiceID
// используется и авторизационным middleware (permission.ListForUserAndService),
// поэтому фильтр/сортировку не добавляем в сам SQL-запрос, а применяем в Go
// только здесь, к уже полученному (обычно небольшому) списку.
func (s *Service) MePermissions(ctx context.Context, userID, serviceID string, search *string, sortBy, sortDir string) ([]permission.Permission, error) {
	svc, err := s.queries.GetServiceByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []permission.Permission{}, nil
		}
		return nil, err
	}

	rows, err := s.queries.ListPermissionsByUserIDAndServiceID(ctx, repo.ListPermissionsByUserIDAndServiceIDParams{
		UserID:      userID,
		ServiceID:   nullable.String(serviceID),
		ServiceName: svc.Name,
	})
	if err != nil {
		return nil, err
	}

	items := make([]permission.Permission, 0, len(rows))
	for _, r := range rows {
		items = append(items, permission.Permission{
			ID:          r.ID,
			ServiceID:   nullable.StringOrNil(r.ServiceID),
			Code:        r.Code,
			Name:        r.Name,
			Description: r.Description,
			CreatedAt:   r.CreatedAt,
		})
	}

	items = filterPermissions(items, search)
	sortPermissions(items, sortBy, sortDir)

	return items, nil
}

// buildResponse собирает UserResponse: gender и роли подтягиваются отдельными
// запросами на каждого пользователя — как ленивая ORM-загрузка .gender/.roles
// в Python-версии (см. аналогичный приём в service.List).
func (s *Service) buildResponse(ctx context.Context, u userRow) (UserResponse, error) {
	gender, err := s.queries.GetGenderByID(ctx, u.GenderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserResponse{}, ErrGenderNotFound
		}
		return UserResponse{}, err
	}

	roles, err := s.queries.ListRolesForUser(ctx, u.ID)
	if err != nil {
		return UserResponse{}, err
	}

	return fromUserRow(u, gender, roles), nil
}

// buildResponseByService — как buildResponse, но оставляет у пользователя только
// роли указанного сервиса (Python: user_response.roles = [r for r in roles
// if r.service_id == service_id]).
func (s *Service) buildResponseByService(ctx context.Context, u userRow, serviceID string) (UserResponse, error) {
	gender, err := s.queries.GetGenderByID(ctx, u.GenderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserResponse{}, ErrGenderNotFound
		}
		return UserResponse{}, err
	}

	roles, err := s.queries.ListRolesForUser(ctx, u.ID)
	if err != nil {
		return UserResponse{}, err
	}

	filtered := make([]repo.Role, 0, len(roles))
	for _, r := range roles {
		if r.ServiceID.Valid && r.ServiceID.String == serviceID {
			filtered = append(filtered, r)
		}
	}

	return fromUserRow(u, gender, filtered), nil
}

// filterUserRows — регистронезависимый поиск по name/surname/patronymic/
// username. Используется там, где список пользователей не пагинируется в SQL
// (ListAll) — поэтому фильтрация безопасно выполняется в Go после выборки.
func filterUserRows(rows []userRow, search *string) []userRow {
	if search == nil {
		return rows
	}
	needle := strings.ToLower(*search)
	filtered := make([]userRow, 0, len(rows))
	for _, r := range rows {
		if strings.Contains(strings.ToLower(r.Name), needle) ||
			strings.Contains(strings.ToLower(r.Surname), needle) ||
			strings.Contains(strings.ToLower(r.Patronymic.String), needle) ||
			strings.Contains(strings.ToLower(r.Username), needle) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// sortUserRows сортирует rows "на месте" по одной из user.SortableColumns.
func sortUserRows(rows []userRow, sortBy, sortDir string) {
	less := func(i, j int) bool {
		switch sortBy {
		case "name":
			return rows[i].Name < rows[j].Name
		case "surname":
			return rows[i].Surname < rows[j].Surname
		case "patronymic":
			return rows[i].Patronymic.String < rows[j].Patronymic.String
		case "username":
			return rows[i].Username < rows[j].Username
		case "birthday":
			return rows[i].Birthday.Before(rows[j].Birthday)
		case "status":
			return rows[i].Status < rows[j].Status
		default: // "created_at" и любое неизвестное значение
			return rows[i].CreatedAt.Before(rows[j].CreatedAt)
		}
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if sortDir == "asc" {
			return less(i, j)
		}
		return less(j, i)
	})
}

// filterPermissions — регистронезависимый поиск по code/name/description.
// Используется в MePermissions, где список не пагинируется в SQL.
func filterPermissions(items []permission.Permission, search *string) []permission.Permission {
	if search == nil {
		return items
	}
	needle := strings.ToLower(*search)
	filtered := make([]permission.Permission, 0, len(items))
	for _, p := range items {
		if strings.Contains(strings.ToLower(p.Code), needle) ||
			strings.Contains(strings.ToLower(p.Name), needle) ||
			strings.Contains(strings.ToLower(p.Description), needle) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

// sortPermissions сортирует items "на месте" по code/name/description/created_at.
func sortPermissions(items []permission.Permission, sortBy, sortDir string) {
	less := func(i, j int) bool {
		switch sortBy {
		case "code":
			return items[i].Code < items[j].Code
		case "name":
			return items[i].Name < items[j].Name
		case "description":
			return items[i].Description < items[j].Description
		default: // "created_at" и любое неизвестное значение
			return items[i].CreatedAt.Before(items[j].CreatedAt)
		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		if sortDir == "asc" {
			return less(i, j)
		}
		return less(j, i)
	})
}
