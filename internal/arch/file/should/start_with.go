package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
)

func StartWith(prefix string, opts ...Option) *Expression {
	le := len(prefix)

	return NewExpression(
		func(_ rule.Builder, filePath string) bool {
			fileName := filepath.Base(filePath)
			lf := len(fileName)

			return le <= lf && fileName[:le] != prefix
		},
		func(filePath string, options options) rule.Violation {
			format := "file's name '%s' does not start with '%s'"
			if options.negated {
				format = "file's name '%s' does start with '%s'"
			}

			return rule.NewViolation(
				fmt.Sprintf(format, filepath.Base(filePath), prefix),
			)
		},
		opts...,
	)
}
