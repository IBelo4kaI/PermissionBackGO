package auth

const (
	minSessionTTLSeconds = 60
	maxSessionTTLSeconds = 60 * 24 * 60 * 60
)

func (r LoginRequest) Validate() error {
	if r.Login == "" {
		return ErrUsernameRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	if r.SessionTTLSeconds != nil && (*r.SessionTTLSeconds < minSessionTTLSeconds || *r.SessionTTLSeconds > maxSessionTTLSeconds) {
		return ErrInvalidSessionTTL
	}
	return nil
}

func (r ForgotPasswordRequest) Validate() error {
	if r.Username == "" {
		return ErrUsernameRequired
	}
	return nil
}

func (r ResetPasswordRequest) Validate() error {
	if r.Token == "" {
		return ErrTokenRequired
	}
	if r.NewPassword == "" {
		return ErrNewPasswordRequired
	}
	return nil
}
