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
		options     []should.Option
		want        []rule.Violation
	}{
		{
			desc:        "content of file 'foobar.txt' matches expected content",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "foo bar baz quux\n",
			want:        nil,
		},
		{
			desc:        "content of file 'foobar.txt' matches expected content, ignoring newlines at the end of file",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "foo bar baz quux",
			options: []should.Option{
				should.IgnoreNewLinesAtTheEndOfFile{},
			},
			want: nil,
		},
		{
			desc:        "content of file 'foobar.txt' matches expected content, ignoring case",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "FOO BAR BAZ QUUX\n",
			options: []should.Option{
				should.IgnoreCase{},
			},
			want: nil,
		},
		{
			desc:        "content of file 'foobar.txt' does not match expected content",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "something else",
			want: []rule.Violation{
				rule.NewViolation("file 'foobar.txt' does not have content matching 'something else'"),
			},
		},
		{
			desc:        "every line of file 'baz.txt' matches expected content",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "foo bar baz quux",
			options: []should.Option{
				should.IgnoreNewLinesAtTheEndOfFile{},
				should.MatchSingleLines{},
			},
			want: nil,
		},
		{
			desc:        "not every line of file 'baz.txt' matches regex",
			ruleBuilder: file.One(filepath.Join(basePath, "test/baz.txt")),
			content:     "something else",
			options: []should.Option{
				should.IgnoreNewLinesAtTheEndOfFile{},
				should.MatchSingleLines{},
			},
			want: []rule.Violation{
				rule.NewViolation("file 'baz.txt' does not have all lines matching 'something else'"),
			},
		},
		{
			desc:        "negated: content of file 'foobar.txt' does not match expected content",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "something else\n",
			options: []should.Option{
				should.Negated{},
			},
			want: nil,
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