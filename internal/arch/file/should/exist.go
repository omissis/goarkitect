package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
)

func Exist() *Expression {
	return &Expression{
		evaluate: func(rb rule.Builder, filePath string) bool {
			if _, err := os.Stat(filePath); err != nil {
				if os.IsNotExist(err) {
					return true
				}

				rb.AddError(err)
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
