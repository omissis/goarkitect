package should_test

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/file/should"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_MatchGlob(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	newRuleBuilder := func() *file.RuleBuilder {
		rb := file.All()
		rb.SetFiles([]string{
			filepath.Join(basePath, "test/project3/baz.go"),
			filepath.Join(basePath, "test/project3/quux.go"),
		})
		return rb
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		glob        string
		want        []rule.Violation
	}{
		{
			desc:        "project3 matches '*.go'",
			ruleBuilder: newRuleBuilder(),
			glob:        "*/*/*.go",
			want:        nil,
		},
		{
			desc:        "project3 matches 'foo/*/*.go'",
			ruleBuilder: newRuleBuilder(),
			glob:        "test/*/*.go",
			want:        nil,
		},
		{
			desc:        "project3 does not match '**/*.ts'",
			ruleBuilder: newRuleBuilder(),
			glob:        "**/*.ts",
			want: []rule.Violation{
				rule.NewViolation("file's path 'baz.go' does not match glob pattern '**/*.ts'"),
				rule.NewViolation("file's path 'quux.go' does not match glob pattern '**/*.ts'"),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mg := should.MatchGlob(tC.glob, basePath)
			got := mg.Evaluate(tC.ruleBuilder)
			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
