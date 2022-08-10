package should

import (
	"bytes"
	"fmt"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
	"regexp"
)

func HaveContentMatchingRegex(regex string, opts ...Option) *haveContentMatchingRegex {
	expr := &haveContentMatchingRegex{
		regex: regex,
	}

	for _, opt := range opts {
		opt.apply(&expr.options)
	}

	return expr
}

type haveContentMatchingRegex struct {
	baseExpression

	regex string
}

func (e haveContentMatchingRegex) Evaluate(rb rule.Builder) []rule.Violation {
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

func (e haveContentMatchingRegex) getViolation(filePath string) rule.Violation {
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

	return rule.NewViolation(
		fmt.Sprintf(format, filepath.Base(filePath), e.regex),
	)
}
