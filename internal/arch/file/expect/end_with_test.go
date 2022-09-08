package expect_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/file/expect"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func Test_EndWith(t *testing.T) {
	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		suffix      string
		options     []expect.Option
		want        []rule.CoreViolation
	}{
		{
			desc:        "foobar ends with bar",
			ruleBuilder: file.One("foobar"),
			suffix:      "bar",
			want:        nil,
		},
		{
			desc:        "foobar does not end with baz",
			ruleBuilder: file.One("foobar"),
			suffix:      "baz",
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file's name 'foobar' does not end with 'baz'"),
			},
		},
		{
			desc:        "negated: foobar does not end with baz",
			ruleBuilder: file.One("foobar"),
			suffix:      "baz",
			options:     []expect.Option{expect.Negated{}},
			want:        nil,
		},
		{
			desc:        "negated: foobar ends with bar",
			ruleBuilder: file.One("foobar"),
			suffix:      "bar",
			options:     []expect.Option{expect.Negated{}},
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file's name 'foobar' does end with 'bar'"),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ew := expect.EndWith(tC.suffix, tC.options...)
			got := ew.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.CoreViolation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
