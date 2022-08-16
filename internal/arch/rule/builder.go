package rule

type Builder interface {
	That(t That) Builder
	AndThat(t That) Builder
	Except(s ...Except) Builder
	Must(e Expect) Builder
	AndMust(e Expect) Builder
	Should(e Expect) Builder
	AndShould(e Expect) Builder
	Could(e Expect) Builder
	AndCould(e Expect) Builder
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

type Expect interface {
	Evaluate(rule Builder) []CoreViolation
	GetErrors() []error
}

type Because string
