package invite

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"
	"time"

	repo "permisson/internal/database/sqlc"
	corporatedb "permisson/internal/database/sqlc_corporate"
	"permisson/internal/pkg/nullable"
	"permisson/internal/pkg/token"
	"permisson/internal/user"

	"github.com/google/uuid"
)

const defaultExpiresInHours = 72

// nowUTC — единая точка времени для сравнения expires_at. БД сравнивает
// через UTC_TIMESTAMP() (см. query/invites.sql), Go-сторона — этой функцией,
// чтобы не разъезжаться при разных часовых поясах сервера.
func nowUTC() time.Time {
	return time.Now().UTC()
}

// Service — держит оба Querier (сервисная БД и corporate-БД, см.
// internal/database/sqlc и internal/database/sqlc_corporate) плюс *user.Service,
// который переиспользуется для создания auth-аккаунта в Accept — так не
// дублируется валидация username/gender_id и хэширование пароля.
type Service struct {
	queries          repo.Querier
	corporateQueries corporatedb.Querier
	userService      *user.Service
}

func NewService(queries repo.Querier, corporateQueries corporatedb.Querier, userService *user.Service) *Service {
	return &Service{queries: queries, corporateQueries: corporateQueries, userService: userService}
}

func (s *Service) Create(ctx context.Context, req CreateRequest, createdBy string) (InviteResponse, error) {
	expiresInHours := defaultExpiresInHours
	if req.ExpiresInHours != nil {
		expiresInHours = *req.ExpiresInHours
	}

	// Повторная выдача на тот же email — тихо отзываем прежний активный инвайт.
	if req.Email != nil && *req.Email != "" {
		existing, err := s.queries.GetActiveInviteByEmail(ctx, nullable.String(*req.Email))
		if err == nil {
			if err := s.queries.RevokeInvite(ctx, existing.ID); err != nil {
				return InviteResponse{}, err
			}
		} else if !errors.Is(err, sql.ErrNoRows) {
			return InviteResponse{}, err
		}
	}

	code, err := token.Generate()
	if err != nil {
		return InviteResponse{}, err
	}

	id := uuid.NewString()
	err = s.queries.CreateInvite(ctx, repo.CreateInviteParams{
		ID:           id,
		Code:         code,
		Email:        nullable.StringPtr(req.Email),
		PersonID:     nullable.StringPtr(req.PersonID),
		CompanyID:    nullable.StringPtr(req.CompanyID),
		DepartmentID: nullable.StringPtr(req.DepartmentID),
		PositionID:   nullable.StringPtr(req.PositionID),
		CreatedBy:    createdBy,
		ExpiresAt:    nowUTC().Add(time.Duration(expiresInHours) * time.Hour),
	})
	if err != nil {
		return InviteResponse{}, err
	}

	return s.GetByID(ctx, id)
}

func (s *Service) GetByID(ctx context.Context, id string) (InviteResponse, error) {
	inv, err := s.queries.GetInviteByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return InviteResponse{}, ErrNotFound
		}
		return InviteResponse{}, err
	}
	return fromInvite(inv), nil
}

// List — без пагинации (см. user.Service.ListAll): фильтрация и сортировка
// применяются в Go после выборки всех строк.
func (s *Service) List(ctx context.Context, search, companyID, departmentID, positionID *string, sortBy, sortDir string) ([]InviteResponse, error) {
	rows, err := s.queries.ListInvites(ctx)
	if err != nil {
		return nil, err
	}

	rows = filterInvites(rows, search, companyID, departmentID, positionID)
	sortInvites(rows, sortBy, sortDir)

	items := make([]InviteResponse, 0, len(rows))
	for _, r := range rows {
		items = append(items, fromInvite(r))
	}
	return items, nil
}

func (s *Service) Revoke(ctx context.Context, id string) (InviteResponse, error) {
	inv, err := s.queries.GetInviteByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return InviteResponse{}, ErrNotFound
		}
		return InviteResponse{}, err
	}
	if inv.Used {
		return InviteResponse{}, ErrInviteUsed
	}
	if inv.Revoked {
		return InviteResponse{}, ErrInviteRevoked
	}

	if err := s.queries.RevokeInvite(ctx, id); err != nil {
		return InviteResponse{}, err
	}
	return s.GetByID(ctx, id)
}

// ValidateCode — публичная, без побочных эффектов проверка ссылки перед
// показом формы регистрации. Невалидный код/использованный/отозванный/
// истёкший инвайт — не ошибка, а просто {valid: false}, без утечки причины.
func (s *Service) ValidateCode(ctx context.Context, code string) (ValidateResponse, error) {
	inv, err := s.queries.GetInviteByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ValidateResponse{Valid: false}, nil
		}
		return ValidateResponse{}, err
	}

	if inv.Used || inv.Revoked || !inv.ExpiresAt.After(nowUTC()) {
		return ValidateResponse{Valid: false}, nil
	}

	resp := ValidateResponse{
		Valid:     true,
		Email:     nullable.StringOrNil(inv.Email),
		ExpiresAt: &inv.ExpiresAt,
	}

	// Если инвайт выпущен с уже существующим person_id — подтягиваем снапшот
	// профиля, чтобы фронт мог предзаполнить форму регистрации (см.
	// ValidateResponse.Person).
	if inv.PersonID.Valid && inv.PersonID.String != "" {
		person, err := s.corporateQueries.GetPersonByID(ctx, inv.PersonID.String)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return ValidateResponse{}, err
		}
		if err == nil {
			resp.HasPerson = true
			resp.Person = &PersonSnapshot{
				FirstName:  person.FirstName,
				LastName:   person.LastName,
				Patronymic: nullable.StringOrNil(person.Patronymic),
				Birthday:   nullable.TimeOrNil(person.Birthday),
				Phone:      nullable.StringOrNil(person.Phone),
				Email:      nullable.StringOrNil(person.Email),
			}
		}
	}

	return resp, nil
}

// Accept — принимает инвайт: создаёт auth-аккаунт (users, через
// user.Service.Create — как обычная регистрация) и, отдельно, профиль в
// corporate-БД: создаёт нового person либо линкует уже существующего (если
// инвайт был выпущен с person_id), плюс employment, если в инвайте заданы
// и company_id, и position_id (в employments они NOT NULL).
//
// Сервисная и corporate-БД — разные подключения, атомарности между ними нет:
// если запись в corporate-БД не удалась после успешного создания users, аккаунт
// уже создан, а профиля ещё нет. Это известное ограничение текущей
// архитектуры (нет двухфазного commit между базами), не решается на этом
// этапе.
func (s *Service) Accept(ctx context.Context, req AcceptRequest) (user.UserResponse, error) {
	inv, err := s.queries.GetInviteByCode(ctx, req.Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.UserResponse{}, ErrCodeNotFound
		}
		return user.UserResponse{}, err
	}
	if inv.Used {
		return user.UserResponse{}, ErrInviteUsed
	}
	if inv.Revoked {
		return user.UserResponse{}, ErrInviteRevoked
	}
	if !inv.ExpiresAt.After(nowUTC()) {
		return user.UserResponse{}, ErrInviteExpired
	}

	createReq := user.CreateRequest{
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
		Username:   req.Username,
		Birthday:   req.Birthday,
		GenderID:   req.GenderID,
		Password:   req.Password,
	}
	if err := createReq.Validate(); err != nil {
		return user.UserResponse{}, err
	}

	newUser, err := s.userService.Create(ctx, createReq)
	if err != nil {
		return user.UserResponse{}, err
	}

	if err := s.linkOrCreatePerson(ctx, inv, req, newUser.ID); err != nil {
		return user.UserResponse{}, err
	}

	if err := s.queries.MarkInviteUsed(ctx, repo.MarkInviteUsedParams{
		UserID: nullable.String(newUser.ID),
		ID:     inv.ID,
	}); err != nil {
		return user.UserResponse{}, err
	}

	return newUser, nil
}

// linkOrCreatePerson — часть Accept, вынесенная отдельно ради читаемости.
func (s *Service) linkOrCreatePerson(ctx context.Context, inv repo.Invite, req AcceptRequest, authUserID string) error {
	personID := inv.PersonID.String
	if inv.PersonID.Valid && personID != "" {
		// Существующий person: перезаписываем профиль тем, что пользователь
		// подтвердил/поправил в форме (она была предзаполнена из
		// ValidateResponse.Person), плюс линкуем auth-аккаунт.
		if err := s.corporateQueries.UpdatePersonOnLink(ctx, corporatedb.UpdatePersonOnLinkParams{
			FirstName:  req.Name,
			LastName:   req.Surname,
			Patronymic: nullable.StringPtr(&req.Patronymic),
			Birthday:   nullable.Time(req.Birthday),
			Phone:      nullable.String(req.Phone),
			Email:      nullable.String(req.Email),
			AuthUserID: nullable.String(authUserID),
			ID:         personID,
		}); err != nil {
			return err
		}
	} else {
		personID = uuid.NewString()
		email := inv.Email
		if req.Email != "" {
			email = nullable.String(req.Email)
		}
		if err := s.corporateQueries.CreatePerson(ctx, corporatedb.CreatePersonParams{
			ID:         personID,
			FirstName:  req.Name,
			LastName:   req.Surname,
			Patronymic: nullable.StringPtr(&req.Patronymic),
			Birthday:   nullable.Time(req.Birthday),
			Phone:      nullable.String(req.Phone),
			Email:      email,
			AuthUserID: nullable.String(authUserID),
		}); err != nil {
			return err
		}
	}

	// employments.company_id и position_id — NOT NULL: создаём трудоустройство
	// только если оба заданы в инвайте.
	if inv.CompanyID.Valid && inv.CompanyID.String != "" && inv.PositionID.Valid && inv.PositionID.String != "" {
		if err := s.corporateQueries.CreateEmployment(ctx, corporatedb.CreateEmploymentParams{
			ID:           uuid.NewString(),
			PersonID:     personID,
			CompanyID:    inv.CompanyID.String,
			DepartmentID: inv.DepartmentID,
			PositionID:   inv.PositionID.String,
			Status:       "active",
			HiredAt:      nowUTC(),
		}); err != nil {
			return err
		}
	}

	return nil
}

// filterInvites — регистронезависимый поиск по email + точные фильтры по
// company_id/department_id/position_id (см. user.filterUserRows — тот же приём
// для списка без SQL-пагинации).
func filterInvites(rows []repo.Invite, search, companyID, departmentID, positionID *string) []repo.Invite {
	filtered := make([]repo.Invite, 0, len(rows))
	for _, r := range rows {
		if search != nil && !strings.Contains(strings.ToLower(r.Email.String), strings.ToLower(*search)) {
			continue
		}
		if companyID != nil && *companyID != "" && r.CompanyID.String != *companyID {
			continue
		}
		if departmentID != nil && *departmentID != "" && r.DepartmentID.String != *departmentID {
			continue
		}
		if positionID != nil && *positionID != "" && r.PositionID.String != *positionID {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}

// sortInvites сортирует rows "на месте" по одной из SortableColumns
// (см. user.sortUserRows).
func sortInvites(rows []repo.Invite, sortBy, sortDir string) {
	less := func(i, j int) bool {
		switch sortBy {
		case "email":
			return rows[i].Email.String < rows[j].Email.String
		case "expires_at":
			return rows[i].ExpiresAt.Before(rows[j].ExpiresAt)
		default: // "created_at" и любое неизвестное значение
			return rows[i].CreatedAt.Before(rows[j].CreatedAt)
		}
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if sortDir == "asc" {
			return less(i, j)
		}
		return less(j, i)
	})
}
