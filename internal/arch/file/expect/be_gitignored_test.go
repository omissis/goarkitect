package expect_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/file/expect"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func Test_BeGitignored(t *testing.T) {
	t.Parallel()

	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		options     []expect.Option
		want        []rule.CoreViolation
	}{
		{
			desc:        "file 'ignored.txt' expect be gitignored",
			ruleBuilder: file.One(filepath.Join(basePath, "test/ignored.txt")),
			want:        nil,
		},
		{
			desc:        "file 'not_ignored.txt' expect be gitignored",
			ruleBuilder: file.One(filepath.Join(basePath, "test/not_ignored.txt")),
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file 'not_ignored.txt' is not gitignored"),
			},
		},
		{
			desc:        "negated: file 'ignored.txt' expect not be gitignored",
			ruleBuilder: file.One(filepath.Join(basePath, "test/ignored.txt")),
			options:     []expect.Option{expect.Negated{}},
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file 'ignored.txt' is gitignored"),
			},
		},
		{
			desc:        "negated: file 'not_ignored.txt' expect not be gitignored",
			ruleBuilder: file.One(filepath.Join(basePath, "test/not_ignored.txt")),
			options:     []expect.Option{expect.Negated{}},
			want:        nil,
		},
	}

	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			hcm := expect.BeGitignored(tC.options...)
			got := hcm.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.CoreViolation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
