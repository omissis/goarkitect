package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
)

func StartWith(prefix string, opts ...Option) *startWithExpression {
	expr := &startWithExpression{
		prefix: prefix,
	}

	for _, opt := range opts {
		opt.apply(&expr.options)
	}

	return expr
}

type startWithExpression struct {
	baseExpression

	prefix string
}

func (e startWithExpression) Evaluate(rb rule.Builder) []rule.Violation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e startWithExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	fileName := filepath.Base(filePath)

	le := len(e.prefix)
	lf := len(fileName)

	return le <= lf && fileName[:le] != e.prefix
}

func (e startWithExpression) getViolation(filePath string) rule.Violation {
	format := "file's name '%s' does not start with '%s'"
	if e.options.negated {
		format = "file's name '%s' does start with '%s'"
	}

	return rule.NewViolation(
		fmt.Sprintf(format, filepath.Base(filePath), e.prefix),
	)
}
