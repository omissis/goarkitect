package should_test

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/file/should"
	"goarkitect/internal/arch/rule"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_NotMatchRegex(t *testing.T) {
	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		regexp      string
		want        []rule.Violation
	}{
		{
			desc:        "foobar matches '[a-z]+'",
			ruleBuilder: file.One("foobar"),
			regexp:      "[a-z]+",
			want: []rule.Violation{
				rule.NewViolation("file's name 'foobar' does not match regex '[a-z]+'"),
			},
		},
		{
			desc:        "foobar matches 'foobar'",
			ruleBuilder: file.One("foobar"),
			regexp:      "foobar",
			want: []rule.Violation{
				rule.NewViolation("file's name 'foobar' does not match regex 'foobar'"),
			},
		},
		{
			desc:        "foobar does not match '[0-9]+'",
			ruleBuilder: file.One("foobar"),
			regexp:      "[0-9]+",
			want:        nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ew := should.NotMatchRegex(tC.regexp)
			got := ew.Evaluate(tC.ruleBuilder)
			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
