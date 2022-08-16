package config

import (
	"fmt"
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"

	fe "goarkitect/internal/arch/file/except"
	"goarkitect/internal/arch/file/should"
	fs "goarkitect/internal/arch/file/should"
	ft "goarkitect/internal/arch/file/that"
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
	Shoulds []Should `yaml:"shoulds" json:"shoulds"`
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
type Should struct {
	Kind        string         `yaml:"kind" json:"kind"`
	Value       string         `yaml:"value" json:"value"`
	Suffix      string         `yaml:"suffix" json:"suffix"`
	Regex       string         `yaml:"regex" json:"regex"`
	Permissions string         `yaml:"permissions" json:"permissions"`
	Glob        string         `yaml:"glob" json:"glob"`
	Prefix      string         `yaml:"prefix" json:"prefix"`
	BasePath    string         `yaml:"basePath" json:"basePath"`
	Options     []ShouldOption `yaml:"options" json:"options"`
}
type ShouldOption struct {
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

	for _, s := range conf.Shoulds {
		opts := make([]should.Option, len(s.Options))
		for i, opt := range s.Options {
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
				return nil, []error{fmt.Errorf("unknown 'should.option' kind: '%s'", opt.Kind)}
			}
		}

		switch s.Kind {
		case "be_gitencrypted":
			rb.Should(fs.BeGitencrypted(opts...))
		case "be_gitignored":
			rb.Should(fs.BeGitignored(opts...))
		case "contain_value":
			rb.Should(fs.ContainValue([]byte(s.Value), opts...))
		case "end_with":
			rb.Should(fs.EndWith(s.Suffix, opts...))
		case "exist":
			rb.Should(fs.Exist())
		case "have_content_matching_regex":
			rb.Should(fs.HaveContentMatchingRegex(s.Regex, opts...))
		case "have_content_matching":
			rb.Should(fs.HaveContentMatching([]byte(s.Value), opts...))
		case "have_permissions":
			rb.Should(fs.HavePermissions(s.Permissions, opts...))
		case "match_glob":
			rb.Should(fs.MatchGlob(s.Glob, s.BasePath, opts...))
		case "match_regex":
			rb.Should(fs.MatchRegex(s.Regex, opts...))
		case "start_with":
			rb.Should(fs.StartWith(s.Prefix, opts...))
		default:
			return nil, []error{fmt.Errorf("unknown 'should' kind: '%s'", s.Kind)}
		}
	}

	return rb.Because(rule.Because(conf.Because))
}
