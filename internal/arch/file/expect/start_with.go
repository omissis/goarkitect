package expect

import (
	"fmt"
	"path/filepath"

	"github.com/omissis/goarkitect/internal/arch/rule"
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

func (e startWithExpression) Evaluate(rb rule.Builder) []rule.CoreViolation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e startWithExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	fileName := filepath.Base(filePath)

	le := len(e.prefix)
	lf := len(fileName)

	return le <= lf && fileName[:le] != e.prefix
}

func (e startWithExpression) getViolation(filePath string) rule.CoreViolation {
	format := "file's name '%s' does not start with '%s'"
	if e.options.negated {
		format = "file's name '%s' does start with '%s'"
	}

	return rule.NewCoreViolation(
		fmt.Sprintf(format, filepath.Base(filePath), e.prefix),
	)
}
