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

func Test_BeGitignored(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		options     []should.Option
		want        []rule.Violation
	}{
		{
			desc:        "file 'ignored.txt' should be gitignored",
			ruleBuilder: file.One(filepath.Join(basePath, "test/ignored.txt")),
			want:        nil,
		},
		{
			desc:        "file 'not_ignored.txt' should be gitignored",
			ruleBuilder: file.One(filepath.Join(basePath, "test/not_ignored.txt")),
			want: []rule.Violation{
				rule.NewViolation("file 'not_ignored.txt' is not gitignored"),
			},
		},
		{
			desc:        "negated: file 'ignored.txt' should not be gitignored",
			ruleBuilder: file.One(filepath.Join(basePath, "test/ignored.txt")),
			options:     []should.Option{should.Negated{}},
			want: []rule.Violation{
				rule.NewViolation("file 'ignored.txt' is gitignored"),
			},
		},
		{
			desc:        "negated: file 'not_ignored.txt' should not be gitignored",
			ruleBuilder: file.One(filepath.Join(basePath, "test/not_ignored.txt")),
			options:     []should.Option{should.Negated{}},
			want:        nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			hcm := should.BeGitignored(tC.options...)
			got := hcm.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
