package except

import "path/filepath"

func This(value string) *Expression {
	return &Expression{
		evaluate: func(filePath string) bool {
			return filepath.Base(filePath) != value
		},
	}
}
