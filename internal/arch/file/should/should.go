package should

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
)

type Expression struct {
	negated      bool
	evaluate     func(rb rule.Builder, filePath string) bool
	getViolation func(filePath string, negated bool) rule.Violation
}

func (e Expression) Evaluate(rb rule.Builder) []rule.Violation {
	violations := make([]rule.Violation, 0)
	for _, fp := range rb.(*file.RuleBuilder).GetFiles() {
		result := e.evaluate(rb, fp)
		if (!e.negated && result) || (e.negated && !result) {
			violations = append(violations, e.getViolation(fp, e.negated))
		}
	}

	return violations
}

func Not(expr *Expression) *Expression {
	expr.negated = !expr.negated

	return expr
}
