package expect_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/file/expect"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func Test_StartWith(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		prefix      string
		options     []expect.Option
		want        []rule.CoreViolation
	}{
		{
			desc:        "foobar starts with foo",
			ruleBuilder: file.One("foobar"),
			prefix:      "foo",
			want:        nil,
		},
		{
			desc:        "foobar does not start with baz",
			ruleBuilder: file.One("foobar"),
			prefix:      "baz",
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file's name 'foobar' does not start with 'baz'"),
			},
		},
		{
			desc:        "negated: foobar does not start with baz",
			ruleBuilder: file.One("foobar"),
			prefix:      "baz",
			options:     []expect.Option{expect.Negated{}},
			want:        nil,
		},
		{
			desc:        "negated: foobar starts with foo",
			ruleBuilder: file.One("foobar"),
			prefix:      "foo",
			options:     []expect.Option{expect.Negated{}},
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file's name 'foobar' does start with 'foo'"),
			},
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			ew := expect.StartWith(tC.prefix, tC.options...)
			got := ew.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.CoreViolation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
