package auth

import "time"

type LoginRequest struct {
	Login    string `json:"login" validate:"required,max=100" example:"admin" description:"User email"`
	Password string `json:"password" validate:"required,max=255" example:"P@ssw0rd123" description:"User password"`
}

type LoginResult struct {
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time `json:"expires_at" format:"date-time" description:"Token expiration time"`
}

type ValidateSessionResponse struct {
	Valid bool `json:"valid" example:"true" description:"Whether the provided session is valid"`
}

type LogoutResponse struct {
	Message string `json:"message" example:"Вы успешно вышли из системы" description:"Result message"`
}

// ForgotPasswordRequest.Username — логин (email). Ответ всегда одинаковый,
// есть такой пользователь или нет (см. Service.ForgotPassword) — иначе
// ручка становится оракулом существования логинов.
type ForgotPasswordRequest struct {
	Username string `json:"username" validate:"required,max=100" example:"user@example.com"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message" example:"Если такой пользователь существует, письмо с инструкциями отправлено"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,max=255"`
}

type ResetPasswordResponse struct {
	Message string `json:"message" example:"Пароль изменён"`
}
