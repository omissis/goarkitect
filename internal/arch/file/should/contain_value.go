package should

import (
	"bytes"
	"fmt"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
)

func ContainValue(value []byte, opts ...Option) *containValueExpression {
	expr := &containValueExpression{
		value: value,
	}

	for _, opt := range opts {
		opt.apply(&expr.options)
	}

	return expr
}

type containValueExpression struct {
	baseExpression

	value []byte
}

func (e containValueExpression) Evaluate(rb rule.Builder) []rule.Violation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e containValueExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		rb.AddError(err)

		return true
	}

	if e.options.ignoreCase {
		data = bytes.ToLower(data)
		e.value = bytes.ToLower(e.value)
	}

	return !bytes.Contains(data, e.value)
}

func (e containValueExpression) getViolation(filePath string) rule.Violation {
	format := "file '%s' does not contain the value '%s'"

	if e.options.negated {
		format = "file '%s' does contain the value '%s'"
	}

	return rule.NewViolation(
		fmt.Sprintf(format, filepath.Base(filePath), e.value),
	)
}
