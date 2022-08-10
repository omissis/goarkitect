package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
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

func (e matchGlobExpression) Evaluate(rb rule.Builder) []rule.Violation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e matchGlobExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	match, err := filepath.Match(filepath.Join(e.basePath, e.glob), filePath)
	if err != nil {
		rb.AddError(err)
	}

	return !match
}

func (e matchGlobExpression) getViolation(filePath string) rule.Violation {
	format := "file's path '%s' does not match glob pattern '%s'"
	if e.options.negated {
		format = "file's path '%s' does match glob pattern '%s'"
	}

	return rule.NewViolation(
		fmt.Sprintf(
			format,
			filepath.Base(filePath),
			e.glob,
		),
	)
}
