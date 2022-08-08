package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
)

func StartWith(prefix string) *Expression {
	le := len(prefix)

	return &Expression{
		checkViolation: func(filePath string) bool {
			fileName := filepath.Base(filePath)
			lf := len(fileName)

			return le <= lf && fileName[:le] != prefix
		},
		getViolation: func(filePath string) rule.Violation {
			return rule.NewViolation(
				fmt.Sprintf(
					"file's name '%s' does not start with '%s'",
					filepath.Base(filePath),
					prefix,
				),
			)
		},
	}
}
