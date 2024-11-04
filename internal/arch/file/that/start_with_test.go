package that_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/file/that"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func Test_StartWith(t *testing.T) {
	t.Parallel()

	rb := func() *file.RuleBuilder {
		rb := file.All()
		rb.SetFiles([]string{"Dockerfile", "Makefile", "foo/bar.go"})

		return rb
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		prefix      string
		want        []string
	}{
		{
			desc:        "files starting with 'foo'",
			ruleBuilder: rb(),
			prefix:      "foo",
			want:        nil,
		},
		{
			desc:        "files starting with 'Make'",
			ruleBuilder: rb(),
			prefix:      "Make",
			want:        []string{"Makefile"},
		},
		{
			desc:        "files in a subdirectory starting with 'bar",
			ruleBuilder: rb(),
			prefix:      "bar",
			want:        []string{"foo/bar.go"},
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			ew := that.StartWith(tC.prefix)
			ew.Evaluate(tC.ruleBuilder)

			got := tC.ruleBuilder.GetFiles()
			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
