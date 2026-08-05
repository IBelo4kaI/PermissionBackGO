package role

import "time"

// service_name никогда не заполняется: в Python-версии поле есть в схеме,
// но не наполняется ни в одном списковом эндпоинте (поведение сохранено 1:1).
//
// is_global — в Python это int (0/1), здесь bool: так его отдаёт sqlc для
// TINYINT-колонки (единственное сознательное отступление от контракта).
type RoleResponse struct {
	ID               string    `json:"id"`
	ServiceID        *string   `json:"service_id,omitempty"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	IsGlobal         bool      `json:"is_global"`
	ServiceName      string    `json:"service_name,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UserCount        int64     `json:"user_count"`
	PermissionsCount int64     `json:"permissions_count"`
}

type UpsertRequest struct {
	ServiceID   *string `json:"service_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsGlobal    bool    `json:"is_global"`
}

type AddPermissionRequest struct {
	RoleID string `json:"role_id"`
	PermID string `json:"perm_id"`
}

type ListByServiceRequest struct {
	ServiceID string `uri:"service_id"`
	Page      int    `query:"page" validate:"omitempty,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

type ListResponse struct {
	Items []RoleResponse `json:"items"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Pages int            `json:"pages"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

// CreatedAt — строка в формате "dd.mm.yyyy HH:MM" (не ISO 8601), как в Python-версии.
type PermissionWithUse struct {
	ID          string  `json:"id"`
	ServiceID   *string `json:"service_id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	ServiceName *string `json:"service_name"`
	Use         bool    `json:"use"`
}

type UserRoleInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
}

type DetailedResponse struct {
	ID                   string                         `json:"id"`
	Name                 string                         `json:"name"`
	Description          string                         `json:"description"`
	IsGlobal             bool                           `json:"is_global"`
	CreatedAt            string                         `json:"created_at"`
	UsedPermissionsCount int                            `json:"used_permissions_count"`
	PermissionsByService map[string][]PermissionWithUse `json:"permissions_by_service"`
	Users                []UserRoleInfo                 `json:"users"`
}
