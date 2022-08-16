package cmd

import (
	"flag"
	"fmt"
	"goarkitect/internal/config"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

// configs contains a list of files where ruleset are specified
type configFiles []string

func (i *configFiles) String() string {
	return strings.Join(*i, ",")
}

func (i *configFiles) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func VerifyFactory() (cli.Command, error) {
	return &verifyCommand{}, nil
}

type verifyCommand struct {
	configFiles configFiles
	output      string
}

func (vc *verifyCommand) Help() string {
	return "Usage: goarkitect verify [options] [ruleset(s)]"
}

func (vc *verifyCommand) Run(args []string) int {
	vc.parseFlags()

	for _, configFile := range vc.configFiles {
		fmt.Printf("CONFIG FILE %s\n", configFile)
		// TODO: recognize if config is relative or absolute, then adjust configFile accordingly
		configData, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}

		var conf config.Root
		if err := yaml.Unmarshal(configData, &conf); err != nil {
			log.Fatal(err)
		}

		results := config.Execute(conf)

		vc.printResults(results)
	}

	return 0
}

func (vc *verifyCommand) Synopsis() string {
	return "Verify the ruleset against a project"
}

func (vc *verifyCommand) getCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return cwd
}

// parseFlags returns the list of config files, the output format and the base path
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
		cfs = []string{filepath.Join(vc.getCwd(), ".goarkitect.yaml")}
	}

	vc.output = out
	vc.configFiles = vc.listConfigFiles(cfs)
}

func (vc *verifyCommand) listConfigFiles(cfs []string) []string {
	configFiles := make([]string, 0)

	for _, cf := range cfs {
		fileInfo, err := os.Stat(cf)
		if err != nil {
			log.Fatal(err)
		}

		if !fileInfo.IsDir() {
			configFiles = append(configFiles, cf)
			continue
		}

		if err := filepath.Walk(cf, vc.visitConfigFolder(&configFiles)); err != nil {
			log.Fatal(err)
		}
	}

	slices.Sort(configFiles)

	return slices.Compact(configFiles)
}

func (vc *verifyCommand) visitConfigFolder(files *[]string) filepath.WalkFunc {
	return func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !file.IsDir() && filepath.Ext(path) == ".yaml" {
			*files = append(*files, path)
		}

		return nil
	}
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
