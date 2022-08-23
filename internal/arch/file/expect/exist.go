package expect

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/omissis/goarkitect/internal/arch/rule"
)

func Exist(opts ...Option) *existExpression {
	expr := &existExpression{}

	for _, opt := range opts {
		opt.apply(&expr.options)
	}

	return expr
}

type existExpression struct {
	baseExpression
}

func (e existExpression) Evaluate(rb rule.Builder) []rule.CoreViolation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e existExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if !os.IsNotExist(err) {
			rb.AddError(err)
		}

		return true
	}

	return false
}

func (e existExpression) getViolation(filePath string) rule.CoreViolation {
	format := "file '%s' does not exist"
	if e.options.negated {
		format = "file '%s' does exist"
	}

	return rule.NewCoreViolation(
		fmt.Sprintf(format, filepath.Base(filePath)),
	)
}
