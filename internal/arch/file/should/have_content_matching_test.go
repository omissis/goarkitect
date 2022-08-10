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

func Test_HaveContentMatching(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		content     string
		options     []should.ContentMatchOption
		want        []rule.Violation
	}{
		{
			desc:        "content of file 'foobar' matches expected content",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "foo bar baz quux\n",
			want:        nil,
		},
		{
			desc:        "content of file 'foobar' matches expected content, ignoring newlines at the end of file",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "foo bar baz quux",
			options: []should.ContentMatchOption{
				should.IgnoreNewLinesAtTheEndOfFile{},
			},
			want: nil,
		},
		{
			desc:        "content of file 'foobar' matches expected content, ignoring case",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "FOO BAR BAZ QUUX\n",
			options: []should.ContentMatchOption{
				should.IgnoreCase{},
			},
			want: nil,
		},
		{
			desc:        "content of file 'foobar' does not match expected content",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "something else",
			want: []rule.Violation{
				rule.NewViolation("file 'foobar.txt' does not have content matching 'something else'"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			hcm := should.HaveContentMatching([]byte(tC.content), tC.options...)
			got := hcm.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
