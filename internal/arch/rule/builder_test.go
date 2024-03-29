package rule_test

import (
	"testing"

	"github.com/omissis/goarkitect/internal/arch/rule"
)

func Test_NewViolation(t *testing.T) {
	t.Parallel()

	v := rule.NewViolation("message", rule.Error)
	if v.String() != "[ERROR] message" {
		t.Errorf("NewViolation(\"message\") returns %v, want %v", v.String(), "[ERROR] message")
	}
}

func Test_NewCoreViolation(t *testing.T) {
	t.Parallel()

	v := rule.NewCoreViolation("message")
	if v.String() != "message" {
		t.Errorf("NewViolation(\"message\") returns %v, want %v", v.String(), "message")
	}
}
