package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
	"regexp"
)

func MatchRegex(res string) *Expression {
	rx := regexp.MustCompile(res)

	return &Expression{
		evaluate: func(_ rule.Builder, filePath string) bool {
			return !rx.MatchString(
				filepath.Base(filePath),
			)
		},
		getViolation: func(filePath string, negated bool) rule.Violation {
			format := "file's name '%s' does not match regex '%s'"
			if negated {
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
	}
}
