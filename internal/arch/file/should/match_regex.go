package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"path/filepath"
	"regexp"
)

func MatchRegex(res string, opts ...Option) *matchRegexExpression {
	expr := &matchRegexExpression{
		regex: regexp.MustCompile(res),
	}

	for _, opt := range opts {
		opt.apply(&expr.options)
	}

	return expr
}

type matchRegexExpression struct {
	baseExpression

	regex *regexp.Regexp
}

func (e matchRegexExpression) Evaluate(rb rule.Builder) []rule.Violation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e matchRegexExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	return !e.regex.MatchString(
		filepath.Base(filePath),
	)
}

func (e matchRegexExpression) getViolation(filePath string) rule.Violation {
	format := "file's name '%s' does not match regex '%s'"
	if e.options.negated {
		format = "file's name '%s' does match regex '%s'"
	}

	return rule.NewViolation(
		fmt.Sprintf(
			format,
			filepath.Base(filePath),
			e.regex,
		),
	)
}
