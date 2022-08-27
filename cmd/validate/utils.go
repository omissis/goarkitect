package validate

import (
	"errors"
	"fmt"

	"github.com/omissis/goarkitect/internal/jsonx"
	"github.com/omissis/goarkitect/internal/logx"
	"github.com/omissis/goarkitect/internal/schema/santhosh"
)

var ErrHasValidationErrors = errors.New("schema has validation errors")

func PrintResults(output string, err error, conf any, configFile string) {
	ptrPaths := santhosh.GetPtrPaths(err)

	switch output {
	case "text":
		// TODO: improve formatting
		fmt.Printf("CONFIG FILE %s\n", configFile)

		for _, path := range ptrPaths {
			value, serr := santhosh.GetValueAtPath(conf, path)
			if serr != nil {
				logx.Fatal(serr)
			}

			// TODO: improve santhosh.JoinPtrPath output
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
		logx.Fatal(fmt.Errorf("unknown output format: '%s', supported ones are: json, text", output))
	}
}
