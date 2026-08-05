package user

import (
	"database/sql"
	"time"

	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/nullable"
	"permisson/internal/role"
)

// userRow — общий промежуточный тип: у всех четырёх Row-типов sqlc
// (GetUserByIDRow, ListUsersRow, ListAllUsersRow, ListUsersByServiceIDRow)
// одинаковый набор колонок, но разные сгенерированные имена структур.
type userRow struct {
	ID         string
	Name       string
	Surname    string
	Patronymic sql.NullString
	Username   string
	GenderID   string
	Birthday   time.Time
	Password   string
	Status     repo.UsersStatus
	CreatedAt  time.Time
}

func fromGetUserRow(r repo.GetUserByIDRow) userRow {
	return userRow{
		ID: r.ID, Name: r.Name, Surname: r.Surname, Patronymic: r.Patronymic,
		Username: r.Username, GenderID: r.GenderID, Birthday: r.Birthday,
		Password: r.Password, Status: r.Status, CreatedAt: r.CreatedAt,
	}
}

func fromListUsersRow(r repo.ListUsersRow) userRow {
	return userRow{
		ID: r.ID, Name: r.Name, Surname: r.Surname, Patronymic: r.Patronymic,
		Username: r.Username, GenderID: r.GenderID, Birthday: r.Birthday,
		Password: r.Password, Status: r.Status, CreatedAt: r.CreatedAt,
	}
}

func fromListAllUsersRow(r repo.ListAllUsersRow) userRow {
	return userRow{
		ID: r.ID, Name: r.Name, Surname: r.Surname, Patronymic: r.Patronymic,
		Username: r.Username, GenderID: r.GenderID, Birthday: r.Birthday,
		Password: r.Password, Status: r.Status, CreatedAt: r.CreatedAt,
	}
}

func fromListUsersByServiceIDRow(r repo.ListUsersByServiceIDRow) userRow {
	return userRow{
		ID: r.ID, Name: r.Name, Surname: r.Surname, Patronymic: r.Patronymic,
		Username: r.Username, GenderID: r.GenderID, Birthday: r.Birthday,
		Password: r.Password, Status: r.Status, CreatedAt: r.CreatedAt,
	}
}

// fromRole — роль пользователя в том же формате, что и в role-пакете:
// счётчики (user_count, permissions_count) остаются нулевыми, как в Python-версии.
func fromRole(r repo.Role) role.RoleResponse {
	return role.RoleResponse{
		ID:          r.ID,
		ServiceID:   nullable.StringOrNil(r.ServiceID),
		Name:        r.Name,
		Description: r.Description,
		IsGlobal:    r.IsGlobal,
		CreatedAt:   r.CreatedAt,
	}
}

func fromUserRow(u userRow, gender repo.Gender, roles []repo.Role) UserResponse {
	roleItems := make([]role.RoleResponse, 0, len(roles))
	for _, r := range roles {
		roleItems = append(roleItems, fromRole(r))
	}

	return UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Surname:    u.Surname,
		Patronymic: nullable.StringOrNil(u.Patronymic),
		Username:   u.Username,
		Birthday:   u.Birthday,
		Status:     u.Status,
		CreatedAt:  u.CreatedAt,
		Gender:     gender,
		Roles:      roleItems,
		RolesCount: len(roleItems),
	}
}
