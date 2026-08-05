package permission

import "time"

type Permission struct {
	ID          string    `json:"id"`
	ServiceID   *string   `json:"service_id,omitempty"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	ServiceName string    `json:"service_name"`
}

type UpsertRequest struct {
	ServiceID   *string `json:"service_id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
}

type ListByServiceRequest struct {
	ServiceID string `uri:"service_id"`
	Page      int    `query:"page" validate:"omitempty,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}
