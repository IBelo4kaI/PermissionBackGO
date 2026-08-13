package role

import (
	"permisson/internal/pkg/apidoc"
	"time"
)

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

// RoleIDRequest — ID роли в path-параметре :role_id.
type RoleIDRequest struct {
	RoleID string `uri:"role_id"`
}

// UpdateRoleRequest — как UpsertRequest, но с ID роли в path-параметре.
type UpdateRoleRequest struct {
	RoleID      string  `uri:"role_id"`
	ServiceID   *string `json:"service_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsGlobal    bool    `json:"is_global"`
}

// ListRequest — query-параметры GET /roles/. Расширяет apidoc.Pagination
// параметром is_global — фильтром "глобальная / привязанная к сервису роль"
// (не задан — без фильтра). Биндинг всё так же вручную в хендлере
// (query.QueryBoolPtr), тег нужен только для документации OpenAPI.
type ListRequest struct {
	apidoc.Pagination
	IsGlobal *bool `query:"is_global" description:"true — только глобальные роли, false — только привязанные к сервису"`
}

type AddPermissionRequest struct {
	RoleID string `json:"role_id"`
	PermID string `json:"perm_id"`
}

type ListByServiceRequest struct {
	ServiceID string `uri:"service_id"`
	Page      int    `query:"page" validate:"omitempty,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search    string `query:"search" validate:"omitempty,max=255" description:"Поиск по name/description"`
	SortBy    string `query:"sort_by" description:"name, description или created_at (по умолчанию created_at)"`
	SortDir   string `query:"sort_dir" validate:"omitempty,oneof=asc desc" description:"asc или desc (по умолчанию desc)"`
}

// SortableColumns — белый список колонок для query-параметра sort_by
// (см. query.QuerySort). service_name сюда не входит: в списковых
// эндпоинтах он никогда не заполняется (см. RoleResponse).
var SortableColumns = []string{"name", "description", "created_at"}

const DefaultSortColumn = "created_at"

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
	ServiceID            *string                        `json:"service_id"`
	Name                 string                         `json:"name"`
	Description          string                         `json:"description"`
	IsGlobal             bool                           `json:"is_global"`
	CreatedAt            string                         `json:"created_at"`
	UsedPermissionsCount int                            `json:"used_permissions_count"`
	PermissionsByService map[string][]PermissionWithUse `json:"permissions_by_service"`
	Users                []UserRoleInfo                 `json:"users"`
}
