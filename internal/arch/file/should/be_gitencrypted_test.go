package should_test

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/file/should"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_BeGitencrypted(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		options     []should.Option
		want        []rule.Violation
	}{
		{
			desc:        "file 'encrypted.txt' should be gitencrypted",
			ruleBuilder: file.One(filepath.Join(basePath, "test/encrypted.txt")),
			want:        nil,
		},
		{
			desc:        "file 'not_encrypted.txt' should be gitencrypted",
			ruleBuilder: file.One(filepath.Join(basePath, "test/not_encrypted.txt")),
			want: []rule.Violation{
				rule.NewViolation("file 'not_encrypted.txt' is not gitencrypted"),
			},
		},
		{
			desc:        "negated: file 'encrypted.txt' should not be gitencrypted",
			ruleBuilder: file.One(filepath.Join(basePath, "test/encrypted.txt")),
			options: []should.Option{
				should.Negated{},
			},
			want: []rule.Violation{
				rule.NewViolation("file 'encrypted.txt' is gitencrypted"),
			},
		},
		{
			desc:        "negated: file 'not_encrypted.txt' should not be gitencrypted",
			ruleBuilder: file.One(filepath.Join(basePath, "test/not_encrypted.txt")),
			options: []should.Option{
				should.Negated{},
			},
			want: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			hcm := should.BeGitencrypted(tC.options...)
			got := hcm.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
