package expect

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"goarkitect/internal/arch/rule"
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

func (e containValueExpression) Evaluate(rb rule.Builder) []rule.CoreViolation {
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

func (e containValueExpression) getViolation(filePath string) rule.CoreViolation {
	format := "file '%s' does not contain the value '%s'"

	if e.options.negated {
		format = "file '%s' does contain the value '%s'"
	}

	return rule.NewCoreViolation(
		fmt.Sprintf(format, filepath.Base(filePath), e.value),
	)
}
