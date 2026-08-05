package role

// Validate — Python-модель RoleCreate требует наличия name/description
// (pydantic допустил бы и пустую строку, но здесь чуть строже — как и в
// остальных сущностях).
func (r UpsertRequest) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Description == "" {
		return ErrDescriptionRequired
	}
	return nil
}
