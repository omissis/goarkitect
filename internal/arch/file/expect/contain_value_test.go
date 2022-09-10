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

func Test_ContainValue(t *testing.T) {
	t.Parallel()

	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		value       string
		options     []expect.Option
		want        []rule.CoreViolation
	}{
		{
			desc:        "file 'foobar.txt' contains the value 'bar'",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			value:       "bar",
			want:        nil,
		},
		{
			desc:        "file 'foobar.txt' contains the value 'bar', ignoring case",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			value:       "BAR",
			options: []expect.Option{
				expect.IgnoreCase{},
			},
			want: nil,
		},
		{
			desc:        "file 'foobar.txt' does not contain the value 'something else'",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			value:       "something else",
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file 'foobar.txt' does not contain the value 'something else'"),
			},
		},
		{
			desc:        "negated: file 'foobar.txt' contains the value 'bar'",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			value:       "bar",
			options:     []expect.Option{expect.Negated{}},
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file 'foobar.txt' does contain the value 'bar'"),
			},
		},
		{
			desc:        "negated: file 'foobar.txt' contains the value 'bar', ignoring case",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			value:       "BAR",
			options: []expect.Option{
				expect.IgnoreCase{},
				expect.Negated{},
			},
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file 'foobar.txt' does contain the value 'BAR'"),
			},
		},
		{
			desc:        "negated: file 'foobar.txt' does not contain the value 'something else'",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			value:       "something else",
			options:     []expect.Option{expect.Negated{}},
			want:        nil,
		},
	}

	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			hcm := expect.ContainValue([]byte(tC.value), tC.options...)
			got := hcm.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.CoreViolation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
