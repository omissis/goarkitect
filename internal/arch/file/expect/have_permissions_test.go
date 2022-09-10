package expect_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/file/expect"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

func Test_HavePermissions(t *testing.T) {
	t.Parallel()

	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		permissions string
		options     []expect.Option
		want        []rule.CoreViolation
		wantErrs    []error
	}{
		{
			desc:        "wrong permissions string",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions/0700.txt")),
			permissions: "foobarbaz-",
			want:        nil,
			wantErrs:    []error{expect.ErrInvalidPermissions},
		},
		{
			desc:        "permissions of directory 'test/permissions' match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions")),
			permissions: "drwxr-xr-x",
			want:        nil,
		},
		{
			desc:        "permissions of file '0700.txt' match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions/0700.txt")),
			permissions: "-rwx------",
			want:        nil,
		},
		{
			desc:        "permissions of directory 'test/permissions' do not match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions")),
			permissions: "dr--r--r--",
			want: []rule.CoreViolation{
				rule.NewCoreViolation("directory 'permissions' does not have permissions matching 'dr--r--r--'"),
			},
		},
		{
			desc:        "permissions of file '0700.txt' do not match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions/0700.txt")),
			permissions: "-rwxrwxrwx",
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file '0700.txt' does not have permissions matching '-rwxrwxrwx'"),
			},
		},
		{
			desc:        "negated: permissions of directory 'test/permissions' match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions")),
			permissions: "drwxr-xr-x",
			options:     []expect.Option{expect.Negated{}},
			want: []rule.CoreViolation{
				rule.NewCoreViolation("directory 'permissions' does have permissions matching 'drwxr-xr-x'"),
			},
		},
		{
			desc:        "negated: permissions of file '0700.txt' match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions/0700.txt")),
			permissions: "-rwx------",
			options:     []expect.Option{expect.Negated{}},
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file '0700.txt' does have permissions matching '-rwx------'"),
			},
		},
		{
			desc:        "negated: permissions of directory 'test/permissions' do not match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions")),
			permissions: "dr--r--r--",
			options:     []expect.Option{expect.Negated{}},
			want:        nil,
		},
		{
			desc:        "negated: permissions of file '0700.txt' do not match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions/0700.txt")),
			permissions: "-rwxrwxrwx",
			options:     []expect.Option{expect.Negated{}},
			want:        nil,
		},
	}

	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			hcm := expect.HavePermissions(tC.permissions, tC.options...)
			got := hcm.Evaluate(tC.ruleBuilder)
			gotErrs := hcm.GetErrors()

			if !cmp.Equal(gotErrs, tC.wantErrs, cmpopts.EquateErrors(), cmpopts.EquateEmpty()) {
				t.Errorf("wantErr = %+v, gotErr = %+v", tC.wantErrs, gotErrs)
			}

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.CoreViolation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
