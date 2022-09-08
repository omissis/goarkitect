package config

import (
	"fmt"

	"github.com/omissis/goarkitect/internal/arch/file"
	fe "github.com/omissis/goarkitect/internal/arch/file/except"
	fs "github.com/omissis/goarkitect/internal/arch/file/expect"
	ft "github.com/omissis/goarkitect/internal/arch/file/that"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

type Root struct {
	Rules []Rule `yaml:"rules" json:"rules"`
}
type Rule struct {
	Name    string   `yaml:"name" json:"name"`
	Kind    string   `yaml:"kind" json:"kind"`
	Matcher Matcher  `yaml:"matcher" json:"matcher"`
	Thats   []That   `yaml:"thats" json:"thats"`
	Excepts []Except `yaml:"excepts" json:"excepts"`
	Musts   []Expect `yaml:"musts" json:"musts"`
	Shoulds []Expect `yaml:"shoulds" json:"shoulds"`
	Coulds  []Expect `yaml:"coulds" json:"coulds"`
	Because string   `yaml:"because" json:"because"`
}
type Matcher struct {
	Kind      string   `yaml:"kind" json:"kind"`
	FilePath  string   `yaml:"filePath" json:"filePath"`
	FilePaths []string `yaml:"filePaths" json:"filePaths"`
}
type That struct {
	Kind      string `yaml:"kind" json:"kind"`
	Folder    string `yaml:"folder" json:"folder"`
	Recursive bool   `yaml:"recursive" json:"recursive"`
	Suffix    string `yaml:"suffix" json:"suffix"`
}
type Except struct {
	Kind     string `yaml:"kind" json:"kind"`
	FilePath string `yaml:"filePath" json:"filePath"`
}
type Expect struct {
	Kind        string         `yaml:"kind" json:"kind"`
	Value       string         `yaml:"value" json:"value"`
	Suffix      string         `yaml:"suffix" json:"suffix"`
	Regex       string         `yaml:"regex" json:"regex"`
	Permissions string         `yaml:"permissions" json:"permissions"`
	Glob        string         `yaml:"glob" json:"glob"`
	Prefix      string         `yaml:"prefix" json:"prefix"`
	BasePath    string         `yaml:"basePath" json:"basePath"`
	Options     []ExpectOption `yaml:"options" json:"options"`
}
type ExpectOption struct {
	Kind      string `yaml:"kind" json:"kind"`
	Separator string `yaml:"separator" json:"separator"`
}

type RuleExecutionResult struct {
	RuleName   string
	Violations []rule.Violation
	Errors     []error
}

// Execute will take the config struct and execute the rules described in it against the given subjects.
func Execute(conf Root) []RuleExecutionResult {
	var rers []RuleExecutionResult

	for _, r := range conf.Rules {
		var vs []rule.Violation
		var errs []error

		switch r.Kind {
		case "file":
			vs, errs = ExecuteFileRule(r)
		default:
			vs = nil
			errs = []error{fmt.Errorf("unknown 'rule' kind: '%s'", r.Kind)}
		}

		rers = append(rers, RuleExecutionResult{
			RuleName:   r.Name,
			Violations: vs,
			Errors:     errs,
		})
	}

	return rers
}

func ExecuteFileRule(conf Rule) ([]rule.Violation, []error) {
	var rb *file.RuleBuilder

	switch conf.Matcher.Kind {
	case "one":
		rb = file.One(conf.Matcher.FilePath)
	case "set":
		rb = file.Set(conf.Matcher.FilePaths...)
	case "all":
		rb = file.All()
	default:
		return nil, []error{fmt.Errorf("unknown 'matcher' kind: '%s'", conf.Matcher.Kind)}
	}

	for _, t := range conf.Thats {
		switch t.Kind {
		case "are_in_folder":
			rb.That(ft.AreInFolder(t.Folder, t.Recursive))
		case "end_with":
			rb.That(ft.EndWith(t.Suffix))
		default:
			return nil, []error{fmt.Errorf("unknown 'that' kind: '%s'", t.Kind)}
		}
	}

	for _, e := range conf.Excepts {
		switch e.Kind {
		case "this":
			rb.Except(fe.This(e.FilePath))
		default:
			return nil, []error{fmt.Errorf("unknown 'except' kind: '%s'", e.Kind)}
		}
	}

	errs := []error{}

	for _, m := range conf.Musts {
		opts, err := getOpts(m)
		if err != nil {
			errs = append(errs, err)
		}

		if err := applyExpects(rb.Must, m, opts); err != nil {
			errs = append(errs, err)
		}
	}

	for _, s := range conf.Shoulds {
		opts, err := getOpts(s)
		if err != nil {
			errs = append(errs, err)
		}

		if err := applyExpects(rb.Should, s, opts); err != nil {
			errs = append(errs, err)
		}
	}

	for _, c := range conf.Coulds {
		opts, err := getOpts(c)
		if err != nil {
			errs = append(errs, err)
		}

		if err := applyExpects(rb.Could, c, opts); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return rb.Because(rule.Because(conf.Because))
}

func applyExpects(expectFn func(e rule.Expect) rule.Builder, expect Expect, opts []fs.Option) error {
	switch expect.Kind {
	case "be_gitencrypted":
		expectFn(fs.BeGitencrypted(opts...))
	case "be_gitignored":
		expectFn(fs.BeGitignored(opts...))
	case "contain_value":
		expectFn(fs.ContainValue([]byte(expect.Value), opts...))
	case "end_with":
		expectFn(fs.EndWith(expect.Suffix, opts...))
	case "exist":
		expectFn(fs.Exist())
	case "have_content_matching_regex":
		expectFn(fs.HaveContentMatchingRegex(expect.Regex, opts...))
	case "have_content_matching":
		expectFn(fs.HaveContentMatching([]byte(expect.Value), opts...))
	case "have_permissions":
		expectFn(fs.HavePermissions(expect.Permissions, opts...))
	case "match_glob":
		expectFn(fs.MatchGlob(expect.Glob, expect.BasePath, opts...))
	case "match_regex":
		expectFn(fs.MatchRegex(expect.Regex, opts...))
	case "start_with":
		expectFn(fs.StartWith(expect.Prefix, opts...))
	default:
		return fmt.Errorf("unknown 'should' kind: '%s'", expect.Kind)
	}

	return nil
}

func getOpts(e Expect) ([]fs.Option, error) {
	opts := make([]fs.Option, len(e.Options))
	for i, opt := range e.Options {
		switch opt.Kind {
		case "negated":
			opts[i] = fs.Negated{}
		case "ignore_case":
			opts[i] = fs.IgnoreCase{}
		case "ignore_new_lines_at_the_end_of_file":
			opts[i] = fs.IgnoreNewLinesAtTheEndOfFile{}
		case "match_single_lines":
			opts[i] = fs.MatchSingleLines{
				Separator: opt.Separator,
			}
		default:
			return nil, fmt.Errorf("unknown 'should.option' kind: '%s'", opt.Kind)
		}
	}

	return opts, nil
}
