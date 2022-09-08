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

func Test_HaveContentMatching(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc        string
		ruleBuilder *file.RuleBuilder
		content     string
		options     []expect.Option
		want        []rule.CoreViolation
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
			options: []expect.Option{
				expect.IgnoreNewLinesAtTheEndOfFile{},
			},
			want: nil,
		},
		{
			desc:        "content of file 'foobar.txt' matches expected content, ignoring case",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "FOO BAR BAZ QUUX\n",
			options: []expect.Option{
				expect.IgnoreCase{},
			},
			want: nil,
		},
		{
			desc:        "content of file 'foobar.txt' does not match expected content",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "something else",
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file 'foobar.txt' does not have content matching 'something else'"),
			},
		},
		{
			desc:        "every line of file 'baz.txt' matches expected content",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "foo bar baz quux",
			options: []expect.Option{
				expect.IgnoreNewLinesAtTheEndOfFile{},
				expect.MatchSingleLines{},
			},
			want: nil,
		},
		{
			desc:        "not every line of file 'baz.txt' matches regex",
			ruleBuilder: file.One(filepath.Join(basePath, "test/baz.txt")),
			content:     "something else",
			options: []expect.Option{
				expect.IgnoreNewLinesAtTheEndOfFile{},
				expect.MatchSingleLines{},
			},
			want: []rule.CoreViolation{
				rule.NewCoreViolation("file 'baz.txt' does not have all lines matching 'something else'"),
			},
		},
		{
			desc:        "negated: content of file 'foobar.txt' does not match expected content",
			ruleBuilder: file.One(filepath.Join(basePath, "test/foobar.txt")),
			content:     "something else\n",
			options:     []expect.Option{expect.Negated{}},
			want:        nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			hcm := expect.HaveContentMatching([]byte(tC.content), tC.options...)
			got := hcm.Evaluate(tC.ruleBuilder)

			if !cmp.Equal(got, tC.want, cmp.AllowUnexported(rule.CoreViolation{}), cmpopts.EquateEmpty()) {
				t.Errorf("want = %+v, got = %+v", tC.want, got)
			}
		})
	}
}
