package should

import (
	"bytes"
	"fmt"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
	"regexp"
)

func HaveContentMatchingRegex(res string, opts ...Option) *Expression {
	rx := regexp.MustCompile(res)

	return NewExpression(
		func(rb rule.Builder, filePath string) bool {
			data, err := os.ReadFile(filePath)
			if err != nil {
				rb.AddError(err)

				return true
			}

			match := "SINGLE"
			separator := []byte("\n")
			for _, opt := range opts {
				switch opt.(type) {
				case IgnoreNewLinesAtTheEndOfFile:
					data = bytes.TrimRight(data, "\n")
				case IgnoreCase:
					data = bytes.ToLower(data)
				case MatchSingleLines:
					match = "MULTIPLE"
					if sep := opt.(MatchSingleLines).Separator; sep != "" {
						separator = []byte(sep)
					}
				}
			}

			if match == "SINGLE" {
				return !rx.Match(data)
			}

			linesData := bytes.Split(data, separator)
			for _, ld := range linesData {
				if !rx.Match(ld) {
					return true
				}
			}

			return false
		},
		func(filePath string, options options) rule.Violation {
			format := "file '%s' does not have content matching regex '%s'"

			if options.matchSingleLines {
				format = "file '%s' does not have all lines matching regex '%s'"
			}

			if options.negated {
				format = "file '%s' does have content matching regex '%s'"
			}

			if options.negated && options.matchSingleLines {
				format = "file '%s' does have all lines matching regex '%s'"
			}

			return rule.NewViolation(
				fmt.Sprintf(format, filepath.Base(filePath), res),
			)
		},
		opts...,
	)
}
