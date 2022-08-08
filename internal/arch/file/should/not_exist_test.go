package should_test

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/file/should"
	"goarkitect/internal/arch/rule"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_NotExist(t *testing.T) {
	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		want        []rule.Violation
	}{
		{
			desc:        "exist.go exists",
			ruleBuilder: file.One("exist.go"),
			want: []rule.Violation{
				rule.NewViolation("file 'exist.go' does exist"),
			},
		},
		{
			desc:        "abc.xyz does not exist",
			ruleBuilder: file.One("abc.xyz"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			e := should.NotExist()
			got := e.Evaluate(tC.ruleBuilder)
			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
