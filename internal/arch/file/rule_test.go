package file_test

import (
	"goarkitect/internal/arch/file"
	"testing"
)

func Test_It_Adds_ErrRuleBuilderLocked_Only_Once(t *testing.T) {
	rb := file.NewRuleBuilder()

	rb.AddError(file.ErrRuleBuilderLocked)
	rb.AddError(file.ErrRuleBuilderLocked)

	if errs := rb.GetErrors(); len(errs) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errs))
	}
}
