package should

import (
	"bytes"
	"fmt"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"
)

func HaveContentMatching(want []byte, opts ...Option) *haveContentMatchingExpression {
	expr := &haveContentMatchingExpression{
		want: want,
	}

	for _, opt := range opts {
		opt.apply(&expr.options)
	}

	return expr
}

type haveContentMatchingExpression struct {
	baseExpression

	want []byte
}

func (e haveContentMatchingExpression) Evaluate(rb rule.Builder) []rule.Violation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e haveContentMatchingExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		rb.AddError(err)

		return true
	}

	if e.options.ignoreNewLinesAtTheEndOfFile {
		data = bytes.TrimRight(data, "\n")
		e.want = bytes.TrimRight(e.want, "\n")
	}

	if e.options.ignoreCase {
		data = bytes.ToLower(data)
		e.want = bytes.ToLower(e.want)
	}

	if e.options.matchSingleLines {
		linesData := bytes.Split(data, []byte(e.options.matchSingleLinesSeparator))
		for _, ld := range linesData {
			if slices.Compare(ld, e.want) != 0 {
				return true
			}
		}

		return false
	}

	return slices.Compare(data, e.want) != 0
}

func (e haveContentMatchingExpression) getViolation(filePath string) rule.Violation {
	format := "file '%s' does not have content matching '%s'"

	if e.options.matchSingleLines {
		format = "file '%s' does not have all lines matching '%s'"
	}

	if e.options.negated {
		format = "file '%s' does have content matching '%s'"
	}

	if e.options.negated && e.options.matchSingleLines {
		format = "file '%s' does have all lines matching '%s'"
	}

	return rule.NewViolation(
		fmt.Sprintf(format, filepath.Base(filePath), e.want),
	)
}
