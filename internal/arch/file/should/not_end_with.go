package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
)

func NotEndWith(suffix string) *Expression {
	ls := len(suffix)

	return &Expression{
		checkViolation: func(filePath string) bool {
			fileName := filepath.Base(filePath)

			lf := len(fileName)

			return ls > lf || fileName[lf-ls:] == suffix
		},
		getViolation: func(filePath string) rule.Violation {
			return rule.NewViolation(
				fmt.Sprintf(
					"file's name '%s' does end with '%s'",
					filepath.Base(filePath),
					suffix,
				),
			)
		},
	}
}
