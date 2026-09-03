package settings

// Validate: hasExisting — уже есть сохранённые настройки (тогда Password
// можно не передавать, старый пароль сохранится) или это первоначальная
// настройка (тогда Password обязателен — иначе нечего шифровать и сохранять).
func (r UpsertSMTPSettingsRequest) Validate(hasExisting bool) error {
	if r.Host == "" {
		return ErrHostRequired
	}
	if r.Port <= 0 {
		return ErrPortRequired
	}
	if r.Username == "" {
		return ErrUsernameRequired
	}
	if r.FromAddress == "" {
		return ErrFromAddressRequired
	}
	if !hasExisting && (r.Password == nil || *r.Password == "") {
		return ErrPasswordRequired
	}
	return nil
}
