package cmd

import (
	"bytes"
	"encoding/json"
	"goarkitect/internal/schema/santhosh"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/cli"
	"github.com/santhosh-tekuri/jsonschema"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

func ValidateFactory() (cli.Command, error) {
	return &validateCommand{}, nil
}

type validateCommand struct {
}

func (vc *validateCommand) Help() string {
	return "TBD"
}

func (vc *validateCommand) Run(args []string) int {
	basePath := vc.getCwd()

	ruleset := vc.parseArgs(basePath, args)

	schema := vc.loadSchema(basePath)

	for _, subject := range ruleset {
		// TODO: print in debug mode
		// log.Printf("validating %s", subject)

		configPath := filepath.Join(basePath, subject)

		configData, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatal(err)
		}

		var conf interface{}
		if err := yaml.Unmarshal(configData, &conf); err != nil {
			log.Fatal(err)
		}

		if err := schema.ValidateInterface(conf); err != nil {
			vc.logValidationError(err, subject, conf)

			log.Fatal(err)
		}
	}

	return 0
}

func (vc *validateCommand) Synopsis() string {
	return "Validate the configuration file(s)"
}

func (vc *validateCommand) getCwd() string {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return basePath
}

func (vc *validateCommand) parseArgs(basePath string, args []string) []string {
	if len(args) < 1 {
		return []string{".goarkitect.yaml"}
	}

	ruleset := make([]string, 0)

	for _, arg := range args {
		fileInfo, err := os.Stat(arg)
		if err != nil {
			log.Fatal(err)
		}

		if !fileInfo.IsDir() {
			ruleset = append(ruleset, arg)
			continue
		}

		if err := filepath.Walk(arg, vc.visitFolder(&ruleset)); err != nil {
			log.Fatal(err)
		}
	}

	slices.Sort(ruleset)

	return slices.Compact(ruleset)
}

func (e *validateCommand) visitFolder(files *[]string) filepath.WalkFunc {
	return func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !file.IsDir() {
			*files = append(*files, path)
		}

		return nil
	}
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

func (vc *validateCommand) logValidationError(err error, subject string, conf any) {
	ptrPaths := santhosh.GetPtrPaths(err.(*jsonschema.ValidationError))
	for _, path := range ptrPaths {
		value, err := json.Marshal(santhosh.GetValueAtPath(conf, path))
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(
			"file '%s': path '%s' contains an invalid configuration value: %+v\n",
			subject,
			santhosh.JoinPtrPath(path),
			string(value),
		)
	}
}
