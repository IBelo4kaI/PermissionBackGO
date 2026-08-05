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

type ListResponse struct {
	Items []ServiceResponse `json:"items"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
	Pages int               `json:"pages"`
}
