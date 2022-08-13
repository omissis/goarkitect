package file_test

import (
	"goarkitect/internal/arch/file"
	fe "goarkitect/internal/arch/file/except"
	fs "goarkitect/internal/arch/file/should"
	"goarkitect/internal/arch/rule"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_It_Checks_All_Conditions(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc           string
		ruleBuilder    rule.Builder
		wantViolations []rule.Violation
	}{
		{
			desc:           "check that a file matches all conditions",
			ruleBuilder:    file.One("./test/one/Testfile"),
			wantViolations: nil,
		},
		{
			desc: "check that a set of files matches all conditions",
			ruleBuilder: file.Set(
				"./test/set/Test1file",
				"./test/set/Test2file",
				"./test/set/Test3file",
			),
			wantViolations: nil,
		},
		{
			desc:           "check that all files in folder except one match all conditions",
			ruleBuilder:    file.All().Except(fe.This("./test/set/Test3file")),
			wantViolations: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := tC.ruleBuilder.
				Should(fs.Not(fs.BeGitencrypted())).
				Should(fs.Not(fs.BeGitignored())).
				Should(fs.ContainValue([]byte("foo"))).
				Should(fs.EndWith("file")).
				Should(fs.Exist()).
				Should(fs.HaveContentMatching([]byte("foo"), fs.IgnoreNewLinesAtTheEndOfFile{})).
				Should(fs.HaveContentMatchingRegex("[0-9]+", fs.IgnoreNewLinesAtTheEndOfFile{})).
				Should(fs.HavePermissions("-rwxr-xr-x")).
				Should(fs.MatchGlob("test/one/*", basePath)).
				Should(fs.MatchRegex("[A-z0-9]+")).
				Should(fs.StartWith("Test")).
				Because("I want to test all expressions together")

			if !cmp.Equal(vs, tC.wantViolations, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("Expected %v, got %v", tC.wantViolations, vs)
			}

			if errs != nil {
				t.Errorf("Expected errs to be nil, got: %+v", errs)
			}
		})
	}
}

func Test_It_Adds_ErrRuleBuilderLocked_Only_Once(t *testing.T) {
	rb := file.NewRuleBuilder()

	rb.AddError(file.ErrRuleBuilderLocked)
	rb.AddError(file.ErrRuleBuilderLocked)

	if errs := rb.GetErrors(); len(errs) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errs))
	}
}
