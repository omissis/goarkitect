package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/omissis/goarkitect/cmd/cmdutil"
	"github.com/omissis/goarkitect/internal/jsonx"
	"github.com/omissis/goarkitect/internal/logx"
)

func NewVersionCommand(output *string, versions map[string]string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information about goarkitect",
		Args:  cobra.ExactArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			if output == nil {
				return cmdutil.ErrNoOutputFormat
			}

			switch *output {
			case "text":
				for k, v := range versions {
					fmt.Printf("%s: %s\n", k, v)
				}
			case "json":
				fmt.Println(
					jsonx.Marshal(
						versions,
					),
				)
			default:
				logx.Fatal(fmt.Errorf("'%s': %w", *output, cmdutil.ErrUnknownOutputFormat))
			}

			return nil
		},
	}
}
