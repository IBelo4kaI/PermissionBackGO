package user

func (r CreateRequest) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Surname == "" {
		return ErrSurnameRequired
	}
	if r.Username == "" {
		return ErrUsernameRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	if r.Birthday.IsZero() {
		return ErrBirthdayRequired
	}
	if r.GenderID == "" {
		return ErrGenderRequired
	}
	return nil
}
