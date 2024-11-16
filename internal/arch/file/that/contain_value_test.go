package that_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/file/that"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func Test_ContainValue(t *testing.T) {
	t.Parallel()

	rb := func() *file.RuleBuilder {
		rb := file.All()
		rb.SetFiles([]string{"Dockerfile", "Makefile", "foo/bar.go"})

		return rb
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		value       string
		want        []string
	}{
		{
			desc:        "files contains value 'foo'",
			ruleBuilder: rb(),
			value:       "foo",
			want:        []string{"foo/bar.go"},
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			ew := that.ContainValue(tC.value)
			ew.Evaluate(tC.ruleBuilder)

			got := tC.ruleBuilder.GetFiles()
			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
