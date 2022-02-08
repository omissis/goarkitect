package arch

type That interface {
	Evaluate(rule *RuleSubjectBuilder)
}

type Should interface {
	Evaluate(rule *RuleSubjectBuilder) []Violation
}

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
