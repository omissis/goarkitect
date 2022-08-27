package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/omissis/goarkitect/internal/cli"
	"github.com/omissis/goarkitect/internal/jsonx"
	"github.com/omissis/goarkitect/internal/logx"
	"github.com/omissis/goarkitect/internal/schema/santhosh"
)

var ErrHasValidationErrors = errors.New("schema has validation errors")

func NewValidateCommand(output *string) cli.Command {
	return &validateCommand{
		output: output,
	}
}

type validateCommand struct {
	configFiles configFiles
	output      *string
}

func (vc *validateCommand) Name() string {
	return "validate"
}

func (vc *validateCommand) Help() string {
	return "TBD"
}

func (vc *validateCommand) Run(args []string) error {
	basePath := getWd()

	vc.parseFlags()

	if len(vc.configFiles) == 0 {
		return errors.New("no config files found")
	}

	schema, err := santhosh.LoadSchema(basePath)
	if err != nil {
		return err
	}

	hasErrors := error(nil)
	for _, configFile := range vc.configFiles {
		conf := loadConfig[any](configFile)

		if err := schema.ValidateInterface(conf); err != nil {
			vc.printResults(err, conf, configFile)

			hasErrors = ErrHasValidationErrors
		}
	}

	return hasErrors
}

func (vc *validateCommand) Synopsis() string {
	return "Validate the configuration file(s)"
}

func (vc *validateCommand) parseFlags() {
	flagSet := flag.NewFlagSet("validate", flag.ContinueOnError)

	if err := flagSet.Parse(cli.GetArgs(os.Args, 2)); err != nil {
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

	switch *vc.output {
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
