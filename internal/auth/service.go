package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/mailer"
	"permisson/internal/pkg/nullable"
	"permisson/internal/pkg/passwordhash"
	"permisson/internal/pkg/token"
	"permisson/internal/settings"
	"time"

	"github.com/google/uuid"
)

// passwordResetTTL — срок жизни ссылки восстановления пароля. Не вынесено в
// конфиг: в отличие от SESSION_TTL это не эксплуатационный параметр, а
// константа политики безопасности.
const passwordResetTTL = time.Hour

type Service struct {
	queries    repo.Querier
	sessionTTL time.Duration
	settings   *settings.Service
	appBaseURL string
}

func NewService(queries repo.Querier, sessionTTL time.Duration, settingsService *settings.Service, appBaseURL string) *Service {
	return &Service{
		queries:    queries,
		sessionTTL: sessionTTL,
		settings:   settingsService,
		appBaseURL: appBaseURL,
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResult, error) {
	user, err := s.queries.GetUserByUsername(ctx, req.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, err
	}

	if !passwordhash.Verify(user.Password, req.Password) {
		return LoginResult{}, ErrInvalidCredentials
	}

	if user.Status != repo.UsersStatusActive {
		return LoginResult{}, ErrUserInactive
	}

	rawToken, err := token.Generate()
	if err != nil {
		return LoginResult{}, err
	}

	expiresAt := time.Now().UTC().Add(s.sessionTTL)

	err = s.queries.CreateSession(ctx, repo.CreateSessionParams{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: token.Hash(rawToken),
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: rawToken, ExpiresAt: expiresAt}, nil
}

func (s *Service) ValidateSession(ctx context.Context, rawToken string) (bool, error) {
	if rawToken == "" {
		return false, nil
	}

	_, err := s.queries.GetSessionByTokenHash(ctx, token.Hash(rawToken))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Logout удаляет сессию по сырому токену (в БД хранится только его SHA-256 хэш).
// Идемпотентен: если сессия уже не существует (например, истекла), ошибки нет.
func (s *Service) Logout(ctx context.Context, rawToken string) error {
	if rawToken == "" {
		return ErrSessionTokenMissing
	}

	return s.queries.DeleteSessionByTokenHash(ctx, token.Hash(rawToken))
}

func (s *Service) SessionUserID(ctx context.Context, rawToken string) (string, error) {
	if rawToken == "" {
		return "", ErrInvalidSession
	}

	session, err := s.queries.GetSessionByTokenHash(ctx, token.Hash(rawToken))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidSession
		}
		return "", err
	}

	return session.UserID, nil
}

// ForgotPassword намеренно не возвращает вызывающему информацию о том,
// найден ли пользователь и отправлено ли письмо: если различать ответы
// (200 vs 404, или 200 vs 500 на сбое отправки), ручка становится оракулом
// существования логинов. Пользователь не найден, SMTP не настроен, само
// письмо не отправилось — во всех случаях либо тихо выходим, либо пишем в
// лог и всё равно возвращаем nil.
func (s *Service) ForgotPassword(ctx context.Context, username string) error {
	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	rawToken, err := token.Generate()
	if err != nil {
		return err
	}

	if err := s.queries.CreatePasswordReset(ctx, repo.CreatePasswordResetParams{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: token.Hash(rawToken),
		ExpiresAt: time.Now().UTC().Add(passwordResetTTL),
	}); err != nil {
		return err
	}

	mailCfg, err := s.settings.GetForSending(ctx)
	if err != nil {
		log.Printf("forgot-password: SMTP не настроен или недоступен, письмо не отправлено: %v", err)
		return nil
	}

	resetLink := s.appBaseURL + "/reset-password/" + rawToken
	body := fmt.Sprintf(
		`<p>Для сброса пароля перейдите по ссылке: <a href="%s">%s</a></p><p>Ссылка действительна 1 час. Если вы не запрашивали восстановление пароля — просто проигнорируйте это письмо.</p>`,
		resetLink, resetLink,
	)

	if err := mailer.Send(mailCfg, username, "Восстановление пароля", body); err != nil {
		log.Printf("forgot-password: не удалось отправить письмо: %v", err)
	}

	return nil
}

// ResetPassword — принимает токен из письма и новый пароль. После успешной
// смены удаляет все сессии пользователя (принудительный релогин везде).
func (s *Service) ResetPassword(ctx context.Context, rawToken, newPassword string) error {
	reset, err := s.queries.GetPasswordResetByTokenHash(ctx, token.Hash(rawToken))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidResetToken
		}
		return err
	}

	hashed, err := passwordhash.Hash(newPassword)
	if err != nil {
		return err
	}

	if err := s.queries.UpdateUser(ctx, repo.UpdateUserParams{
		ID:       reset.UserID,
		Password: nullable.String(hashed),
	}); err != nil {
		return err
	}

	if err := s.queries.MarkPasswordResetUsed(ctx, reset.ID); err != nil {
		return err
	}

	return s.queries.DeleteSessionsByUserID(ctx, reset.UserID)
}
