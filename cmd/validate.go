package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/omissis/goarkitect/internal/jsonx"
	"github.com/omissis/goarkitect/internal/logx"
	"github.com/omissis/goarkitect/internal/schema/santhosh"

	"github.com/mitchellh/cli"
)

func ValidateFactory(output string) (cli.Command, error) {
	return &validateCommand{
		output: output,
	}, nil
}

type validateCommand struct {
	configFiles configFiles
	output      string
}

func (vc *validateCommand) Help() string {
	return "TBD"
}

func (vc *validateCommand) Run(args []string) int {
	exitCode := 0
	basePath := getWd()

	vc.parseFlags()

	if len(vc.configFiles) == 0 {
		logx.Fatal(errors.New("no config files found"))
	}

	schema, err := santhosh.LoadSchema(basePath)
	if err != nil {
		logx.Fatal(err)
	}

	for _, configFile := range vc.configFiles {
		conf := loadConfig[any](configFile)

		if err := schema.ValidateInterface(conf); err != nil {
			vc.printResults(err, conf, configFile)

			exitCode = 1
		}
	}

	return exitCode
}

func (vc *validateCommand) Synopsis() string {
	return "Validate the configuration file(s)"
}

func (vc *validateCommand) parseFlags() {
	flagSet := flag.NewFlagSet("validate", flag.ContinueOnError)

	if err := flagSet.Parse(os.Args[2:]); err != nil {
		logx.Fatal(err)
	}

	cfs := flagSet.Args()
	if len(cfs) < 1 {
		cfs = []string{".goarkitect.yaml"}
	}

	vc.configFiles = listConfigFiles(cfs)
}

func (vc *validateCommand) printResults(err error, conf any, configFile string) {
	ptrPaths := santhosh.GetPtrPaths(err)

	switch vc.output {
	case "text":
		// TODO: improve formatting
		fmt.Printf("CONFIG FILE %s\n", configFile)

		for _, path := range ptrPaths {
			value, serr := santhosh.GetValueAtPath(conf, path)
			if serr != nil {
				logx.Fatal(serr)
			}

			// TODO: improve this output
			fmt.Printf(
				"path '%s' contains an invalid configuration value: %+v\n",
				santhosh.JoinPtrPath(path),
				value,
			)
		}

		fmt.Println(err)
	case "json":
		for _, path := range ptrPaths {
			value, serr := santhosh.GetValueAtPath(conf, path)
			if serr != nil {
				logx.Fatal(serr)
			}

			fmt.Println(
				jsonx.Marshal(
					map[string]any{
						"file":    configFile,
						"message": "path contains an invalid configuration value",
						"path":    santhosh.JoinPtrPath(path),
						"value":   value,
					},
				),
			)
		}

		fmt.Println(jsonx.Marshal(err))
	default:
		logx.Fatal(fmt.Errorf("unknown output format: '%s'", vc.output))
	}
}
