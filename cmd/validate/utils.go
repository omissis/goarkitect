package validate

import (
	"errors"
	"fmt"

	"github.com/omissis/goarkitect/cmd/cmdutil"
	"github.com/omissis/goarkitect/internal/jsonx"
	"github.com/omissis/goarkitect/internal/logx"
	"github.com/omissis/goarkitect/internal/schema/santhosh"
)

var ErrHasValidationErrors = errors.New("schema has validation errors")

func PrintSummary(output string, hasErrors bool) {
	switch output {
	case "text":
		if hasErrors {
			fmt.Println("Validation failed")
		} else {
			fmt.Println("Validation succeeded")
		}

	case "json":
		if hasErrors {
			fmt.Println("{\"result\":\"Validation failed\"}")
		} else {
			fmt.Println("{\"result\":\"Validation succeeded\"}")
		}

	default:
		logx.Fatal(fmt.Errorf("'%s': %w", output, cmdutil.ErrUnknownOutputFormat))
	}
}

func PrintResults(output string, err error, conf any, configFile string) {
	ptrPaths := santhosh.GetPtrPaths(err)

	switch output {
	case "text":
		printTextResults(ptrPaths, err, conf, configFile)

	case "json":
		printJSONResults(ptrPaths, err, conf, configFile)

	default:
		logx.Fatal(fmt.Errorf("'%s': %w", output, cmdutil.ErrUnknownOutputFormat))
	}
}

func printTextResults(ptrPaths [][]any, err error, conf any, configFile string) {
	fmt.Printf("CONFIG FILE %s\n", configFile)

	for _, path := range ptrPaths {
		value, serr := santhosh.GetValueAtPath(conf, path)
		if serr != nil {
			logx.Fatal(serr)
		}

		fmt.Printf(
			"path '%s' contains an invalid configuration value: %+v\n",
			santhosh.JoinPtrPath(path),
			value,
		)
	}

	fmt.Println(err)
}

func printJSONResults(ptrPaths [][]any, err error, conf any, configFile string) {
	for _, path := range ptrPaths {
		value, serr := santhosh.GetValueAtPath(conf, path)
		if serr != nil {
			logx.Fatal(serr)
		}

		jv, jerr := jsonx.Marshal(
			map[string]any{
				"file":    configFile,
				"message": "path contains an invalid configuration value",
				"path":    santhosh.JoinPtrPath(path),
				"value":   value,
			},
		)
		if jerr != nil {
			logx.Fatal(jerr)
		}

		fmt.Println(jv)
	}

	jv, jerr := jsonx.Marshal(err)
	if jerr != nil {
		logx.Fatal(jerr)
	}

	fmt.Println(jv)
}
