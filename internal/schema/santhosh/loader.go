package santhosh

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema"
)

func LoadSchema(basePath string) (*jsonschema.Schema, error) {
	schemaPath := filepath.Join(basePath, "api/config_schema.json")

	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	if err = compiler.AddResource(schemaPath, bytes.NewReader(data)); err != nil {
		return nil, fmt.Errorf("failed to add resource to json schema compiler: %w", err)
	}

	schema, err := compiler.Compile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to compile json schema: %w", err)
	}

	return schema, nil
}
