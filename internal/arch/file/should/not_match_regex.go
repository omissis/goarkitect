package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
	"regexp"
)

func NotMatchRegex(res string) *Expression {
	rx := regexp.MustCompile(res)

	return &Expression{
		checkViolation: func(filePath string) bool {
			fileName := filepath.Base(filePath)
			return rx.MatchString(fileName)
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
