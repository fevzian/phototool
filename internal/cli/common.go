package cli

type Validator interface {
	Validate() error
}

func Validate(validators []Validator) error {
	for _, v := range validators {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}
