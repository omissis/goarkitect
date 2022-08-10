package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
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

func (e existExpression) Evaluate(rb rule.Builder) []rule.Violation {
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

func (e existExpression) getViolation(filePath string) rule.Violation {
	format := "file '%s' does not exist"
	if e.options.negated {
		format = "file '%s' does exist"
	}

	return rule.NewViolation(
		fmt.Sprintf(format, filepath.Base(filePath)),
	)
}
