package auth

func (r LoginRequest) Validate() error {
	if r.Login == "" {
		return ErrUsernameRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
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
