package except

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
)

type evaluateFunc func(filePath string) bool

func NewExpression(
	evaluate evaluateFunc,
) *Expression {
	return &Expression{
		evaluate: evaluate,
	}
}

type Expression struct {
	evaluate func(filePath string) bool
}

func (e Expression) Evaluate(rb rule.Builder) {
	frb := rb.(*file.RuleBuilder)

	nf := make([]string, 0)
	for _, filePath := range frb.GetFiles() {
		if e.evaluate(filePath) {
			nf = append(nf, filePath)
		}
	}

	frb.SetFiles(nf)
}
