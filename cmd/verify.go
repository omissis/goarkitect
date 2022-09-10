package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/omissis/goarkitect/cmd/cmdutil"
	"github.com/omissis/goarkitect/cmd/verify"
	"github.com/omissis/goarkitect/internal/config"
)

func NewVerifyCommand(output *string) *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Verify the ruleset against a project",
		RunE: func(_ *cobra.Command, args []string) error {
			if output == nil {
				return cmdutil.ErrNoOutputFormat
			}

			if len(args) == 0 {
				args = append(args, filepath.Join(cmdutil.GetWd(), ".goarkitect.yaml"))
			}

			cfs := cmdutil.ListConfigFiles(args)
			if len(cfs) == 0 {
				return cmdutil.ErrNoConfigFileFound
			}

			hasErrors := error(nil)
			for _, cf := range cfs {
				conf := cmdutil.LoadConfig[config.Root](cf)

				results := config.Execute(conf)

				verify.PrintResults(*output, cf, results)

				if verify.HasErrors(results) {
					hasErrors = verify.ErrProjectDoesNotRespectRules
				}
			}

			return hasErrors
		},
	}
}
