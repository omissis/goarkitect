package that

import (
	"strings"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func EndWith(s string) *EndWithExpression {
	return &EndWithExpression{
		suffix: s,
	}
}

type EndWithExpression struct {
	suffix string

	errors []error
}

func (e *EndWithExpression) GetErrors() []error {
	return e.errors
}

func (e EndWithExpression) Evaluate(rb rule.Builder) {
	frb, ok := rb.(*file.RuleBuilder)
	if !ok {
		e.errors = append(e.errors, file.ErrInvalidRuleBuilder)

		return
	}

	files := make([]string, 0)

	for _, f := range frb.GetFiles() {
		if strings.HasSuffix(f, e.suffix) {
			files = append(files, f)
		}
	}

	frb.SetFiles(files)
}
