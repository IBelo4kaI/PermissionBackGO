package service

func (r UpsertRequest) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	if len(r.Name) > 100 {
		return ErrNameTooLong
	}
	if r.Description == "" {
		return ErrDescriptionRequired
	}
	if r.Prefix == "" {
		return ErrPrefixRequired
	}
	if len(r.Prefix) > 5 {
		return ErrPrefixTooLong
	}
	return nil
}
