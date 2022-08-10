package should_test

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/file/should"
	"goarkitect/internal/arch/rule"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_StartWith(t *testing.T) {
	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		prefix      string
		want        []rule.Violation
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
			want: []rule.Violation{
				rule.NewViolation("file's name 'foobar' does not start with 'baz'"),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ew := should.StartWith(tC.prefix)
			got := ew.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
