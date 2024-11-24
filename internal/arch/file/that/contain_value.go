package that

import (
	"strings"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func ContainValue(s string) *ContainValueExpression {
	return &ContainValueExpression{
		value: s,
	}
}

type ContainValueExpression struct {
	value string

	errors []error
}

func (e *ContainValueExpression) GetErrors() []error {
	return e.errors
}

func (e *ContainValueExpression) Evaluate(rb rule.Builder) {
	frb, ok := rb.(*file.RuleBuilder)
	if !ok {
		e.errors = append(e.errors, file.ErrInvalidRuleBuilder)

		return
	}

	files := make([]string, 0)

	for _, f := range frb.GetFiles() {
		if strings.Contains(f, e.value) {
			files = append(files, f)
		}
	}

	frb.SetFiles(files)
}
