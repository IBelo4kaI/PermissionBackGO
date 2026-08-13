package service

import "time"

// api_key_hash намеренно не входит в структуру — наружу он не отдаётся.
//
// Name и ServiceName дублируют друг друга намеренно: в JSON-ответе всегда
// присутствуют оба ключа с одинаковым значением (сохранено из Python-версии).
type ServiceResponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ServiceName      string    `json:"service_name"`
	Description      string    `json:"description"`
	ImageURL         *string   `json:"image_url,omitempty"`
	URL              *string   `json:"url,omitempty"`
	Theme            *string   `json:"theme,omitempty"`
	Prefix           string    `json:"prefix"`
	CreatedAt        time.Time `json:"created_at"`
	PermissionsCount int       `json:"permissions_count"`
}

type UpsertRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    *string `json:"image_url"`
	URL         *string `json:"url"`
	Theme       *string `json:"theme"`
	Prefix      string  `json:"prefix"`
}

// ServiceIDRequest — ID сервиса в path-параметре :service_id.
type ServiceIDRequest struct {
	ServiceID string `uri:"service_id"`
}

// UpdateServiceRequest — как UpsertRequest, но с ID сервиса в path-параметре.
type UpdateServiceRequest struct {
	ServiceID   string  `uri:"service_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    *string `json:"image_url"`
	URL         *string `json:"url"`
	Theme       *string `json:"theme"`
	Prefix      string  `json:"prefix"`
}

// api_key отдаётся ровно один раз, при выпуске/перевыпуске — дальше в БД
// хранится только его хэш.
type APIKeyResponse struct {
	ServiceID string `json:"service_id"`
	APIKey    string `json:"api_key"`
}

type AccessResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    *string `json:"image_url,omitempty"`
	URL         *string `json:"url,omitempty"`
	Theme       *string `json:"theme,omitempty"`
}

// ListUserAccessibleRequest — query-параметры GET /services/user-accessible.
// Список небольшой (сервисы, доступные конкретному пользователю) и без
// пагинации, поэтому поиск и сортировка применяются на стороне Go, а не в SQL
// (см. Service.ListAccessibleForUser).
type ListUserAccessibleRequest struct {
	Search  string `query:"search" validate:"omitempty,max=255" description:"Поиск по name/description"`
	SortBy  string `query:"sort_by" description:"name или description (по умолчанию name)"`
	SortDir string `query:"sort_dir" validate:"omitempty,oneof=asc desc" description:"asc или desc (по умолчанию desc)"`
}

// SortableColumns — белый список для sort_by в GET /services/ (см. query.QuerySort).
var SortableColumns = []string{"name", "description", "prefix", "created_at"}

const DefaultSortColumn = "created_at"

// AccessibleSortableColumns — белый список для sort_by в GET /services/user-accessible.
var AccessibleSortableColumns = []string{"name", "description"}

const AccessibleDefaultSortColumn = "name"
