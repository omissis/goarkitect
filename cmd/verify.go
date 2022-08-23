package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/omissis/goarkitect/internal/arch/rule"
	"github.com/omissis/goarkitect/internal/config"
	"github.com/omissis/goarkitect/internal/jsonx"
	"github.com/omissis/goarkitect/internal/logx"

	"github.com/mitchellh/cli"
)

func VerifyFactory(output string) (cli.Command, error) {
	return &validateCommand{
		output: output,
	}, nil
}

type verifyCommand struct {
	configFiles configFiles
	output      string
}

func (vc *verifyCommand) Help() string {
	return "TBD"
}

func (vc *verifyCommand) Run(args []string) int {
	exitCode := 0

	vc.parseFlags()

	if len(vc.configFiles) == 0 {
		logx.Fatal(errors.New("no config files found"))
	}

	for _, configFile := range vc.configFiles {
		conf := loadConfig[config.Root](configFile)

		results := config.Execute(conf)

		vc.printResults(configFile, results)

		if vc.hasErrors(results) {
			exitCode = 1
		}
	}

	return exitCode
}

func (vc *verifyCommand) Synopsis() string {
	return "Verify the ruleset against a project"
}

func (vc *verifyCommand) parseFlags() {
	cfs := configFiles{}

	flagSet := flag.NewFlagSet("verify", flag.ContinueOnError)

	flagSet.Var(&cfs, "config", "path to the config file or folder")

	if err := flagSet.Parse(os.Args[2:]); err != nil {
		logx.Fatal(err)
	}

	if len(cfs) < 1 {
		cfs = []string{filepath.Join(getWd(), ".goarkitect.yaml")}
	}

	vc.configFiles = listConfigFiles(cfs)
}

func (vc *verifyCommand) printResults(configFile string, results []config.RuleExecutionResult) {
	switch vc.output {
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
		logx.Fatal(fmt.Errorf("unknown output format: '%s'", vc.output))
	}
}

func (vc *verifyCommand) hasErrors(results []config.RuleExecutionResult) bool {
	for _, r := range results {
		for _, v := range r.Violations {
			if v.Severity() == rule.Error.String() {
				return true
			}
		}
	}

	return false
}
