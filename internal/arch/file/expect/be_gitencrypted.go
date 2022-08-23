package expect

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/omissis/goarkitect/internal/arch/rule"
)

func BeGitencrypted(opts ...Option) *gitEncryptedExpression {
	expr := &gitEncryptedExpression{}

	for _, opt := range opts {
		opt.apply(&expr.options)
	}

	return expr
}

type gitEncryptedExpression struct {
	baseExpression
}

func (e gitEncryptedExpression) Evaluate(rb rule.Builder) []rule.CoreViolation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e gitEncryptedExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	cmd := exec.Command("git", "crypt", "status", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		rb.AddError(err)

		return true
	}

	return bytes.Contains(out, []byte("not encrypted"))
}

func (e gitEncryptedExpression) getViolation(filePath string) rule.CoreViolation {
	format := "file '%s' is not gitencrypted"

	if e.options.negated {
		format = "file '%s' is gitencrypted"
	}

	return rule.NewCoreViolation(
		fmt.Sprintf(format, filepath.Base(filePath)),
	)
}
