package cmd

import (
	"fmt"
	"strings"

	"github.com/omissis/goarkitect/internal/jsonx"
	"github.com/omissis/goarkitect/internal/logx"
	"github.com/spf13/cobra"
)

func NewVersionCommand(output *string, versions map[string]string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information about goarkitect",
		Args:  cobra.ExactArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			if output == nil {
				return ErrNoOutputFormat
			}

			switch *output {
			case "text":
				for k, v := range versions {
					fmt.Printf("%s: %s\n", strings.Title(k), v)
				}
			case "json":
				fmt.Println(
					jsonx.Marshal(
						versions,
					),
				)
			default:
				logx.Fatal(fmt.Errorf("unknown output format: '%s', supported ones are: json, text", *output))
			}

			return nil
		},
	}
}
