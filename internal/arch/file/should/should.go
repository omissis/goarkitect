package should

import (
	"fmt"
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
)

var (
	ErrEmptyOpts = fmt.Errorf("empty options")
)

type Expression interface {
	Evaluate(rb rule.Builder) []rule.Violation
	applyOption(opt Option) error
	doEvaluate(rb rule.Builder, filePath string) bool
	getViolation(filePath string) rule.Violation
}

type evaluateFunc func(rb rule.Builder, filePath string) bool
type getViolationFunc func(filePath string) rule.Violation
type options struct {
	negated                      bool
	ignoreCase                   bool
	ignoreNewLinesAtTheEndOfFile bool
	matchSingleLines             bool
	matchSingleLinesSeparator    string
}

type baseExpression struct {
	options      options
	getViolation getViolationFunc
}

func (e *baseExpression) evaluate(
	rb rule.Builder,
	evaluate evaluateFunc,
	getViolation getViolationFunc,
) []rule.Violation {
	violations := make([]rule.Violation, 0)
	for _, fp := range rb.(*file.RuleBuilder).GetFiles() {
		result := evaluate(rb, fp)
		if (!e.options.negated && result) || (e.options.negated && !result) {
			violations = append(violations, getViolation(fp))
		}
	}

	return violations
}

func (e *baseExpression) applyOption(opt Option) error {
	return opt.apply(&e.options)
}

func Not(expr Expression) Expression {
	if err := expr.applyOption(Negated{}); err != nil {
		panic(err)
	}

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
