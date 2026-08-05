package user

import (
	"time"

	repo "permisson/internal/database/sqlc"
	"permisson/internal/role"
)

// patronymic — *string (не string): колонка в БД nullable, и в Python-версии
// при NULL сериализуется как null, а не пустая строка.
//
// roles и roles_count повторяют Python-контракт: roles — полный RoleResponse,
// roles_count вычисляется как len(roles).
type UserResponse struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	Surname    string              `json:"surname"`
	Patronymic *string             `json:"patronymic"`
	Username   string              `json:"username"`
	Birthday   time.Time           `json:"birthday"`
	Status     repo.UsersStatus    `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
	Gender     repo.Gender         `json:"gender"`
	Roles      []role.RoleResponse `json:"roles"`
	RolesCount int                 `json:"roles_count"`
}

type CreateRequest struct {
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Patronymic string    `json:"patronymic"`
	Username   string    `json:"username"`
	Birthday   time.Time `json:"birthday"`
	GenderID   string    `json:"gender_id"`
	Password   string    `json:"password"`
}

// Все поля опциональны — частичное обновление (sqlc.narg + COALESCE).
type UpdateRequest struct {
	Name       *string    `json:"name"`
	Surname    *string    `json:"surname"`
	Patronymic *string    `json:"patronymic"`
	Username   *string    `json:"username"`
	Birthday   *time.Time `json:"birthday"`
	Status     *string    `json:"status"`
	GenderID   *string    `json:"gender_id"`
	Password   *string    `json:"password"`
}

// UserIDRequest — ID пользователя в path-параметре :user_id.
type UserIDRequest struct {
	UserID string `uri:"user_id"`
}

// UpdateUserRequest — как UpdateRequest, но с ID пользователя в path-параметре.
type UpdateUserRequest struct {
	UserID     string     `uri:"user_id"`
	Name       *string    `json:"name"`
	Surname    *string    `json:"surname"`
	Patronymic *string    `json:"patronymic"`
	Username   *string    `json:"username"`
	Birthday   *time.Time `json:"birthday"`
	Status     *string    `json:"status"`
	GenderID   *string    `json:"gender_id"`
	Password   *string    `json:"password"`
}

type RoleRequest struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}

type ListByServiceRequest struct {
	ServiceID string `uri:"service_id"`
	Page      int    `query:"page" validate:"omitempty,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

type MePermissionsRequest struct {
	ServiceID string `uri:"service_id"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}
