package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
)

func MatchGlob(glob string, basePath string, opts ...Option) *Expression {
	return NewExpression(
		func(rb rule.Builder, filePath string) bool {
			match, err := filepath.Match(filepath.Join(basePath, glob), filePath)
			if err != nil {
				rb.AddError(err)
			}

			return !match
		},
		func(filePath string, options options) rule.Violation {
			format := "file's path '%s' does not match glob pattern '%s'"
			if options.negated {
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
		opts...,
	)
}
