package auth

import "time"

// LoginRequest describes user credentials.
type LoginRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=100" example:"admin" description:"User email"`
	Password string `json:"password" validate:"required,min=8,max=255" example:"P@ssw0rd123" description:"User password"`
}

// LoginResponse describes successful authentication.
type LoginResult struct {
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time `json:"expires_at" format:"date-time" description:"Token expiration time"`
}

// ValidateSessionResponse indicates whether the session is valid.
type ValidateSessionResponse struct {
	Valid bool `json:"valid" example:"true" description:"Whether the provided session is valid"`
}
