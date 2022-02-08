package should

import (
	"goarkitect/internal/arch"
	"fmt"
	"os"
)

func Exist() *ExistExpression {
	return &ExistExpression{}
}

type ExistExpression struct {
}

func (e ExistExpression) Evaluate(rule *arch.RuleSubjectBuilder) []arch.Violation {
	if rule.Kind() != arch.FilesRuleKind {
		panic("ExistExpression can only be used with files rules.")
	}

	violations := make([]arch.Violation, len(rule.Files))
	for _, file := range rule.Files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			violations = append(violations,
				arch.NewViolation(
					fmt.Sprintf("File %s does not exist.", file),
				),
			)
		}
	}

	if len(violations) > 0 {
		return violations
	}

	return nil
}
