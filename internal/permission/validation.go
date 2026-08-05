package permission

func (r UpsertRequest) Validate() error {
	if r.Code == "" {
		return ErrCodeRequired
	}
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Description == "" {
		return ErrDescriptionRequired
	}
	return nil
}
