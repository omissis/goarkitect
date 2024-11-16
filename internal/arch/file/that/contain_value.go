package that

import (
	"strings"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func ContainsValue(s string) *ContainsValueExpression {
	return &ContainsValueExpression{
		value: s,
	}
}

type ContainsValueExpression struct {
	value string

	errors []error
}

func (e *ContainsValueExpression) GetErrors() []error {
	return e.errors
}

func (e *ContainsValueExpression) Evaluate(rb rule.Builder) {
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
