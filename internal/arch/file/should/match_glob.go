package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
)

func MatchGlob(glob string, basePath string) *Expression {
	return &Expression{
		evaluate: func(rb rule.Builder, filePath string) bool {
			match, err := filepath.Match(filepath.Join(basePath, glob), filePath)
			if err != nil {
				rb.AddError(err)
			}

			return !match
		},
		getViolation: func(filePath string, negated bool) rule.Violation {
			format := "file's path '%s' does not match glob pattern '%s'"
			if negated {
				format = "file's path '%s' does match glob pattern '%s'"
			}

			return rule.NewViolation(
				fmt.Sprintf(
					format,
					filepath.Base(filePath),
					glob,
				),
			)
		},
	}
}
