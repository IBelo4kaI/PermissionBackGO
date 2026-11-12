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
	if err := validateNickname(&r.Nickname); err != nil {
		return err
	}
	return validatePhone(r.Phone)
}

func validateNickname(n *string) error {
	if n != nil && utf8.RuneCountInString(*n) > 100 {
		return ErrNicknameTooLong
	}
	return nil
}

func validatePhone(p *string) error {
	if p != nil && utf8.RuneCountInString(*p) > 18 {
		return ErrPhoneTooLong
	}
	return nil
}
