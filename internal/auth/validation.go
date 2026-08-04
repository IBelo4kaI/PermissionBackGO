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
