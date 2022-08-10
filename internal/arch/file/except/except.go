package except

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
)

type Expression interface {
	Evaluate(rb rule.Builder)
}

type evaluateFunc func(filePath string) bool

type baseExpression struct{}

func (e baseExpression) evaluate(rb rule.Builder, eval evaluateFunc) {
	frb := rb.(*file.RuleBuilder)

	nf := make([]string, 0)
	for _, filePath := range frb.GetFiles() {
		if eval(filePath) {
			nf = append(nf, filePath)
		}
	}

	frb.SetFiles(nf)
}
