package config

import (
	"errors"
	"fmt"

	"github.com/omissis/goarkitect/internal/arch/file"
	fe "github.com/omissis/goarkitect/internal/arch/file/except"
	fs "github.com/omissis/goarkitect/internal/arch/file/expect"
	ft "github.com/omissis/goarkitect/internal/arch/file/that"
	"github.com/omissis/goarkitect/internal/arch/rule"
)

var (
	ErrUnknonwnExcept      = errors.New("unknown 'except'")
	ErrUnknonwnThat        = errors.New("unknown 'that'")
	ErrUnknownExpect       = errors.New("unknown 'expect'")
	ErrUnknownExpectOption = errors.New("unknown 'expect' option")
	ErrUnknownMatcher      = errors.New("unknown 'matcher'")
	ErrUnknownRule         = errors.New("unknown 'rule'")
)

type Root struct {
	Rules []Rule `json:"rules" yaml:"rules"`
}
type Rule struct {
	Name    string   `json:"name"    yaml:"name"`
	Kind    string   `json:"kind"    yaml:"kind"`
	Matcher Matcher  `json:"matcher" yaml:"matcher"`
	Thats   []That   `json:"thats"   yaml:"thats"`
	Excepts []Except `json:"excepts" yaml:"excepts"`
	Musts   []Expect `json:"musts"   yaml:"musts"`
	Shoulds []Expect `json:"shoulds" yaml:"shoulds"`
	Coulds  []Expect `json:"coulds"  yaml:"coulds"`
	Because string   `json:"because" yaml:"because"`
}
type Matcher struct {
	Kind      string   `json:"kind"      yaml:"kind"`
	FilePath  string   `json:"filePath"  yaml:"filePath"`
	FilePaths []string `json:"filePaths" yaml:"filePaths"`
}
type That struct {
	Kind      string `json:"kind"      yaml:"kind"`
	Folder    string `json:"folder"    yaml:"folder"`
	Recursive bool   `json:"recursive" yaml:"recursive"`
	Suffix    string `json:"suffix"    yaml:"suffix"`
	Value     string `json:"value"     yaml:"value"`
}
type Except struct {
	Kind     string `json:"kind"     yaml:"kind"`
	FilePath string `json:"filePath" yaml:"filePath"`
}
type Expect struct {
	Kind        string         `json:"kind"        yaml:"kind"`
	Value       string         `json:"value"       yaml:"value"`
	Suffix      string         `json:"suffix"      yaml:"suffix"`
	Regex       string         `json:"regex"       yaml:"regex"`
	Permissions string         `json:"permissions" yaml:"permissions"`
	Glob        string         `json:"glob"        yaml:"glob"`
	Prefix      string         `json:"prefix"      yaml:"prefix"`
	BasePath    string         `json:"basePath"    yaml:"basePath"`
	Options     []ExpectOption `json:"options"     yaml:"options"`
}
type ExpectOption struct {
	Kind      string `json:"kind"      yaml:"kind"`
	Separator string `json:"separator" yaml:"separator"`
}

type ExpectFunc func(e rule.Expect) rule.Builder

type RuleExecutionResult struct {
	RuleName   string
	Violations []rule.Violation
	Errors     []error
}

// Execute will take the config struct and execute the rules described in it against the given subjects.
func Execute(conf Root) []RuleExecutionResult {
	rers := make([]RuleExecutionResult, len(conf.Rules))

	for i, r := range conf.Rules {
		var (
			vs   []rule.Violation
			errs []error
		)

		switch r.Kind {
		case "file":
			vs, errs = ExecuteFileRule(r)

		default:
			vs = nil
			errs = []error{fmt.Errorf("'%s': %w", r.Kind, ErrUnknownRule)}
		}

		rers[i] = RuleExecutionResult{
			RuleName:   r.Name,
			Violations: vs,
			Errors:     errs,
		}
	}

	return rers
}

func ExecuteFileRule(conf Rule) ([]rule.Violation, []error) {
	rb, err := createFileBuilder(conf)
	if err != nil {
		return nil, []error{err}
	}

	if err := applyThats(rb, conf.Thats); err != nil {
		return nil, []error{err}
	}

	if err := applyExcepts(rb, conf.Excepts); err != nil {
		return nil, []error{err}
	}

	errs := []error{}

	applyExpects(conf.Musts, rb.Must, &errs)
	applyExpects(conf.Shoulds, rb.Should, &errs)
	applyExpects(conf.Coulds, rb.Could, &errs)

	if len(errs) > 0 {
		return nil, errs
	}

	return rb.Because(rule.Because(conf.Because))
}

func createFileBuilder(conf Rule) (*file.RuleBuilder, error) {
	var rb *file.RuleBuilder

	switch conf.Matcher.Kind {
	case "one":
		rb = file.One(conf.Matcher.FilePath)

	case "set":
		rb = file.Set(conf.Matcher.FilePaths...)

	case "all":
		rb = file.All()

	default:
		return nil, fmt.Errorf("'%s': %w", conf.Matcher.Kind, ErrUnknownMatcher)
	}

	return rb, nil
}

func applyThats(rb *file.RuleBuilder, ts []That) error {
	for _, t := range ts {
		switch t.Kind {
		case "are_in_folder":
			rb.That(ft.AreInFolder(t.Folder, t.Recursive))

		case "end_with":
			rb.That(ft.EndWith(t.Suffix))

		case "contain_value":
			rb.That(ft.ContainValue(t.Value))

		default:
			return fmt.Errorf("'%s': %w", t.Kind, ErrUnknonwnThat)
		}
	}

	return nil
}

func applyExcepts(rb *file.RuleBuilder, es []Except) error {
	for _, e := range es {
		switch e.Kind {
		case "this":
			rb.Except(fe.This(e.FilePath))

		default:
			return fmt.Errorf("'%s': %w", e.Kind, ErrUnknonwnExcept)
		}
	}

	return nil
}

func applyExpects(es []Expect, fn ExpectFunc, errs *[]error) {
	for _, e := range es {
		opts, err := getOpts(e)
		if err != nil {
			*errs = append(*errs, err)
		}

		if err := applyExpect(fn, e, opts); err != nil {
			*errs = append(*errs, err)
		}
	}
}

//nolint:cyclop // this will stay here until we have a better solution
func applyExpect(fn ExpectFunc, expect Expect, opts []fs.Option) error {
	switch expect.Kind {
	case "be_gitencrypted":
		fn(fs.BeGitencrypted(opts...))

	case "be_gitignored":
		fn(fs.BeGitignored(opts...))

	case "contain_value":
		fn(fs.ContainValue([]byte(expect.Value), opts...))

	case "end_with":
		fn(fs.EndWith(expect.Suffix, opts...))

	case "exist":
		fn(fs.Exist())

	case "have_content_matching_regex":
		fn(fs.HaveContentMatchingRegex(expect.Regex, opts...))

	case "have_content_matching":
		fn(fs.HaveContentMatching([]byte(expect.Value), opts...))

	case "have_permissions":
		fn(fs.HavePermissions(expect.Permissions, opts...))

	case "match_glob":
		fn(fs.MatchGlob(expect.Glob, expect.BasePath, opts...))

	case "match_regex":
		fn(fs.MatchRegex(expect.Regex, opts...))

	case "start_with":
		fn(fs.StartWith(expect.Prefix, opts...))

	default:
		return fmt.Errorf("'%s': %w", expect.Kind, ErrUnknownExpect)
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
			return nil, fmt.Errorf("'%s': %w", opt.Kind, ErrUnknownExpectOption)
		}
	}

	return opts, nil
}
