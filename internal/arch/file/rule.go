package file

import (
	"errors"

	"github.com/omissis/goarkitect/internal/arch/rule"
)

var (
	ErrRuleBuilderLocked = errors.New(
		"this rule builder has been already used: create a new one if you want to test a new ruleset",
	)
	ErrInvalidRuleBuilder = errors.New("invalid rule builder type")
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
	musts      []rule.Expect
	shoulds    []rule.Expect
	coulds     []rule.Expect
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

func (rb *RuleBuilder) Must(e rule.Expect) rule.Builder {
	return rb.AndMust(e)
}

func (rb *RuleBuilder) AndMust(e rule.Expect) rule.Builder {
	if rb.locked {
		rb.addLockError()

		return rb
	}

	rb.musts = append(rb.musts, e)

	return rb
}

func (rb *RuleBuilder) Should(e rule.Expect) rule.Builder {
	return rb.AndShould(e)
}

func (rb *RuleBuilder) AndShould(e rule.Expect) rule.Builder {
	if rb.locked {
		rb.addLockError()

		return rb
	}

	rb.shoulds = append(rb.shoulds, e)

	return rb
}

func (rb *RuleBuilder) Could(e rule.Expect) rule.Builder {
	return rb.AndCould(e)
}

func (rb *RuleBuilder) AndCould(e rule.Expect) rule.Builder {
	if rb.locked {
		rb.addLockError()

		return rb
	}

	rb.coulds = append(rb.coulds, e)

	return rb
}

func (rb *RuleBuilder) Because(b rule.Because) ([]rule.Violation, []error) {
	if rb.locked {
		rb.addLockError()

		return nil, rb.GetErrors()
	}

	rb.lock()

	rb.because = b

	if errs := rb.evaluateThats(); len(errs) > 0 {
		return nil, errs
	}

	if errs := rb.evaluateExcepts(); len(errs) > 0 {
		return nil, errs
	}

	if errs := rb.evaluateExpects(rb.musts, rule.Error); len(errs) > 0 {
		return nil, errs
	}

	if errs := rb.evaluateExpects(rb.shoulds, rule.Warning); len(errs) > 0 {
		return nil, errs
	}

	if errs := rb.evaluateExpects(rb.coulds, rule.Info); len(errs) > 0 {
		return nil, errs
	}

	return rb.violations, rb.errors
}

func (rb *RuleBuilder) evaluateThats() []error {
	for _, that := range rb.thats {
		if len(that.GetErrors()) > 0 {
			return that.GetErrors()
		}

		that.Evaluate(rb)
	}

	return nil
}

func (rb *RuleBuilder) evaluateExcepts() []error {
	for _, except := range rb.excepts {
		if len(except.GetErrors()) > 0 {
			return except.GetErrors()
		}

		except.Evaluate(rb)
	}

	return nil
}

func (rb *RuleBuilder) evaluateExpects(es []rule.Expect, severity rule.Severity) []error {
	for _, must := range es {
		if len(must.GetErrors()) > 0 {
			return must.GetErrors()
		}

		if vs := must.Evaluate(rb); len(vs) > 0 {
			rb.violations = append(rb.violations, rb.wrapCoreViolations(vs, severity)...)
		}
	}

	return nil
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

func (rb *RuleBuilder) wrapCoreViolations(cvs []rule.CoreViolation, severity rule.Severity) []rule.Violation {
	vs := make([]rule.Violation, len(cvs))

	for i, cv := range cvs {
		vs[i] = rule.NewViolation(cv.String(), severity)
	}

	return vs
}
