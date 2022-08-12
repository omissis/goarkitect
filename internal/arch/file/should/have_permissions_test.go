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

func Test_HavePermissions(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		permissions string
		options     []should.Option
		want        []rule.Violation
		wantErrs    []error
	}{
		{
			desc:        "wrong permissions string",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions/0700.txt")),
			permissions: "foobarbaz-",
			want:        nil,
			wantErrs:    []error{should.ErrInvalidPermissions},
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
			want: []rule.Violation{
				rule.NewViolation("directory 'permissions' does not have permissions matching 'dr--r--r--'"),
			},
		},
		{
			desc:        "permissions of file '0700.txt' do not match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions/0700.txt")),
			permissions: "-rwxrwxrwx",
			want: []rule.Violation{
				rule.NewViolation("file '0700.txt' does not have permissions matching '-rwxrwxrwx'"),
			},
		},
		{
			desc:        "negated: permissions of directory 'test/permissions' match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions")),
			permissions: "drwxr-xr-x",
			options:     []should.Option{should.Negated{}},
			want: []rule.Violation{
				rule.NewViolation("directory 'permissions' does have permissions matching 'drwxr-xr-x'"),
			},
		},
		{
			desc:        "negated: permissions of file '0700.txt' match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions/0700.txt")),
			permissions: "-rwx------",
			options:     []should.Option{should.Negated{}},
			want: []rule.Violation{
				rule.NewViolation("file '0700.txt' does have permissions matching '-rwx------'"),
			},
		},
		{
			desc:        "negated: permissions of directory 'test/permissions' do not match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions")),
			permissions: "dr--r--r--",
			options:     []should.Option{should.Negated{}},
			want:        nil,
		},
		{
			desc:        "negated: permissions of file '0700.txt' do not match expected one",
			ruleBuilder: file.One(filepath.Join(basePath, "test/permissions/0700.txt")),
			permissions: "-rwxrwxrwx",
			options:     []should.Option{should.Negated{}},
			want:        nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			hcm := should.HavePermissions(tC.permissions, tC.options...)
			got := hcm.Evaluate(tC.ruleBuilder)
			gotErrs := hcm.GetErrors()

			if !cmp.Equal(gotErrs, tC.wantErrs, cmpopts.EquateErrors(), cmpopts.EquateEmpty()) {
				t.Errorf("wantErr = %+v, gotErr = %+v", tC.wantErrs, gotErrs)
			}

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
