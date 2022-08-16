package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"goarkitect/internal/arch/rule"
	"goarkitect/internal/config"

	"github.com/mitchellh/cli"
)

func VerifyFactory() (cli.Command, error) {
	return &verifyCommand{}, nil
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

	for _, configFile := range vc.configFiles {
		fmt.Printf("CONFIG FILE %s\n", configFile)

		conf := loadConfig[config.Root](configFile)

		results := config.Execute(conf)

		vc.printResults(results)

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
	out := ""

	flagSet := flag.NewFlagSet("verify", flag.ExitOnError)

	flagSet.Var(&cfs, "config", "path to the config file or folder")
	flagSet.StringVar(&out, "output", "text", "format of the output")

	if err := flagSet.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	if len(cfs) < 1 {
		cfs = []string{filepath.Join(getWd(), ".goarkitect.yaml")}
	}

	vc.output = out
	vc.configFiles = listConfigFiles(cfs)
}

func (vc *verifyCommand) printResults(results []config.RuleExecutionResult) {
	for _, r := range results {
		fmt.Printf("\nRULE '%s'\n", r.RuleName)

		fmt.Println("Violations:")
		for _, v := range r.Violations {
			fmt.Printf("- %s\n", v)
		}
		if len(r.Violations) == 0 {
			fmt.Println("- None")
		}

		fmt.Println("Errors:")
		for _, v := range r.Errors {
			fmt.Printf("- %s\n", v)
		}
		if len(r.Errors) == 0 {
			fmt.Println("- None")
		}
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
