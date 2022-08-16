package that_test

import (
	"testing"

	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/file/that"
	"goarkitect/internal/arch/rule"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_EndWith(t *testing.T) {
	rb := func() *file.RuleBuilder {
		rb := file.All()
		rb.SetFiles([]string{"Dockerfile", "Makefile"})

		return rb
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		suffix      string
		want        []string
	}{
		{
			desc:        "files ending with 'foo'",
			ruleBuilder: rb(),
			suffix:      "foo",
			want:        nil,
		},
		{
			desc:        "files ending with 'file'",
			ruleBuilder: rb(),
			suffix:      "file",
			want:        []string{"Dockerfile", "Makefile"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ew := that.EndWith(tC.suffix)
			ew.Evaluate(tC.ruleBuilder)

			got := tC.ruleBuilder.GetFiles()
			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
