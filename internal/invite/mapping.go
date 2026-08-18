package invite

import (
	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/nullable"
)

// statusOf — вычисляемый статус инвайта: used/revoked/expired/pending.
// В БД отдельно не хранится (см. schema/003_add_invites.sql).
func statusOf(i repo.Invite) string {
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

func fromInvite(i repo.Invite) InviteResponse {
	return InviteResponse{
		ID:           i.ID,
		Code:         i.Code,
		Email:        nullable.StringOrNil(i.Email),
		PersonID:     nullable.StringOrNil(i.PersonID),
		CompanyID:    nullable.StringOrNil(i.CompanyID),
		DepartmentID: nullable.StringOrNil(i.DepartmentID),
		PositionID:   nullable.StringOrNil(i.PositionID),
		CreatedBy:    i.CreatedBy,
		UserID:       nullable.StringOrNil(i.UserID),
		Used:         i.Used,
		Revoked:      i.Revoked,
		Status:       statusOf(i),
		ExpiresAt:    i.ExpiresAt,
		CreatedAt:    i.CreatedAt,
	}
}
