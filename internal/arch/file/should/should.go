package should

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
)

type Expression struct {
	checkViolation func(filePath string) bool
	getViolation   func(filePath string) rule.Violation
}

func (e Expression) Evaluate(rb rule.Builder) []rule.Violation {
	violations := make([]rule.Violation, 0)
	for _, filePath := range rb.(*file.RuleBuilder).GetFiles() {
		if e.checkViolation(filePath) {
			violations = append(violations, e.getViolation(filePath))
		}
	}

	return violations
}
