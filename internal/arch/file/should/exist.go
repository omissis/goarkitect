package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
)

func Exist() *Expression {
	return &Expression{
		checkViolation: func(filePath string) bool {
			if _, err := os.Stat(filePath); err != nil {
				if os.IsNotExist(err) {
					return true
				}

				panic(err)
			}

			return false
		},
		getViolation: func(filePath string) rule.Violation {
			return rule.NewViolation(
				fmt.Sprintf(
					"file '%s' does not exist",
					filepath.Base(filePath),
				),
			)
		},
	}
}
