package rule

type Builder interface {
	That(t That) Builder
	AndThat(t That) Builder
	Except(s ...Except) Builder
	Should(e Should) Builder
	AndShould(e Should) Builder
	Because(b Because) ([]Violation, []error)
	AddError(err error)
	GetErrors() []error
}

type That interface {
	Evaluate(rule Builder)
	GetErrors() []error
}

type Except interface {
	Evaluate(rule Builder)
	GetErrors() []error
}

type Should interface {
	Evaluate(rule Builder) []Violation
	GetErrors() []error
}

type Because string

func NewViolation(message string) Violation {
	return Violation{
		message: message,
	}
}

type Violation struct {
	message string
}

func (v Violation) String() string {
	return v.message
}
