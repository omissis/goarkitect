package should_test

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/file/should"
	"goarkitect/internal/arch/rule"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_MatchRegex(t *testing.T) {
	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		regexp      string
		options     []should.Option
		want        []rule.Violation
	}{
		{
			desc:        "foobar matches '[a-z]+'",
			ruleBuilder: file.One("foobar"),
			regexp:      "[a-z]+",
			want:        nil,
		},
		{
			desc:        "foobar matches 'foobar'",
			ruleBuilder: file.One("foobar"),
			regexp:      "foobar",
			want:        nil,
		},
		{
			desc:        "foobar does not match '[0-9]+'",
			ruleBuilder: file.One("foobar"),
			regexp:      "[0-9]+",
			want: []rule.Violation{
				rule.NewViolation("file's name 'foobar' does not match regex '[0-9]+'"),
			},
		},
		{
			desc:        "negated: foobar matches '[a-z]+'",
			ruleBuilder: file.One("foobar"),
			regexp:      "[a-z]+",
			options:     []should.Option{should.Negated{}},
			want: []rule.Violation{
				rule.NewViolation("file's name 'foobar' does match regex '[a-z]+'"),
			},
		},
		{
			desc:        "negated: foobar matches 'foobar'",
			ruleBuilder: file.One("foobar"),
			regexp:      "foobar",
			options:     []should.Option{should.Negated{}},
			want: []rule.Violation{
				rule.NewViolation("file's name 'foobar' does match regex 'foobar'"),
			},
		},
		{
			desc:        "negated: foobar does not match '[0-9]+'",
			ruleBuilder: file.One("foobar"),
			regexp:      "[0-9]+",
			options:     []should.Option{should.Negated{}},
			want:        nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ew := should.MatchRegex(tC.regexp, tC.options...)
			got := ew.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
