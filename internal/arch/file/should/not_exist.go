package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
)

func NotExist() *Expression {
	return &Expression{
		checkViolation: func(filePath string) bool {
			if _, err := os.Stat(filePath); err != nil {
				if os.IsNotExist(err) {
					return false
				}
				panic(err)
			}

			return true
		},
		getViolation: func(filePath string) rule.Violation {
			return rule.NewViolation(
				fmt.Sprintf(
					"file '%s' does exist",
					filepath.Base(filePath),
				),
			)
		},
	}
}
