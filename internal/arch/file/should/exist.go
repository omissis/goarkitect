package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
)

func Exist(opts ...Option) *Expression {
	return NewExpression(
		func(rb rule.Builder, filePath string) bool {
			if _, err := os.Stat(filePath); err != nil {
				if !os.IsNotExist(err) {
					rb.AddError(err)
				}

				return true
			}

			return false
		},
		func(filePath string, options options) rule.Violation {
			format := "file '%s' does not exist"
			if options.negated {
				format = "file '%s' does exist"
			}

			return rule.NewViolation(
				fmt.Sprintf(format, filepath.Base(filePath)),
			)
		},
		opts...,
	)
}
