package rule_test

import (
	"goarkitect/internal/arch/rule"
	"testing"
)

func Test_NewViolation(t *testing.T) {
	v := rule.NewViolation("message")
	if v.String() != "message" {
		t.Errorf("NewViolation(\"message\") returns %v, want %v", v.String(), "message")
	}
}
