package that

import (
	"path/filepath"
	"strings"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func StartWith(s string) *StartWithExpression {
	return &StartWithExpression{
		prefix: s,
	}
}

type StartWithExpression struct {
	prefix string

	errors []error
}

func (e *StartWithExpression) GetErrors() []error {
	return e.errors
}

func (e *StartWithExpression) Evaluate(rb rule.Builder) {
	frb, ok := rb.(*file.RuleBuilder)
	if !ok {
		e.errors = append(e.errors, file.ErrInvalidRuleBuilder)

		return
	}

	files := make([]string, 0)

	for _, f := range frb.GetFiles() {
		if strings.HasPrefix(filepath.Base(f), e.prefix) {
			files = append(files, f)
		}
	}

	frb.SetFiles(files)
}
