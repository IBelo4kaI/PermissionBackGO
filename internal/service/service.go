package service

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"

	"github.com/google/uuid"

	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/nullable"
	pageHelper "permisson/internal/pkg/page"
	"permisson/internal/pkg/response"
	"permisson/internal/pkg/token"
)

type Service struct {
	queries repo.Querier
}

func NewService(queries repo.Querier) *Service {
	return &Service{queries: queries}
}

func (s *Service) List(ctx context.Context, page, limit int, search *string, sortBy, sortDir string) (*response.Page[ServiceResponse], error) {
	page, limit = pageHelper.NormalizePage(page, limit)
	offset := (page - 1) * limit

	rows, err := s.queries.ListServices(ctx, repo.ListServicesParams{
		Search:  search,
		SortBy:  sortBy,
		SortDir: sortDir,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	total, err := s.queries.CountServices(ctx, repo.CountServicesParams{Search: search})
	if err != nil {
		return nil, err
	}

	items := make([]ServiceResponse, 0, len(rows))
	for _, r := range rows {
		// permissions_count считается отдельным запросом на каждый элемент списка
		// (как ленивая ORM-загрузка .permissions в Python-версии).
		count, err := s.queries.CountPermissionsByServiceID(ctx, repo.CountPermissionsByServiceIDParams{ServiceID: nullable.String(r.ID)})
		if err != nil {
			return nil, err
		}
		items = append(items, fromListRow(r, count))
	}

	return &response.Page[ServiceResponse]{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pageHelper.PageCount(total, limit),
	}, nil
}

// GetByID — используется Create/Update/IssueAPIKey/RevokeAPIKey, чтобы собрать
// полноценный ServiceResponse после операции и проверить существование сервиса.
func (s *Service) GetByID(ctx context.Context, id string) (ServiceResponse, error) {
	row, err := s.queries.GetServiceByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ServiceResponse{}, ErrNotFound
		}
		return ServiceResponse{}, err
	}

	count, err := s.queries.CountPermissionsByServiceID(ctx, repo.CountPermissionsByServiceIDParams{ServiceID: nullable.String(id)})
	if err != nil {
		return ServiceResponse{}, err
	}

	return fromGetByIDRow(row, count), nil
}

func (s *Service) Create(ctx context.Context, req UpsertRequest) (ServiceResponse, error) {
	id := uuid.NewString()

	err := s.queries.CreateService(ctx, repo.CreateServiceParams{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		ImageUrl:    nullable.StringPtr(req.ImageURL),
		Url:         nullable.StringPtr(req.URL),
		Theme:       nullable.StringPtr(req.Theme),
		Prefix:      req.Prefix,
	})
	if err != nil {
		return ServiceResponse{}, err
	}

	return s.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, req UpsertRequest) (ServiceResponse, error) {
	if _, err := s.GetByID(ctx, id); err != nil {
		return ServiceResponse{}, err
	}

	err := s.queries.UpdateService(ctx, repo.UpdateServiceParams{
		Name:        req.Name,
		Description: req.Description,
		ImageUrl:    nullable.StringPtr(req.ImageURL),
		Url:         nullable.StringPtr(req.URL),
		Theme:       nullable.StringPtr(req.Theme),
		Prefix:      req.Prefix,
		ID:          id,
	})
	if err != nil {
		return ServiceResponse{}, err
	}

	return s.GetByID(ctx, id)
}

// ListAccessibleForUser: список небольшой и без пагинации (сервисы,
// доступные конкретному пользователю), поэтому search/sort применяются
// на срезе в Go, а не в SQL — в отличие от List, здесь не нужно согласовывать
// фильтр с LIMIT/OFFSET.
func (s *Service) ListAccessibleForUser(ctx context.Context, userID string, search *string, sortBy, sortDir string) ([]AccessResponse, error) {
	rows, err := s.queries.ListServicesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	items := make([]AccessResponse, 0, len(rows))
	for _, r := range rows {
		items = append(items, AccessResponse{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			ImageURL:    nullable.StringOrNil(r.ImageUrl),
			URL:         nullable.StringOrNil(r.Url),
			Theme:       nullable.StringOrNil(r.Theme),
		})
	}

	if search != nil {
		needle := strings.ToLower(*search)
		filtered := make([]AccessResponse, 0, len(items))
		for _, it := range items {
			if strings.Contains(strings.ToLower(it.Name), needle) ||
				strings.Contains(strings.ToLower(it.Description), needle) {
				filtered = append(filtered, it)
			}
		}
		items = filtered
	}

	key := func(a AccessResponse) string {
		if sortBy == "description" {
			return a.Description
		}
		return a.Name
	}
	sort.SliceStable(items, func(i, j int) bool {
		if sortDir == "asc" {
			return key(items[i]) < key(items[j])
		}
		return key(items[i]) > key(items[j])
	})

	return items, nil
}

// Возвращает сырой ключ ровно один раз; перевыпуск затирает предыдущий ключ.
func (s *Service) IssueAPIKey(ctx context.Context, serviceID string) (APIKeyResponse, error) {
	if _, err := s.GetByID(ctx, serviceID); err != nil {
		return APIKeyResponse{}, err
	}

	rawKey, err := token.Generate()
	if err != nil {
		return APIKeyResponse{}, err
	}

	err = s.queries.SetServiceAPIKeyHash(ctx, repo.SetServiceAPIKeyHashParams{
		ApiKeyHash: nullable.String(token.Hash(rawKey)),
		ID:         serviceID,
	})
	if err != nil {
		return APIKeyResponse{}, err
	}

	return APIKeyResponse{ServiceID: serviceID, APIKey: rawKey}, nil
}

func (s *Service) RevokeAPIKey(ctx context.Context, serviceID string) error {
	if _, err := s.GetByID(ctx, serviceID); err != nil {
		return err
	}
	return s.queries.RevokeServiceAPIKey(ctx, serviceID)
}
