package file

import (
	"goarkitect/internal/arch/rule"
)

func All() *RuleBuilder {
	return NewRuleBuilder()
}

func One(filename string) *RuleBuilder {
	rb := NewRuleBuilder()

	rb.files = []string{filename}

	return rb
}

func Set(filenames ...string) *RuleBuilder {
	rb := NewRuleBuilder()

	rb.SetFiles(filenames)

	return rb
}

func NewRuleBuilder() *RuleBuilder {
	return &RuleBuilder{}
}

type RuleBuilder struct {
	thats      []rule.That
	excepts    []rule.Except
	shoulds    []rule.Should
	because    rule.Because
	violations []rule.Violation
	errors     []error
	locked     bool

	files []string
}

func (rb *RuleBuilder) That(t rule.That) rule.Builder {
	rb.assertNotLocked()

	rb.thats = []rule.That{t}

	return rb
}

func (rb *RuleBuilder) AndThat(t rule.That) rule.Builder {
	rb.assertNotLocked()

	rb.thats = append(rb.thats, t)

	return rb
}

func (rb *RuleBuilder) Except(s ...rule.Except) rule.Builder {
	rb.assertNotLocked()

	rb.excepts = s

	return rb
}

func (rb *RuleBuilder) Should(e rule.Should) rule.Builder {
	rb.assertNotLocked()

	rb.shoulds = []rule.Should{e}

	return rb
}

func (rb *RuleBuilder) AndShould(e rule.Should) rule.Builder {
	rb.assertNotLocked()

	rb.shoulds = append(rb.shoulds, e)

	return rb
}

func (rb *RuleBuilder) Because(b rule.Because) ([]rule.Violation, []error) {
	rb.assertNotLocked()

	rb.because = b

	for _, that := range rb.thats {
		if len(that.GetErrors()) > 0 {
			return nil, that.GetErrors()
		}

		that.Evaluate(rb)
	}

	for _, except := range rb.excepts {
		if len(except.GetErrors()) > 0 {
			return nil, except.GetErrors()
		}

		except.Evaluate(rb)
	}

	for _, should := range rb.shoulds {
		if len(should.GetErrors()) > 0 {
			return nil, should.GetErrors()
		}

		if vs := should.Evaluate(rb); len(vs) > 0 {
			rb.violations = append(rb.violations, vs...)
		}
	}

	rb.lock()

	return rb.violations, rb.errors
}

func (rb *RuleBuilder) AddError(err error) {
	rb.assertNotLocked()

	rb.errors = append(rb.errors, err)
}

func (rb *RuleBuilder) GetErrors() []error {
	return rb.errors
}

func (rb *RuleBuilder) GetFiles() []string {
	return rb.files
}

func (rb *RuleBuilder) SetFiles(files []string) {
	rb.assertNotLocked()

	rb.files = files
}

func (rb *RuleBuilder) assertNotLocked() {
	if rb.locked {
		panic("this rule builder has been already used: create a new one if you want to test a new ruleset")
	}
}

func (rb *RuleBuilder) lock() {
	rb.locked = true
}
