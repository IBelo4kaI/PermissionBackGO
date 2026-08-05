package service

import "time"

// ServiceResponse — аналог pydantic ServiceResponse. api_key_hash сознательно
// не входит в структуру — наружу он никогда не отдаётся.
//
// Name и ServiceName дублируют друг друга намеренно: в Python ServiceResponse
// объявляет service_name как alias="name" поверх унаследованного поля name,
// поэтому в JSON-ответе всегда были оба ключа с одинаковым значением.
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

// UpsertRequest — аналог ServiceCreate, используется и для создания, и для
// редактирования сервиса (Python переиспользует одну и ту же модель).
type UpsertRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    *string `json:"image_url"`
	URL         *string `json:"url"`
	Theme       *string `json:"theme"`
	Prefix      string  `json:"prefix"`
}

// APIKeyResponse — аналог ServiceApiKeyResponse. api_key отдаётся ровно один
// раз, в момент выпуска/перевыпуска — дальше в базе хранится только его хэш.
type APIKeyResponse struct {
	ServiceID string `json:"service_id"`
	APIKey    string `json:"api_key"`
}

// AccessResponse — аналог ServiceAccessResponse (для /services/user-accessible).
type AccessResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    *string `json:"image_url,omitempty"`
	URL         *string `json:"url,omitempty"`
	Theme       *string `json:"theme,omitempty"`
}

// ListResponse — аналог PageResponse[ServiceResponse].
type ListResponse struct {
	Items []ServiceResponse `json:"items"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
	Pages int               `json:"pages"`
}
