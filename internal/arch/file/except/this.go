package except

import (
	"goarkitect/internal/arch/rule"
	"path/filepath"
	"strings"
)

func This(filePath string) *ThisExpression {
	return &ThisExpression{
		filePath: filePath,
	}
}

type ThisExpression struct {
	baseExpression

	filePath string
}

func (e ThisExpression) Evaluate(rb rule.Builder) {
	e.evaluate(rb, func(filePath string) bool {
		if filepath.IsAbs(e.filePath) {
			absFilePath, err := filepath.Abs(filePath)
			if err != nil {
				rb.AddError(err)

				return true
			}

			return absFilePath != e.filePath
		}

		return !strings.HasSuffix(filepath.Clean(filePath), filepath.Clean(e.filePath))
	})
}
