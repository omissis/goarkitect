package expect

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"

	"github.com/omissis/goarkitect/internal/arch/rule"
)

func HaveContentMatching(value []byte, opts ...Option) *haveContentMatchingExpression {
	expr := &haveContentMatchingExpression{
		value: value,
	}

	expr.applyOptions(opts)

	return expr
}

type haveContentMatchingExpression struct {
	baseExpression

	value []byte
}

func (e haveContentMatchingExpression) Evaluate(rb rule.Builder) []rule.CoreViolation {
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
		e.value = bytes.TrimRight(e.value, "\n")
	}

	if e.options.ignoreCase {
		data = bytes.ToLower(data)
		e.value = bytes.ToLower(e.value)
	}

	if e.options.matchSingleLines {
		linesData := bytes.Split(data, []byte(e.options.matchSingleLinesSeparator))
		for _, ld := range linesData {
			if slices.Compare(ld, e.value) != 0 {
				return true
			}
		}

		return false
	}

	return slices.Compare(data, e.value) != 0
}

func (e haveContentMatchingExpression) getViolation(filePath string) rule.CoreViolation {
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

	return rule.NewCoreViolation(
		fmt.Sprintf(format, filepath.Base(filePath), e.value),
	)
}
