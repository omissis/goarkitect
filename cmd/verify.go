package cmd

import (
	"path/filepath"

	"github.com/omissis/goarkitect/cmd/verify"
	"github.com/omissis/goarkitect/internal/config"
	"github.com/spf13/cobra"
)

func NewVerifyCommand(output *string) *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Verify the ruleset against a project",
		RunE: func(_ *cobra.Command, args []string) error {
			if output == nil {
				return ErrNoOutputFormat
			}

			if len(args) == 0 {
				args = append(args, filepath.Join(getWd(), ".goarkitect.yaml"))
			}

			cfs := listConfigFiles(args)
			if len(cfs) == 0 {
				return ErrNoConfigFileFound
			}

			hasErrors := error(nil)
			for _, cf := range cfs {
				conf := loadConfig[config.Root](cf)

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
