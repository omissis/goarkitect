package should

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
)

type evaluateFunc func(rb rule.Builder, filePath string) bool
type getViolationFunc func(filePath string, options options) rule.Violation
type options struct {
	negated                      bool
	ignoreCase                   bool
	ignoreNewLinesAtTheEndOfFile bool
	matchSingleLines             bool
}

func NewExpression(
	evaluate evaluateFunc,
	getViolation getViolationFunc,
	opts ...Option,
) *Expression {
	expr := &Expression{
		evaluate:     evaluate,
		getViolation: getViolation,
	}

	for _, opt := range opts {
		opt.apply(expr)
	}

	return expr
}

type Expression struct {
	options      options
	evaluate     evaluateFunc
	getViolation getViolationFunc
}

func (e *Expression) Evaluate(rb rule.Builder) []rule.Violation {
	violations := make([]rule.Violation, 0)
	for _, fp := range rb.(*file.RuleBuilder).GetFiles() {
		result := e.evaluate(rb, fp)
		if (!e.options.negated && result) || (e.options.negated && !result) {
			violations = append(violations, e.getViolation(fp, e.options))
		}
	}

	return violations
}

func Not(expr *Expression) *Expression {
	expr.options.negated = !expr.options.negated

	return expr
}

type Option interface {
	apply(expr *Expression)
}

type Negated struct{}

func (opt Negated) apply(expr *Expression) {
	expr.options.negated = !expr.options.negated
}

type IgnoreCase struct{}

func (opt IgnoreCase) apply(expr *Expression) {
	expr.options.ignoreCase = true
}

type IgnoreNewLinesAtTheEndOfFile struct{}

func (opt IgnoreNewLinesAtTheEndOfFile) apply(expr *Expression) {
	expr.options.ignoreNewLinesAtTheEndOfFile = true
}

type MatchSingleLines struct {
	Separator string
}

func (opt MatchSingleLines) apply(expr *Expression) {
	expr.options.matchSingleLines = true
}
