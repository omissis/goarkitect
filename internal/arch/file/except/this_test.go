package except_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/file/except"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func Test_This(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		except      string
		want        []string
	}{
		{
			desc:        "file list must be empty",
			ruleBuilder: file.One("foobar"),
			except:      "foobar",
			want:        nil,
		},
		{
			desc:        "file list must contain 'foobar'",
			ruleBuilder: file.One("foobar"),
			except:      "bazquux",
			want:        []string{"foobar"},
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			th := except.This(tC.except)

			th.Evaluate(tC.ruleBuilder)

			got := tC.ruleBuilder.GetFiles()

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
