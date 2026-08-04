package auth

import "time"

// LoginRequest — аналог pydantic-модели Login (login, password).
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginResult — то, что раньше возвращалось из AuthService.login() как tuple (token, expires_at).
type LoginResult struct {
	Token     string
	ExpiresAt time.Time
}

// ValidateSessionResponse — аналог {"valid": bool} из /auth/validate-session.
type ValidateSessionResponse struct {
	Valid bool `json:"valid"`
}
