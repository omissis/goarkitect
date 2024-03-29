package that_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/file/that"
)

func Test_AreInFolder(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		folder      string
		recursive   bool
		want        []string
		wantErrs    []string
	}{
		{
			desc:        "files in 'test/project' folder",
			ruleBuilder: file.All(),
			folder:      "../test/project",
			recursive:   false,
			want:        []string{"../test/project/Dockerfile", "../test/project/Makefile"},
			wantErrs:    nil,
		},
		{
			desc:        "files in non-existing folder",
			ruleBuilder: file.All(),
			folder:      "/does/not/exist",
			recursive:   false,
			want:        nil,
			wantErrs:    []string{"open /does/not/exist: no such file or directory"},
		},
		{
			desc:        "files in 'test' folder, recursively",
			ruleBuilder: file.All(),
			folder:      "./test",
			recursive:   true,
			want: []string{
				"test/project1/Dockerfile",
				"test/project1/Makefile",
				"test/project2/Dockerfile.1",
				"test/project2/Dockerfile.2",
				"test/project3/baz.txt",
				"test/project3/quux.txt",
			},
			wantErrs: nil,
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			aif := that.AreInFolder(tC.folder, tC.recursive)
			aif.Evaluate(tC.ruleBuilder)

			got := tC.ruleBuilder.GetFiles()
			if !cmp.Equal(got, tC.want, cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}

			gotErrs := tC.ruleBuilder.GetErrors()
			if len(gotErrs) != len(tC.wantErrs) {
				t.Errorf("want %d errs, got %d", len(tC.wantErrs), len(gotErrs))
			}

			for i := 0; i < len(gotErrs); i++ {
				if !strings.Contains(gotErrs[i].Error(), tC.wantErrs[i]) {
					t.Errorf("wantErr[%d] = %+v, gotErr[%d] = %+v", i, tC.wantErrs[i], i, gotErrs[i].Error())
				}
			}
		})
	}
}
