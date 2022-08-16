package except

import (
	"path/filepath"
	"strings"

	"goarkitect/internal/arch/rule"
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
