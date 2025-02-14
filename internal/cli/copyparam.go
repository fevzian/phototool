package cli

type CopyValidator struct {
	IsCopy bool
}

func (s *CopyValidator) Validate() error {
	return nil
}
