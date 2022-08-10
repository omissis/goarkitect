package should

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
)

type Expression struct {
	evaluate     func(rb rule.Builder, filePath string) bool
	getViolation func(filePath string) rule.Violation
}

func (e Expression) Evaluate(rb rule.Builder) []rule.Violation {
	violations := make([]rule.Violation, 0)
	for _, filePath := range rb.(*file.RuleBuilder).GetFiles() {
		if e.evaluate(rb, filePath) {
			violations = append(violations, e.getViolation(filePath))
		}
	}

	return violations
}
