package invite

import (
	"database/sql"
	"time"

	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/nullable"
)

// inviteRow — общий промежуточный тип: у GetInviteByIDRow/GetInviteByCodeRow/
// GetActiveInviteByEmailRow/ListInvitesRow одинаковый набор колонок, но
// разные сгенерированные структуры (см. тот же приём в user.userRow).
//
// Это не совпадает с repo.Invite — после schema/004_invites_created_by_nullable.sql
// (ALTER ... MODIFY COLUMN без реального изменения позиции) sqlc-движок при
// кодгене переставил created_by в модели в конец структуры, а во всех
// запросах порядок колонок остался как в SELECT — из-за расхождения sqlc
// перестал переиспользовать repo.Invite для этих запросов и сгенерировал
// каждой свой Row-тип.
type inviteRow struct {
	ID           string
	Code         string
	Email        sql.NullString
	PersonID     sql.NullString
	CompanyID    sql.NullString
	DepartmentID sql.NullString
	PositionID   sql.NullString
	CreatedBy    sql.NullString
	UserID       sql.NullString
	Used         bool
	Revoked      bool
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func fromGetInviteByIDRow(r repo.GetInviteByIDRow) inviteRow {
	return inviteRow{
		ID: r.ID, Code: r.Code, Email: r.Email, PersonID: r.PersonID,
		CompanyID: r.CompanyID, DepartmentID: r.DepartmentID, PositionID: r.PositionID,
		CreatedBy: r.CreatedBy, UserID: r.UserID, Used: r.Used, Revoked: r.Revoked,
		ExpiresAt: r.ExpiresAt, CreatedAt: r.CreatedAt,
	}
}

func fromGetInviteByCodeRow(r repo.GetInviteByCodeRow) inviteRow {
	return inviteRow{
		ID: r.ID, Code: r.Code, Email: r.Email, PersonID: r.PersonID,
		CompanyID: r.CompanyID, DepartmentID: r.DepartmentID, PositionID: r.PositionID,
		CreatedBy: r.CreatedBy, UserID: r.UserID, Used: r.Used, Revoked: r.Revoked,
		ExpiresAt: r.ExpiresAt, CreatedAt: r.CreatedAt,
	}
}

func fromGetActiveInviteByEmailRow(r repo.GetActiveInviteByEmailRow) inviteRow {
	return inviteRow{
		ID: r.ID, Code: r.Code, Email: r.Email, PersonID: r.PersonID,
		CompanyID: r.CompanyID, DepartmentID: r.DepartmentID, PositionID: r.PositionID,
		CreatedBy: r.CreatedBy, UserID: r.UserID, Used: r.Used, Revoked: r.Revoked,
		ExpiresAt: r.ExpiresAt, CreatedAt: r.CreatedAt,
	}
}

func fromListInvitesRow(r repo.ListInvitesRow) inviteRow {
	return inviteRow{
		ID: r.ID, Code: r.Code, Email: r.Email, PersonID: r.PersonID,
		CompanyID: r.CompanyID, DepartmentID: r.DepartmentID, PositionID: r.PositionID,
		CreatedBy: r.CreatedBy, UserID: r.UserID, Used: r.Used, Revoked: r.Revoked,
		ExpiresAt: r.ExpiresAt, CreatedAt: r.CreatedAt,
	}
}

// statusOf — вычисляемый статус инвайта: used/revoked/expired/pending.
// В БД отдельно не хранится (см. schema/003_add_invites.sql).
func statusOf(i inviteRow) string {
	switch {
	case i.Used:
		return "used"
	case i.Revoked:
		return "revoked"
	case !i.ExpiresAt.After(nowUTC()):
		return "expired"
	default:
		return "pending"
	}
}

func fromInvite(i inviteRow) InviteResponse {
	return InviteResponse{
		ID:           i.ID,
		Code:         i.Code,
		Email:        nullable.StringOrNil(i.Email),
		PersonID:     nullable.StringOrNil(i.PersonID),
		CompanyID:    nullable.StringOrNil(i.CompanyID),
		DepartmentID: nullable.StringOrNil(i.DepartmentID),
		PositionID:   nullable.StringOrNil(i.PositionID),
		CreatedBy:    nullable.StringOrNil(i.CreatedBy),
		UserID:       nullable.StringOrNil(i.UserID),
		Used:         i.Used,
		Revoked:      i.Revoked,
		Status:       statusOf(i),
		ExpiresAt:    i.ExpiresAt,
		CreatedAt:    i.CreatedAt,
	}
}
