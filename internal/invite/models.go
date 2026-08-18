package invite

import "time"

// CreateRequest — тело POST /invites. Все поля, кроме срока действия,
// необязательны: email — просто пометка "для кого" (доставка — копированием
// ссылки на фронте, не email-рассылкой), person_id/company_id/department_id/
// position_id — непрозрачные id из внешней corporate-БД (без FK и валидации).
type CreateRequest struct {
	Email          *string `json:"email"`
	PersonID       *string `json:"person_id"`
	CompanyID      *string `json:"company_id"`
	DepartmentID   *string `json:"department_id"`
	PositionID     *string `json:"position_id"`
	ExpiresInHours *int    `json:"expires_in_hours" validate:"omitempty,min=1" description:"По умолчанию 72 часа"`
}

// InviteResponse — карточка инвайта. Code отдаётся всегда, в отличие от
// сессий/API-ключей код инвайта не хэшируется (см. schema/003_add_invites.sql).
type InviteResponse struct {
	ID           string    `json:"id"`
	Code         string    `json:"code"`
	Email        *string   `json:"email"`
	PersonID     *string   `json:"person_id"`
	CompanyID    *string   `json:"company_id"`
	DepartmentID *string   `json:"department_id"`
	PositionID   *string   `json:"position_id"`
	CreatedBy    string    `json:"created_by"`
	UserID       *string   `json:"user_id"`
	Used         bool      `json:"used"`
	Revoked      bool      `json:"revoked"`
	Status       string    `json:"status" example:"pending" description:"pending, used, revoked или expired — вычисляется, в БД не хранится"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// InviteIDRequest — ID инвайта в path-параметре :invite_id.
type InviteIDRequest struct {
	InviteID string `uri:"invite_id"`
}

// ListRequest — query-параметры GET /invites. Без пагинации (см.
// user.ListAllRequest) — список отдаётся целиком.
type ListRequest struct {
	Search       string `query:"search" validate:"omitempty,max=255" description:"Поиск по email"`
	CompanyID    string `query:"company_id"`
	DepartmentID string `query:"department_id"`
	PositionID   string `query:"position_id"`
	SortBy       string `query:"sort_by" description:"email, expires_at или created_at (по умолчанию created_at)"`
	SortDir      string `query:"sort_dir" validate:"omitempty,oneof=asc desc" description:"asc или desc (по умолчанию desc)"`
}

// SortableColumns — белый список для sort_by (см. query.QuerySort).
var SortableColumns = []string{"email", "expires_at", "created_at"}

const DefaultSortColumn = "created_at"

// CodeRequest — код инвайта в path-параметре :code (GET /invites/code/:code).
type CodeRequest struct {
	Code string `uri:"code"`
}

// ValidateResponse — публичный, без побочных эффектов ответ на проверку
// ссылки перед показом формы регистрации. Email/ExpiresAt — только если valid.
//
// HasPerson/Person — если инвайт был выпущен с уже существующим person_id
// (см. CreateRequest.PersonID), сюда подтягивается снапшот профиля из
// corporate-БД, чтобы фронт мог предзаполнить форму и дать пользователю
// проверить/поправить данные вместо ввода с нуля (см. invite.Service.Accept).
type ValidateResponse struct {
	Valid     bool            `json:"valid"`
	Email     *string         `json:"email,omitempty"`
	ExpiresAt *time.Time      `json:"expires_at,omitempty"`
	HasPerson bool            `json:"has_person"`
	Person    *PersonSnapshot `json:"person,omitempty"`
}

// PersonSnapshot — часть persons, которую можно показать в форме
// регистрации на проверку/редактирование (см. ValidateResponse.Person).
type PersonSnapshot struct {
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Patronymic *string    `json:"patronymic"`
	Birthday   *time.Time `json:"birthday"`
	Phone      *string    `json:"phone"`
	Email      *string    `json:"email"`
}

// AcceptRequest — тело POST /invites/code/:code/accept. Name/Surname/
// Patronymic/Username/Birthday/GenderID/Password повторяют user.CreateRequest
// один в один (идут в users); Phone/Email — редактируемые поля профиля
// persons (см. Service.linkOrCreatePerson): если у инвайта уже был person_id
// с ValidateResponse.Person — это то, что пользователь подтвердил/поправил в
// форме, если нового человека создаём — просто исходные значения.
type AcceptRequest struct {
	Code       string    `uri:"code"`
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Patronymic string    `json:"patronymic"`
	Username   string    `json:"username"`
	Birthday   time.Time `json:"birthday"`
	GenderID   string    `json:"gender_id"`
	Password   string    `json:"password"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
}
