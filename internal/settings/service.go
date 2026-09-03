package settings

import (
	"context"
	"database/sql"
	"errors"

	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/mailer"
	"permisson/internal/pkg/nullable"
	"permisson/internal/pkg/secretcrypt"
)

// Service держит ключ шифрования (см. env.GetSettingsEncryptionKey) —
// SMTP-пароль в БД лежит только в зашифрованном виде (см.
// internal/pkg/secretcrypt), расшифровывается только в GetForSending.
type Service struct {
	queries       repo.Querier
	encryptionKey []byte
}

func NewService(queries repo.Querier, encryptionKey []byte) *Service {
	return &Service{queries: queries, encryptionKey: encryptionKey}
}

func (s *Service) Get(ctx context.Context) (SMTPSettingsResponse, error) {
	row, err := s.queries.GetSMTPSettings(ctx, SingletonID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SMTPSettingsResponse{Configured: false}, nil
		}
		return SMTPSettingsResponse{}, err
	}
	return fromRow(row), nil
}

func (s *Service) Upsert(ctx context.Context, req UpsertSMTPSettingsRequest) (SMTPSettingsResponse, error) {
	existing, err := s.queries.GetSMTPSettings(ctx, SingletonID)
	hasExisting := true
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			hasExisting = false
		} else {
			return SMTPSettingsResponse{}, err
		}
	}

	if err := req.Validate(hasExisting); err != nil {
		return SMTPSettingsResponse{}, err
	}

	passwordEncrypted := existing.PasswordEncrypted
	if req.Password != nil && *req.Password != "" {
		encrypted, err := secretcrypt.Encrypt(s.encryptionKey, *req.Password)
		if err != nil {
			return SMTPSettingsResponse{}, err
		}
		passwordEncrypted = encrypted
	}

	if err := s.queries.UpsertSMTPSettings(ctx, repo.UpsertSMTPSettingsParams{
		ID:                SingletonID,
		Host:              req.Host,
		Port:              req.Port,
		Username:          req.Username,
		PasswordEncrypted: passwordEncrypted,
		FromAddress:       req.FromAddress,
		FromName:          nullable.StringPtr(req.FromName),
		UseTls:            req.UseTLS,
	}); err != nil {
		return SMTPSettingsResponse{}, err
	}

	return s.Get(ctx)
}

// GetForSending — не наружу (нет ручки): используется только auth.Service
// при отправке письма восстановления пароля. Расшифровывает пароль.
func (s *Service) GetForSending(ctx context.Context) (mailer.Config, error) {
	row, err := s.queries.GetSMTPSettings(ctx, SingletonID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mailer.Config{}, ErrNotConfigured
		}
		return mailer.Config{}, err
	}

	password, err := secretcrypt.Decrypt(s.encryptionKey, row.PasswordEncrypted)
	if err != nil {
		return mailer.Config{}, err
	}

	return mailer.Config{
		Host:        row.Host,
		Port:        int(row.Port),
		Username:    row.Username,
		Password:    password,
		FromAddress: row.FromAddress,
		FromName:    nullable.StringOr(row.FromName, ""),
		UseTLS:      row.UseTls,
	}, nil
}
