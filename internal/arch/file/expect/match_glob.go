package expect

import (
	"fmt"
	"path/filepath"

	"goarkitect/internal/arch/rule"
)

func MatchGlob(glob string, basePath string, opts ...Option) *matchGlobExpression {
	expr := &matchGlobExpression{
		basePath: basePath,
		glob:     glob,
	}

	for _, opt := range opts {
		opt.apply(&expr.options)
	}

	return expr
}

type matchGlobExpression struct {
	baseExpression

	basePath string
	glob     string
}

func (e matchGlobExpression) Evaluate(rb rule.Builder) []rule.CoreViolation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e matchGlobExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	abs1, err := filepath.Abs(filepath.Join(e.basePath, e.glob))
	if err != nil {
		rb.AddError(err)

		return true
	}

	abs2, err := filepath.Abs(filePath)
	if err != nil {
		rb.AddError(err)

		return true
	}

	match, err := filepath.Match(abs1, abs2)
	if err != nil {
		rb.AddError(err)
	}

	return !match
}

func (e matchGlobExpression) getViolation(filePath string) rule.CoreViolation {
	format := "file's path '%s' does not match glob pattern '%s'"
	if e.options.negated {
		format = "file's path '%s' does match glob pattern '%s'"
	}

	return rule.NewCoreViolation(
		fmt.Sprintf(
			format,
			filepath.Base(filePath),
			e.glob,
		),
	)
}
