package validator

type Validatable interface {
	Validate() (bool, error)
}
