package should

import (
	"bytes"
	"fmt"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"
)

type ContentMatchOption interface {
	apply(data []byte) []byte
}

type IgnoreNewLinesAtTheEndOfFile struct{}

func (opt IgnoreNewLinesAtTheEndOfFile) apply(data []byte) []byte {
	return bytes.TrimRight(data, "\n")
}

type IgnoreCase struct{}

func (opt IgnoreCase) apply(data []byte) []byte {
	return bytes.ToLower(data)
}

func HaveContentMatching(want []byte, opts ...ContentMatchOption) *Expression {
	return &Expression{
		evaluate: func(rb rule.Builder, filePath string) bool {
			data, err := os.ReadFile(filePath)
			if err != nil {
				rb.AddError(err)
				return true
			}

			for _, opt := range opts {
				data = opt.apply(data)
			}

			for _, opt := range opts {
				want = opt.apply(want)
			}

			return slices.Compare(data, want) != 0
		},
		getViolation: func(filePath string, negated bool) rule.Violation {
			format := "file '%s' does not have content matching '%s'"
			if negated {
				format = "file '%s' does have content matching '%s'"
			}

			return rule.NewViolation(
				fmt.Sprintf(format, filepath.Base(filePath), want),
			)
		},
	}
}
