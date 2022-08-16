package rule_test

import (
	"testing"

	"goarkitect/internal/arch/rule"
)

func Test_NewViolation(t *testing.T) {
	v := rule.NewViolation("message", rule.Error)
	if v.String() != "message" {
		t.Errorf("NewViolation(\"message\") returns %v, want %v", v.String(), "message")
	}
}

func Test_NewCoreViolation(t *testing.T) {
	v := rule.NewCoreViolation("message")
	if v.String() != "message" {
		t.Errorf("NewViolation(\"message\") returns %v, want %v", v.String(), "message")
	}
}
