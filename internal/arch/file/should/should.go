package should

import (
	"fmt"
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
)

// TODO: use errors
var (
	ErrEmptyOpts      = fmt.Errorf("empty options")
	ErrEmptySeparator = fmt.Errorf("empty separator")
)

type Expression interface {
	Evaluate(rb rule.Builder) []rule.Violation
	applyOption(opt Option)
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

func (e *baseExpression) applyOption(opt Option) {
	opt.apply(&e.options)
}

func Not(expr Expression) Expression {
	expr.applyOption(Negated{})

	return expr
}

type Option interface {
	apply(opts *options)
}

type Negated struct{}

func (opt Negated) apply(opts *options) {
	if opts == nil {
		return
	}

	opts.negated = !opts.negated
}

type IgnoreCase struct{}

func (opt IgnoreCase) apply(opts *options) {
	if opts == nil {
		return
	}

	opts.ignoreCase = true
}

type IgnoreNewLinesAtTheEndOfFile struct{}

func (opt IgnoreNewLinesAtTheEndOfFile) apply(opts *options) {
	if opts == nil {
		return
	}

	opts.ignoreNewLinesAtTheEndOfFile = true
}

type MatchSingleLines struct {
	Separator string
}

func (opt MatchSingleLines) apply(opts *options) {
	if opts == nil {
		return
	}

	opts.matchSingleLines = true
	if opt.Separator == "" {
		opts.matchSingleLinesSeparator = "\n"
	} else {
		opts.matchSingleLinesSeparator = opt.Separator
	}
}
