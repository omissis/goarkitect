package file

import (
	"errors"
	"goarkitect/internal/arch/rule"
)

var (
	ErrRuleBuilderLocked = errors.New(
		"this rule builder has been already used: create a new one if you want to test a new ruleset",
	)
)

func All() *RuleBuilder {
	return NewRuleBuilder()
}

func One(filePath string) *RuleBuilder {
	rb := NewRuleBuilder()

	rb.files = []string{filePath}

	return rb
}

func Set(filePaths ...string) *RuleBuilder {
	rb := NewRuleBuilder()

	rb.SetFiles(filePaths)

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
	return rb.AndThat(t)
}

func (rb *RuleBuilder) AndThat(t rule.That) rule.Builder {
	if rb.locked {
		rb.addLockError()
		return rb
	}

	rb.thats = append(rb.thats, t)

	return rb
}

func (rb *RuleBuilder) Except(s ...rule.Except) rule.Builder {
	if rb.locked {
		rb.addLockError()
		return rb
	}

	rb.excepts = s

	return rb
}

func (rb *RuleBuilder) Should(e rule.Should) rule.Builder {
	return rb.AndShould(e)
}

func (rb *RuleBuilder) AndShould(e rule.Should) rule.Builder {
	if rb.locked {
		rb.addLockError()
		return rb
	}

	rb.shoulds = append(rb.shoulds, e)

	return rb
}

func (rb *RuleBuilder) Because(b rule.Because) ([]rule.Violation, []error) {
	if rb.locked {
		rb.addLockError()

		return nil, rb.GetErrors()
	}

	rb.lock()

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

	return rb.violations, rb.errors
}

func (rb *RuleBuilder) AddError(err error) {
	for _, err := range rb.errors {
		if errors.Is(err, ErrRuleBuilderLocked) {
			return
		}
	}

	rb.errors = append(rb.errors, err)
}

func (rb *RuleBuilder) GetErrors() []error {
	return rb.errors
}

func (rb *RuleBuilder) GetFiles() []string {
	return rb.files
}

func (rb *RuleBuilder) SetFiles(files []string) {
	rb.files = files
}

func (rb *RuleBuilder) addLockError() {
	rb.AddError(ErrRuleBuilderLocked)
}

func (rb *RuleBuilder) lock() {
	rb.locked = true
}
