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
		getViolation: func(filePath string) rule.Violation {
			return rule.NewViolation(
				fmt.Sprintf(
					"file's name '%s' does not match regex '%s'",
					filepath.Base(filePath),
					res,
				),
			)
		},
	}
}
