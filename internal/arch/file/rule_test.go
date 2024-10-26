package file_test

import (
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/omissis/goarkitect/internal/arch/file"
	fe "github.com/omissis/goarkitect/internal/arch/file/except"
	fs "github.com/omissis/goarkitect/internal/arch/file/expect"
	ft "github.com/omissis/goarkitect/internal/arch/file/that"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func Test_It_Checks_All_Conditions(t *testing.T) {
	t.Parallel()

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
			ruleBuilder:    file.All().That(ft.AreInFolder("./test/all", false)).Except(fe.This("./test/all/Test3file")),
			wantViolations: nil,
		},
	}

	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			if _, err := exec.LookPath("git-crypt"); err == nil {
				tC.ruleBuilder.Should(fs.Not(fs.BeGitencrypted()))
			}

			vs, errs := tC.ruleBuilder.
				AndShould(fs.Not(fs.BeGitignored())).
				AndShould(fs.ContainValue([]byte("foo"))).
				AndShould(fs.EndWith("file")).
				AndShould(fs.Exist()).
				AndShould(fs.HaveContentMatching([]byte("foo"), fs.IgnoreNewLinesAtTheEndOfFile{})).
				AndShould(fs.HaveContentMatchingRegex("[A-z0-9]+", fs.IgnoreNewLinesAtTheEndOfFile{})).
				AndShould(fs.HavePermissions("-rw-r--r--")).
				AndShould(fs.MatchGlob("test/*/*", ".")).
				AndShould(fs.MatchRegex("[A-z0-9]+")).
				AndShould(fs.StartWith("Test")).
				Because("I want to test all expressions together")

			if !cmp.Equal(vs, tC.wantViolations, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("expected %v, got %v", tC.wantViolations, vs)
			}

			if errs != nil {
				t.Errorf("expected errs to be nil, got: %+v", errs)
			}
		})
	}
}

func Test_It_Adds_ErrRuleBuilderLocked_Only_Once(t *testing.T) {
	t.Parallel()

	rb := file.NewRuleBuilder()

	rb.AddError(file.ErrRuleBuilderLocked)
	rb.AddError(file.ErrRuleBuilderLocked)

	if errs := rb.GetErrors(); len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
}
