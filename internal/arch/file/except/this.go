package except

import (
	"goarkitect/internal/arch/rule"
	"path/filepath"
)

func This(value string) *ThisExpression {
	return &ThisExpression{
		value: value,
	}
}

type ThisExpression struct {
	baseExpression

	value string
}

func (e ThisExpression) Evaluate(rb rule.Builder) {
	e.evaluate(rb, func(filePath string) bool {
		return filepath.Base(filePath) != e.value
	})
}
