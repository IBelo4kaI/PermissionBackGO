package user

import "unicode/utf8"

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
	return validateNickname(&r.Nickname)
}

func validateNickname(n *string) error {
	if n != nil && utf8.RuneCountInString(*n) > 100 {
		return ErrNicknameTooLong
	}
	return nil
}
