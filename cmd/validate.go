package cmd

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"goarkitect/internal/schema/santhosh"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/cli"
	"github.com/santhosh-tekuri/jsonschema"
)

func ValidateFactory() (cli.Command, error) {
	return &validateCommand{}, nil
}

type validateCommand struct {
	configFiles configFiles
	output      string
}

func (vc *validateCommand) Help() string {
	return "TBD"
}

func (vc *validateCommand) Run(args []string) int {
	basePath := getWd()

	vc.parseFlags()

	schema := vc.loadSchema(basePath)

	for _, configFile := range vc.configFiles {
		fmt.Printf("CONFIG FILE %s\n", configFile)

		conf := loadConfig[any](configFile)

		if err := schema.ValidateInterface(conf); err != nil {
			vc.logValidationError(err, conf)

			log.Fatal(err)
		} else {
			fmt.Println("ok")
		}
	}

	return 0
}

func (vc *validateCommand) Synopsis() string {
	return "Validate the configuration file(s)"
}

// parseFlags returns the list of config files, the output format and the base path
func (vc *validateCommand) parseFlags() {
	out := ""

	flagSet := flag.NewFlagSet("validate", flag.ExitOnError)

	flagSet.StringVar(&out, "output", "text", "format of the output")

	if err := flagSet.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	cfs := flagSet.Args()
	if len(cfs) < 1 {
		cfs = []string{".goarkitect.yaml"}
	}

	vc.output = out
	vc.configFiles = listConfigFiles(cfs)
}

func (vc *validateCommand) loadSchema(basePath string) *jsonschema.Schema {
	schemaPath := filepath.Join(basePath, "api/config_schema.json")

	data, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(schemaPath, bytes.NewReader(data)); err != nil {
		log.Fatal(err)
	}

	schema, err := compiler.Compile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}

	return schema
}

func (vc *validateCommand) logValidationError(err error, conf any) {
	ptrPaths := santhosh.GetPtrPaths(err.(*jsonschema.ValidationError))
	for _, path := range ptrPaths {
		value, err := json.Marshal(santhosh.GetValueAtPath(conf, path))
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(
			"path '%s' contains an invalid configuration value: %+v\n",
			santhosh.JoinPtrPath(path),
			string(value),
		)
	}
}
