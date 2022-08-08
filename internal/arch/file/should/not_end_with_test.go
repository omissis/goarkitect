package should_test

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/file/should"
	"goarkitect/internal/arch/rule"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_NotEndWith(t *testing.T) {
	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		suffix      string
		want        []rule.Violation
	}{
		{
			desc:        "foobar does not end with baz",
			ruleBuilder: file.One("foobar"),
			suffix:      "baz",
			want:        nil,
		},
		{
			desc:        "foobar ends with bar",
			ruleBuilder: file.One("foobar"),
			suffix:      "bar",
			want: []rule.Violation{
				rule.NewViolation("file's name 'foobar' does end with 'bar'"),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ew := should.NotEndWith(tC.suffix)
			got := ew.Evaluate(tC.ruleBuilder)
			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
