package expect

import (
	"fmt"
	"path/filepath"

	"github.com/omissis/goarkitect/internal/arch/rule"
)

func EndWith(suffix string, opts ...Option) *endWithExpression {
	expr := &endWithExpression{
		suffix: suffix,
	}

	for _, opt := range opts {
		opt.apply(&expr.options)
	}

	return expr
}

type endWithExpression struct {
	baseExpression

	suffix string
}

func (e endWithExpression) Evaluate(rb rule.Builder) []rule.CoreViolation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e endWithExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	fileName := filepath.Base(filePath)

	ls := len(e.suffix)
	lf := len(fileName)

	return ls <= lf && fileName[lf-ls:] != e.suffix
}

func (e endWithExpression) getViolation(filePath string) rule.CoreViolation {
	format := "file's name '%s' does not end with '%s'"
	if e.options.negated {
		format = "file's name '%s' does end with '%s'"
	}

	return rule.NewCoreViolation(
		fmt.Sprintf(format, filepath.Base(filePath), e.suffix),
	)
}
