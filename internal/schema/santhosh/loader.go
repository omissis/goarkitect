package santhosh

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema"
)

func LoadSchema(basePath string) (*jsonschema.Schema, error) {
	schemaPath := filepath.Join(basePath, "api/config_schema.json")

	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(schemaPath, bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return compiler.Compile(schemaPath)
}
