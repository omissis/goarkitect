package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/omissis/goarkitect/cmd/cmdutil"
	jsonx "github.com/omissis/goarkitect/internal/x/json"
	logx "github.com/omissis/goarkitect/internal/x/log"
)

func NewVersionCommand(versions map[string]string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information about goarkitect",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			output := cmd.Flag("output").Value.String()

			switch output {
			case "text":
				for k, v := range versions {
					fmt.Printf("%s: %s\n", k, v)
				}
			case "json":
				jv, jerr := jsonx.Marshal(versions)
				if jerr != nil {
					logx.Fatal(jerr)
				}

				fmt.Println(jv)
			default:
				logx.Fatal(fmt.Errorf("'%s': %w", output, cmdutil.ErrUnknownOutputFormat))
			}

			return nil
		},
	}
}
