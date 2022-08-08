package file

import (
	"goarkitect/internal/arch/rule"
)

func All() *RuleBuilder {
	return NewRuleBuilder().AllFiles()
}

func One(filename string) *RuleBuilder {
	return NewRuleBuilder().File(filename)
}

func NewRuleBuilder() *RuleBuilder {
	return &RuleBuilder{}
}

type RuleBuilder struct {
	thats      []rule.That
	excepts    []rule.Except
	shoulds    []rule.Should
	because    rule.Because
	Violations []rule.Violation

	files []string
}

func (rb *RuleBuilder) That(t rule.That) rule.Builder {
	rb.thats = []rule.That{t}

	return rb
}

func (rb *RuleBuilder) AndThat(t rule.That) rule.Builder {
	rb.thats = append(rb.thats, t)

	return rb
}

func (rb *RuleBuilder) Except(s ...rule.Except) rule.Builder {
	rb.excepts = s

	return rb
}

func (rb *RuleBuilder) Should(e rule.Should) rule.Builder {
	rb.shoulds = []rule.Should{e}

	return rb
}

func (rb *RuleBuilder) AndShould(e rule.Should) rule.Builder {
	rb.shoulds = append(rb.shoulds, e)

	return rb
}

func (rb *RuleBuilder) Because(b rule.Because) []rule.Violation {
	rb.because = b

	for _, that := range rb.thats {
		that.Evaluate(rb)
	}

	for _, except := range rb.excepts {
		except.Evaluate(rb)
	}

	for _, should := range rb.shoulds {
		vs := should.Evaluate(rb)
		if len(vs) > 0 {
			rb.Violations = append(rb.Violations, vs...)
		}
	}

	return rb.Violations
}

func (rb *RuleBuilder) GetFiles() []string {
	return rb.files
}

func (rb *RuleBuilder) SetFiles(files []string) {
	rb.files = files
}

func (rb *RuleBuilder) AllFiles() *RuleBuilder {
	return rb
}

func (rb *RuleBuilder) File(filename string) *RuleBuilder {
	rb.files = []string{filename}

	return rb
}
