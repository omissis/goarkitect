package should

import (
	"errors"
	"fmt"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ErrInvalidPermissions = errors.New(
		"permissions must only contain the following characters: 'd', 'r', 'w', 'x', '-'",
	)
)

func HavePermissions(want string, opts ...Option) *havePermissionsExpression {
	rx := regexp.MustCompile("[d-][rwx-]{9}")

	errs := make([]error, 0)
	if !rx.MatchString(want) {
		errs = append(errs, ErrInvalidPermissions)
	}

	expr := &havePermissionsExpression{
		want: want,
		baseExpression: baseExpression{
			errors: errs,
		},
	}

	for _, opt := range opts {
		opt.apply(&expr.options)
	}

	return expr
}

type havePermissionsExpression struct {
	baseExpression

	want string
}

func (e havePermissionsExpression) Evaluate(rb rule.Builder) []rule.Violation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e havePermissionsExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		rb.AddError(err)

		return true
	}

	return e.want != info.Mode().String()
}

func (e havePermissionsExpression) getViolation(filePath string) rule.Violation {
	iNodeType := "file"
	if info, _ := os.Stat(filePath); info.IsDir() {
		iNodeType = "directory"
	}

	format := "%s '%s' does not have permissions matching '%s'"

	if e.options.negated {
		format = "%s '%s' does have permissions matching '%s'"
	}

	return rule.NewViolation(
		fmt.Sprintf(format, iNodeType, filepath.Base(filePath), e.want),
	)
}
