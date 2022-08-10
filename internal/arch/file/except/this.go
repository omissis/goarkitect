package except

import "path/filepath"

func This(value string) *Expression {
	return NewExpression(
		func(filePath string) bool {
			return filepath.Base(filePath) != value
		},
	)
}
