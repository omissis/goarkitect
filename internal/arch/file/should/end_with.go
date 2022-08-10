package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
)

func EndWith(suffix string, opts ...Option) *Expression {
	ls := len(suffix)

	return NewExpression(
		func(_ rule.Builder, filePath string) bool {
			fileName := filepath.Base(filePath)

			lf := len(fileName)

			return ls <= lf && fileName[lf-ls:] != suffix
		},
		func(filePath string, options options) rule.Violation {
			format := "file's name '%s' does not end with '%s'"
			if options.negated {
				format = "file's name '%s' does end with '%s'"
			}

			return rule.NewViolation(
				fmt.Sprintf(format, filepath.Base(filePath), suffix),
			)
		},
		opts...,
	)
}
