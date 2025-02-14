package cli

type GroupByValidator struct {
	GroupBy string
}

func (s *GroupByValidator) Validate() error {
	return nil
}
