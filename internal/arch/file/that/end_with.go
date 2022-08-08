package that

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
	"strings"
)

func EndWith(s string) *EndWithExpression {
	return &EndWithExpression{
		suffix: s,
	}
}

type EndWithExpression struct {
	suffix string
}

func (e EndWithExpression) Evaluate(rb rule.Builder) {
	frb := rb.(*file.RuleBuilder)

	files := make([]string, 0)
	for _, f := range frb.GetFiles() {
		if strings.HasSuffix(f, e.suffix) {
			files = append(files, f)
		}
	}

	frb.SetFiles(files)
}
