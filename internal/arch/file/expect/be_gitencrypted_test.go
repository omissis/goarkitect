package expect_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/file/expect"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func Test_BeGitencrypted(t *testing.T) {
	if _, err := exec.LookPath("git-crypt"); err != nil {
		t.Skip("git-crypt is not installed")
	}

	t.Parallel()

	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		options     []expect.Option
		want        []rule.CoreViolation
	}{
		{
			desc:        "file 'encrypted.txt' is expected to be gitencrypted",
			ruleBuilder: file.One(filepath.Join(basePath, "test/encrypted.txt")),
			want:        nil,
		},
		{
			desc:        "file 'not_encrypted.txt' is expected to be gitencrypted",
			ruleBuilder: file.One(filepath.Join(basePath, "test/not_encrypted.txt")),
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file 'not_encrypted.txt' is not gitencrypted"),
			},
		},
		{
			desc:        "negated: file 'encrypted.txt' is not expected to be gitencrypted",
			ruleBuilder: file.One(filepath.Join(basePath, "test/encrypted.txt")),
			options:     []expect.Option{expect.Negated{}},
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file 'encrypted.txt' is gitencrypted"),
			},
		},
		{
			desc:        "negated: file 'not_encrypted.txt' is not expected to be gitencrypted",
			ruleBuilder: file.One(filepath.Join(basePath, "test/not_encrypted.txt")),
			options:     []expect.Option{expect.Negated{}},
			want:        nil,
		},
	}

	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			hcm := expect.BeGitencrypted(tC.options...)
			got := hcm.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.CoreViolation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
