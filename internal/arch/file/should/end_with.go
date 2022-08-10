package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
)

func EndWith(suffix string) *Expression {
	ls := len(suffix)

	return &Expression{
		evaluate: func(_ rule.Builder, filePath string) bool {
			fileName := filepath.Base(filePath)

			lf := len(fileName)

			return ls <= lf && fileName[lf-ls:] != suffix
		},
		getViolation: func(filePath string, negated bool) rule.Violation {
			format := "file's name '%s' does not end with '%s'"
			if negated {
				format = "file's name '%s' does end with '%s'"
			}

			return rule.NewViolation(
				fmt.Sprintf(format, filepath.Base(filePath), suffix),
			)
		},
	}
}
