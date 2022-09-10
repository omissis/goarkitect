package expect

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/omissis/goarkitect/internal/arch/rule"
)

func HaveContentMatchingRegex(regex string, opts ...Option) *haveContentMatchingRegex {
	expr := &haveContentMatchingRegex{
		regex: regex,
	}

	expr.applyOptions(opts)

	return expr
}

type haveContentMatchingRegex struct {
	baseExpression

	regex string
}

func (e haveContentMatchingRegex) Evaluate(rb rule.Builder) []rule.CoreViolation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e haveContentMatchingRegex) doEvaluate(rb rule.Builder, filePath string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		rb.AddError(err)

		return true
	}

	rx := regexp.MustCompile(e.regex)

	if e.options.ignoreNewLinesAtTheEndOfFile {
		data = bytes.TrimRight(data, "\n")
	}

	if e.options.ignoreCase {
		data = bytes.ToLower(data)
	}

	if e.options.matchSingleLines {
		linesData := bytes.Split(data, []byte(e.options.matchSingleLinesSeparator))
		for _, ld := range linesData {
			if !rx.Match(ld) {
				return true
			}
		}
	}

	return !rx.Match(data)
}

func (e haveContentMatchingRegex) getViolation(filePath string) rule.CoreViolation {
	format := "file '%s' does not have content matching regex '%s'"

	if e.options.matchSingleLines {
		format = "file '%s' does not have all lines matching regex '%s'"
	}

	if e.options.negated {
		format = "file '%s' does have content matching regex '%s'"
	}

	if e.options.negated && e.options.matchSingleLines {
		format = "file '%s' does have all lines matching regex '%s'"
	}

	return rule.NewCoreViolation(
		fmt.Sprintf(format, filepath.Base(filePath), e.regex),
	)
}
