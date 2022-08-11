package should

import (
	"fmt"
	"goarkitect/internal/arch/rule"
	"os/exec"
	"path/filepath"
)

func BeGitignored(opts ...Option) *gitIgnoredExpression {
	expr := &gitIgnoredExpression{}

	for _, opt := range opts {
		opt.apply(&expr.options)
	}

	return expr
}

type gitIgnoredExpression struct {
	baseExpression
}

func (e gitIgnoredExpression) Evaluate(rb rule.Builder) []rule.Violation {
	return e.evaluate(rb, e.doEvaluate, e.getViolation)
}

func (e gitIgnoredExpression) doEvaluate(rb rule.Builder, filePath string) bool {
	cmd := exec.Command("git", "check-ignore", "-q", filePath)
	if err := cmd.Run(); err != nil {
		switch err.(type) {
		case *exec.ExitError:
			return err.(*exec.ExitError).ExitCode() != 0
		default:
			panic(err)
		}
	}

	return false
}

func (e gitIgnoredExpression) getViolation(filePath string) rule.Violation {
	format := "file '%s' is not gitignored"

	if e.options.negated {
		format = "file '%s' is gitignored"
	}

	return rule.NewViolation(
		fmt.Sprintf(format, filepath.Base(filePath)),
	)
}
