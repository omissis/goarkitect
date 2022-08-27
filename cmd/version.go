package cmd

import (
	"fmt"
	"strings"

	"github.com/omissis/goarkitect/internal/cli"
	"github.com/omissis/goarkitect/internal/jsonx"
	"github.com/omissis/goarkitect/internal/logx"
)

func NewVersionCommand(output *string, versions map[string]string) cli.Command {
	return &versionCommand{
		output:   output,
		versions: versions,
	}
}

type versionCommand struct {
	output   *string
	versions map[string]string
}

func (vc *versionCommand) Name() string {
	return "version"
}

func (vc *versionCommand) Help() string {
	return "TBD"
}

func (vc *versionCommand) Run(args []string) error {
	switch *vc.output {
	case "text":
		for k, v := range vc.versions {
			fmt.Printf("%s: %s\n", strings.Title(k), v)
		}
	case "json":
		fmt.Println(
			jsonx.Marshal(
				vc.versions,
			),
		)
	default:
		logx.Fatal(fmt.Errorf("unknown output format: '%s'", vc.output))
	}

	return nil
}

func (vc *versionCommand) Synopsis() string {
	return "Print version information"
}
