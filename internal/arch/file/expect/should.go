package expect

import (
	"fmt"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

var ErrEmptyOpts = fmt.Errorf("empty options")

type Expression interface {
	Evaluate(rb rule.Builder) []rule.CoreViolation
	GetErrors() []error
	applyOptions(opts []Option)
	applyOption(opt Option)
	doEvaluate(rb rule.Builder, filePath string) bool
	getViolation(filePath string) rule.CoreViolation
}

type (
	evaluateFunc     func(rb rule.Builder, filePath string) bool
	getViolationFunc func(filePath string) rule.CoreViolation
	options          struct {
		negated                      bool
		ignoreCase                   bool
		ignoreNewLinesAtTheEndOfFile bool
		matchSingleLines             bool
		matchSingleLinesSeparator    string
	}
)

type baseExpression struct {
	options      options
	getViolation getViolationFunc //nolint:unused // false positive
	errors       []error
}

func (e *baseExpression) GetErrors() []error {
	return e.errors
}

func (e *baseExpression) evaluate(
	rb rule.Builder,
	evaluate evaluateFunc,
	getViolation getViolationFunc,
) []rule.CoreViolation {
	frb, ok := rb.(*file.RuleBuilder)
	if !ok {
		e.errors = append(e.errors, file.ErrInvalidRuleBuilder)

		return nil
	}

	if len(e.errors) > 0 {
		return nil
	}

	violations := make([]rule.CoreViolation, 0)

	for _, fp := range frb.GetFiles() {
		result := evaluate(rb, fp)
		if (!e.options.negated && result) || (e.options.negated && !result) {
			violations = append(violations, getViolation(fp))
		}
	}

	return violations
}

func (e *baseExpression) applyOptions(opts []Option) {
	for _, opt := range opts {
		e.applyOption(opt)
	}
}

func (e *baseExpression) applyOption(opt Option) {
	if err := opt.apply(&e.options); err != nil {
		e.errors = append(e.errors, err)
	}
}

func Not(expr Expression) Expression {
	expr.applyOption(Negated{})

	return expr
}

type Option interface {
	apply(opts *options) error
}

type Negated struct{}

func (opt Negated) apply(opts *options) error {
	if opts == nil {
		return ErrEmptyOpts
	}

	opts.negated = !opts.negated

	return nil
}

type IgnoreCase struct{}

func (opt IgnoreCase) apply(opts *options) error {
	if opts == nil {
		return ErrEmptyOpts
	}

	opts.ignoreCase = true

	return nil
}

type IgnoreNewLinesAtTheEndOfFile struct{}

func (opt IgnoreNewLinesAtTheEndOfFile) apply(opts *options) error {
	if opts == nil {
		return ErrEmptyOpts
	}

	opts.ignoreNewLinesAtTheEndOfFile = true

	return nil
}

type MatchSingleLines struct {
	Separator string
}

func (opt MatchSingleLines) apply(opts *options) error {
	if opts == nil {
		return ErrEmptyOpts
	}

	opts.matchSingleLines = true
	if opt.Separator == "" {
		opts.matchSingleLinesSeparator = "\n"
	} else {
		opts.matchSingleLinesSeparator = opt.Separator
	}

	return nil
}
