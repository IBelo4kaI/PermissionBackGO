package settings

import (
	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/nullable"
)

func fromRow(r repo.SmtpSetting) SMTPSettingsResponse {
	return SMTPSettingsResponse{
		Configured:  true,
		Host:        r.Host,
		Port:        r.Port,
		Username:    r.Username,
		FromAddress: r.FromAddress,
		FromName:    nullable.StringOrNil(r.FromName),
		UseTLS:      r.UseTls,
	}
}
