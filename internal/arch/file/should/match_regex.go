package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
	"regexp"
)

func MatchRegex(res string, opts ...Option) *Expression {
	rx := regexp.MustCompile(res)

	return NewExpression(
		func(_ rule.Builder, filePath string) bool {
			return !rx.MatchString(
				filepath.Base(filePath),
			)
		},
		func(filePath string, options options) rule.Violation {
			format := "file's name '%s' does not match regex '%s'"
			if options.negated {
				format = "file's name '%s' does match regex '%s'"
			}

			return rule.NewViolation(
				fmt.Sprintf(
					format,
					filepath.Base(filePath),
					res,
				),
			)
		},
		opts...,
	)
}
