package expect_test

import (
	"testing"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/file/expect"
	"github.com/omissis/goarkitect/internal/arch/rule"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_Exist(t *testing.T) {
	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		options     []expect.Option
		want        []rule.CoreViolation
	}{
		{
			desc:        "exist.go exists",
			ruleBuilder: file.One("exist.go"),
			want:        nil,
		},
		{
			desc:        "abc.xyz does not exist",
			ruleBuilder: file.One("abc.xyz"),
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file 'abc.xyz' does not exist"),
			},
		},
		{
			desc:        "negated: exist.go exists",
			ruleBuilder: file.One("exist.go"),
			options:     []expect.Option{expect.Negated{}},
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file 'exist.go' does exist"),
			},
		},
		{
			desc:        "negated: abc.xyz does not exist",
			ruleBuilder: file.One("abc.xyz"),
			options:     []expect.Option{expect.Negated{}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			e := expect.Exist(tC.options...)
			got := e.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.CoreViolation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
