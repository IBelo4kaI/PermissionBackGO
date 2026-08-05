package role

import "time"

// RoleResponse — аналог pydantic RoleResponse.
//
// service_name у Python-версии объявлен в схеме, но ни в get_all, ни в
// get_all_by_service_id никогда фактически не заполняется (role_dict
// собирается только из колонок таблицы roles) — так и остаётся тут
// незаполненным (omitempty), поведение сохранено 1:1, а не "починено".
//
// is_global — в Python это int (0/1), здесь bool (так его отдаёт sqlc для
// TINYINT/INT-колонки). Это единственное сознательное отступление от
// побайтовой копии контракта — если клиенту нужен именно int, скажи, вернём.
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

// UpsertRequest — аналог RoleCreate, используется и для создания, и для
// редактирования роли.
type UpsertRequest struct {
	ServiceID   *string `json:"service_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsGlobal    bool    `json:"is_global"`
}

// AddPermissionRequest — аналог RoleAddPermission, используется и для
// /perm/add, и для /perm/remove.
type AddPermissionRequest struct {
	RoleID string `json:"role_id"`
	PermID string `json:"perm_id"`
}

// ListByServiceRequest — page/limit (query) + service_id (path) для
// GET /roles/service/:service_id. Поля объявлены явно (не embedding
// apidoc.Pagination) — см. комментарий в permission.ListByServiceRequest.
type ListByServiceRequest struct {
	ServiceID string `uri:"service_id"`
	Page      int    `query:"page" validate:"omitempty,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

// ListResponse — аналог PageResponse[RoleResponse].
type ListResponse struct {
	Items []RoleResponse `json:"items"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Pages int            `json:"pages"`
}

// DeleteResponse — аналог {"message": "..."} из delete_role.
type DeleteResponse struct {
	Message string `json:"message"`
}

// --- детальная страница роли (RoleService.get_role_detailed) ---

// PermissionWithUse — аналог pydantic PermissionWithUse.
// CreatedAt — строка в формате "dd.mm.yyyy HH:MM", как в Python
// (там дата форматируется перед сериализацией, а не отдаётся как ISO 8601).
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

// UserRoleInfo — аналог pydantic UserRoleInfo.
type UserRoleInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
}

// DetailedResponse — аналог pydantic RoleDetailedResponse.
type DetailedResponse struct {
	ID                   string                          `json:"id"`
	Name                 string                          `json:"name"`
	Description          string                          `json:"description"`
	IsGlobal             bool                             `json:"is_global"`
	CreatedAt            string                           `json:"created_at"`
	UsedPermissionsCount int                              `json:"used_permissions_count"`
	PermissionsByService map[string][]PermissionWithUse   `json:"permissions_by_service"`
	Users                []UserRoleInfo                   `json:"users"`
}
