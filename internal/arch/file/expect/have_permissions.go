package expect

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/omissis/goarkitect/internal/arch/rule"
)

var ErrInvalidPermissions = errors.New(
	"permissions must only contain the following characters: 'd', 'r', 'w', 'x', '-'",
)

func HavePermissions(permissions string, opts ...Option) *havePermissionsExpression {
	rx := regexp.MustCompile("^[d-][rwx-]{9}$")

	errs := make([]error, 0)
	if !rx.MatchString(permissions) {
		errs = append(errs, ErrInvalidPermissions)
	}

	expr := &havePermissionsExpression{
		permissions: permissions,
		baseExpression: baseExpression{
			errors: errs,
		},
	}

	expr.applyOptions(opts)

	return expr
}

type havePermissionsExpression struct {
	baseExpression

	permissions string
}

func (e havePermissionsExpression) Evaluate(rb rule.Builder) []rule.CoreViolation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e havePermissionsExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		rb.AddError(err)

		return true
	}

	return e.permissions != info.Mode().String()
}

func (e havePermissionsExpression) getViolation(filePath string) rule.CoreViolation {
	info, err := os.Stat(filePath)
	if err != nil {
		e.errors = append(e.errors, err)

		return rule.CoreViolation{}
	}

	iNodeType := "file"
	if info.IsDir() {
		iNodeType = "directory"
	}

	format := "%s '%s' does not have permissions matching '%s', '%s' found"

	if e.options.negated {
		format = "%s '%s' does have permissions matching '%s', '%s' found"
	}

	return rule.NewCoreViolation(
		fmt.Sprintf(format, iNodeType, filepath.Base(filePath), e.permissions, info.Mode().String()),
	)
}
