package verify

import (
	"errors"
	"fmt"

	"github.com/omissis/goarkitect/internal/arch/rule"
	"github.com/omissis/goarkitect/internal/config"
	"github.com/omissis/goarkitect/internal/jsonx"
	"github.com/omissis/goarkitect/internal/logx"
)

var ErrProjectDoesNotRespectRules = errors.New("project does not respect defined rules")

func PrintResults(output string, configFile string, results []config.RuleExecutionResult) {
	switch output {
	case "text":
		// TODO: improve formatting
		fmt.Printf("CONFIG FILE %s\n", configFile)

		for _, r := range results {
			fmt.Printf("\nRULE '%s'\n", r.RuleName)

			fmt.Printf("Violations:\n")
			for _, v := range r.Violations {
				fmt.Printf("- %s\n", v)
			}
			if len(r.Violations) == 0 {
				fmt.Printf("- None\n")
			}

			fmt.Printf("Errors:\n")
			for _, v := range r.Errors {
				fmt.Printf("- %s\n", v)
			}
			if len(r.Errors) == 0 {
				fmt.Printf("- None\n")
			}
		}
	case "json":
		fmt.Println(
			jsonx.Marshal(
				map[string]any{
					"configFile": configFile,
					"results":    results,
				},
			),
		)
	default:
		logx.Fatal(fmt.Errorf("unknown output format: '%s', supported ones are: json, text", output))
	}
}

func HasErrors(results []config.RuleExecutionResult) bool {
	for _, r := range results {
		for _, v := range r.Violations {
			if v.Severity() == rule.Error.String() {
				return true
			}
		}
	}

	return false
}
